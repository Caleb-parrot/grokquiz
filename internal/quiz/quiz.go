// Package quiz turns a Grokipedia search hit into a four-choice question.
// It does not fetch or store anything.
package quiz

import (
	"html"
	"math/rand/v2"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	minPrompt = 40
	maxPrompt = 320
)

// Category is one row in the start menu.
// Query is sent to Grokipedia as a live search. It is not a list of questions.
type Category struct {
	Name  string
	Query string
}

// Hit is the slice of a search result the quiz needs.
type Hit struct {
	Title   string
	Slug    string
	Snippet string
}

// Question is one multiple-choice prompt. Answer is an index into Choices.
type Question struct {
	Prompt  string
	Choices []string
	Answer  int
	Topic   string
	Key     string
}

// Fold collapses a title or query for comparison.
func Fold(s string) string {
	s = strings.ToLower(strings.ReplaceAll(s, "_", " "))
	var b strings.Builder
	space := false
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			space = false
			continue
		}
		if !space && b.Len() > 0 {
			b.WriteByte(' ')
			space = true
		}
	}
	return strings.TrimSpace(b.String())
}

// ByName finds a category, ignoring case.
func ByName(name string) (Category, bool) {
	want := Fold(name)
	for _, c := range Categories {
		if Fold(c.Name) == want {
			return c, true
		}
	}
	return Category{}, false
}

var (
	reImg    = regexp.MustCompile(`!\[[^\]]*\]\([^)]*\)`)
	reEmpty  = regexp.MustCompile(`\[\s*\]\([^)]*\)`)
	reLink   = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
	reHead   = regexp.MustCompile(`(?m)^#{1,6}\s*`)
	reTicks  = regexp.MustCompile("`+")
	reSpace  = regexp.MustCompile(`\s+`)
	reTheSub = regexp.MustCompile(`(?i)\b(?:the|a|an)\s+this subject\b`)
)

// Clean strips markdown and citation debris from a snippet.
func Clean(snippet string) string {
	s := html.UnescapeString(snippet)
	s = reImg.ReplaceAllString(s, " ")
	s = reEmpty.ReplaceAllString(s, " ")
	s = reLink.ReplaceAllString(s, "$1")
	s = reHead.ReplaceAllString(s, "")
	s = reTicks.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "\u00a0", " ")
	s = reSpace.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

// Redact replaces whole-word titles with "this subject".
// Longer titles are applied first.
func Redact(s string, titles ...string) string {
	ordered := append([]string{}, titles...)
	sortByLen(ordered)
	out := s
	for _, title := range ordered {
		title = strings.TrimSpace(title)
		if title == "" || utf8.RuneCountInString(title) < 2 {
			continue
		}
		re := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(title) + `(?:'s|’s)?\b`)
		out = re.ReplaceAllStringFunc(out, func(m string) string {
			low := strings.ToLower(m)
			if strings.HasSuffix(low, "'s") || strings.HasSuffix(low, "’s") {
				return "this subject's"
			}
			return "this subject"
		})
	}
	out = reTheSub.ReplaceAllString(out, "this subject")
	out = reSpace.ReplaceAllString(out, " ")
	out = strings.TrimSpace(out)
	if out == "" {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(out), "this subject") {
		return "This subject" + out[len("this subject"):]
	}
	r := []rune(out)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func sortByLen(ss []string) {
	for i := 1; i < len(ss); i++ {
		j := i
		for j > 0 && utf8.RuneCountInString(ss[j]) > utf8.RuneCountInString(ss[j-1]) {
			ss[j], ss[j-1] = ss[j-1], ss[j]
			j--
		}
	}
}

// Build makes a question whose correct choice is answer.Title.
// distractors are other topic names from the same category.
func Build(answer Hit, distractors []string, rng *rand.Rand) (Question, error) {
	if strings.TrimSpace(answer.Title) == "" {
		return Question{}, errShort
	}
	prompt := Clean(answer.Snippet)
	if !leadsWith(prompt, answer.Title) {
		return Question{}, errShort
	}
	choices := make([]string, 0, 4)
	choices = append(choices, answer.Title)
	seen := map[string]bool{Fold(answer.Title): true}
	for _, d := range distractors {
		d = strings.TrimSpace(d)
		k := Fold(d)
		if d == "" || seen[k] {
			continue
		}
		seen[k] = true
		choices = append(choices, d)
		if len(choices) == 4 {
			break
		}
	}
	if len(choices) < 4 {
		return Question{}, errShort
	}
	// Redact only the answer. Replacing distractor names with "this subject"
	// makes a different choice look correct.
	prompt = Redact(prompt, answer.Title)
	prompt = trimPrompt(prompt)
	if utf8.RuneCountInString(prompt) < minPrompt {
		return Question{}, errShort
	}
	if containsFold(prompt, answer.Title) {
		return Question{}, errShort
	}
	rng.Shuffle(len(choices), func(i, j int) {
		choices[i], choices[j] = choices[j], choices[i]
	})
	ans := -1
	want := Fold(answer.Title)
	for i, c := range choices {
		if Fold(c) == want {
			ans = i
			break
		}
	}
	if ans < 0 {
		return Question{}, errShort
	}
	return Question{
		Prompt:  prompt,
		Choices: choices,
		Answer:  ans,
		Topic:   answer.Title,
		Key:     want,
	}, nil
}

func leadsWith(snippet, title string) bool {
	got := Fold(snippet)
	want := Fold(title)
	if want == "" {
		return false
	}
	return strings.HasPrefix(got, want) ||
		strings.HasPrefix(got, "the "+want) ||
		strings.HasPrefix(got, "a "+want) ||
		strings.HasPrefix(got, "an "+want)
}

func trimPrompt(s string) string {
	r := []rune(s)
	if len(r) <= maxPrompt {
		return s
	}
	cut := maxPrompt
	for i := cut; i > minPrompt; i-- {
		if r[i] == ' ' {
			cut = i
			break
		}
	}
	return strings.TrimSpace(string(r[:cut])) + "…"
}

func containsFold(s, title string) bool {
	title = strings.TrimSpace(title)
	if utf8.RuneCountInString(title) < 2 {
		return false
	}
	re := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(title) + `\b`)
	return re.MatchString(s)
}

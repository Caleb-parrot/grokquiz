package quiz

import (
	"math/rand/v2"
	"strings"
	"testing"
)

func TestFold(t *testing.T) {
	if Fold("World_War II!") != "world war ii" {
		t.Fatalf("fold: %q", Fold("World_War II!"))
	}
}

func TestRedactTitle(t *testing.T) {
	got := Redact("The Moon is Earth's only natural satellite, orbiting at 384,400 km.", "Moon")
	if strings.Contains(strings.ToLower(got), "moon") {
		t.Fatalf("title leaked: %s", got)
	}
	if !strings.HasPrefix(got, "This subject") {
		t.Fatalf("lead: %s", got)
	}
}

func TestRedactPluralVerb(t *testing.T) {
	got := Redact("Insect pheromones are chemical substances produced by insects for communication.", "Insect pheromones")
	want := "This subject is chemical substances produced by insects for communication."
	if got != want {
		t.Fatalf("got %q", got)
	}
	mars := Redact("Mars is the fourth planet from the Sun, with a thin atmosphere and two small moons.", "Mars")
	if !strings.HasPrefix(mars, "This subject is ") {
		t.Fatalf("mars %q", mars)
	}
}

func TestLeadRejectsLongerWord(t *testing.T) {
	if leadsWith("Bird nests also function as indicators of pollution in the local environment.", "Bird nest") {
		t.Fatal("bird nests matched bird nest")
	}
	if !leadsWith("Bird nest is a structure built by birds to hold their eggs and young.", "Bird nest") {
		t.Fatal("exact title should lead")
	}
}

func TestRedactPossessive(t *testing.T) {
	got := Redact("Mars's orbit is 687 days long and was measured in 1609.", "Mars")
	if strings.Contains(got, "Mars") {
		t.Fatalf("leaked: %s", got)
	}
	if !strings.Contains(got, "this subject's") && !strings.Contains(got, "This subject's") {
		t.Fatalf("possessive: %s", got)
	}
}

func TestBuildHidesAnswer(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	q, err := Build(Hit{
		Title:   "Mars",
		Slug:    "Mars",
		Snippet: "Mars is the fourth planet from the Sun, with a thin atmosphere and two small moons.",
	}, []string{"Venus", "Jupiter", "Saturn"}, rng)
	if err != nil {
		t.Fatal(err)
	}
	if q.Answer < 0 || q.Answer > 3 {
		t.Fatalf("answer %d", q.Answer)
	}
	if q.Choices[q.Answer] != "Mars" {
		t.Fatalf("choices %v answer %d", q.Choices, q.Answer)
	}
	if containsFold(q.Prompt, "Mars") {
		t.Fatalf("prompt leaks answer: %s", q.Prompt)
	}
	if len(q.Choices) != 4 {
		t.Fatalf("choices %v", q.Choices)
	}
}

func TestBuildKeepsYearAfterTitle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	q, err := Build(Hit{
		Title:   "World War II",
		Snippet: "World War II (1939–1945) was a total war involving the majority of the world's nations and both major coalitions.",
	}, []string{"Cold War", "Roman Empire", "Silk Road"}, rng)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(q.Prompt, "This subject (1939") {
		t.Fatalf("prompt: %s", q.Prompt)
	}
}

func TestBuildRejectsSnippetThatDoesNotLead(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	_, err := Build(Hit{
		Title:   "Isis",
		Snippet: "A hieroglyph symbolizing the name Aset or Isis from the late Fifth Dynasty of Egypt and later worship.",
	}, []string{"Ra", "Anubis", "Horus"}, rng)
	if err == nil {
		t.Fatal("expected a snippet that does not lead with the title to fail")
	}
}

func TestBuildRejectsThinSnippet(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	_, err := Build(Hit{Title: "Mars", Snippet: "Mars."}, []string{"Venus", "Jupiter", "Saturn"}, rng)
	if err == nil {
		t.Fatal("expected thin snippet to fail")
	}
}

func TestCategoriesHaveRoomForChoices(t *testing.T) {
	if len(Categories) != 10 {
		t.Fatalf("menu %d", len(Categories))
	}
	seen := map[string]bool{}
	for _, c := range Categories {
		if c.Name == "" || len(c.Queries) == 0 {
			t.Fatalf("empty category: %+v", c)
		}
		if seen[c.Name] {
			t.Fatalf("duplicate %s", c.Name)
		}
		seen[c.Name] = true
	}
	if _, ok := ByName("science"); !ok {
		t.Fatal("science")
	}
}

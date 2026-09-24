package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/caleb-parrot/quizgrok/internal/game"
	"github.com/caleb-parrot/quizgrok/internal/grok"
	"github.com/caleb-parrot/quizgrok/internal/quiz"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

type phase int

const (
	phaseMenu phase = iota
	phaseLoad
	phasePlay
	phaseOver
	phaseTrouble
)

type drawnMsg struct {
	gen int
	q   quiz.Question
	err error
}

type model struct {
	client *grok.Client
	cats   []quiz.Category
	cursor int
	phase  phase
	cat    *quiz.Category
	run    game.Run
	gen    int
	note   string
	err    string
	width  int
	height int
}

// Run opens the quiz on the alternate screen.
func Run(c *grok.Client) error {
	// Last Horizon remaps the 16-color slots; keep the cream and blue with true color.
	lipgloss.SetColorProfile(termenv.TrueColor)
	m := model{
		client: c,
		cats:   quiz.Categories,
		run:    game.Run{Picked: -1},
	}
	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	return err
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case drawnMsg:
		if msg.gen != m.gen || m.phase != phaseLoad {
			return m, nil
		}
		if msg.err != nil {
			m.phase = phaseTrouble
			m.err = msg.err.Error()
			return m, nil
		}
		m.run.Accept(msg.q)
		m.cursor = 0
		m.phase = phasePlay
		m.err = ""
		return m, nil
	case tea.KeyMsg:
		return m.onKey(msg.String())
	}
	return m, nil
}

func (m model) onKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "q", "ctrl+c":
		return m, tea.Quit
	}
	switch m.phase {
	case phaseMenu:
		return m.onMenu(key)
	case phaseLoad:
		if key == "esc" {
			m.toMenu()
		}
		return m, nil
	case phasePlay:
		return m.onPlay(key)
	case phaseOver:
		switch key {
		case "enter":
			m.run.NewRun()
			m.note = ""
			return m, m.startDraw()
		case "esc":
			m.toMenu()
		}
		return m, nil
	case phaseTrouble:
		switch key {
		case "enter":
			m.note = ""
			return m, m.startDraw()
		case "esc":
			m.toMenu()
		}
		return m, nil
	}
	return m, nil
}

func (m model) onMenu(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "up", "k":
		if len(m.cats) == 0 {
			return m, nil
		}
		m.cursor--
		if m.cursor < 0 {
			m.cursor = len(m.cats) - 1
		}
	case "down", "j":
		if len(m.cats) == 0 {
			return m, nil
		}
		m.cursor++
		if m.cursor >= len(m.cats) {
			m.cursor = 0
		}
	case "enter":
		return m.choose(m.cursor)
	default:
		if n, ok := menuNumber(key); ok {
			return m.choose(n)
		}
	}
	return m, nil
}

func (m model) choose(i int) (tea.Model, tea.Cmd) {
	if i < 0 || i >= len(m.cats) {
		return m, nil
	}
	m.cursor = i
	cat := m.cats[i]
	m.cat = &cat
	m.run.NewRun()
	m.note = ""
	return m, m.startDraw()
}

func (m model) onPlay(key string) (tea.Model, tea.Cmd) {
	if m.run.Q == nil {
		m.toMenu()
		return m, nil
	}
	switch key {
	case "esc":
		m.toMenu()
		return m, nil
	case "up", "k":
		m.cursor--
		if m.cursor < 0 {
			m.cursor = len(m.run.Q.Choices) - 1
		}
		return m, nil
	case "down", "j":
		m.cursor++
		if m.cursor >= len(m.run.Q.Choices) {
			m.cursor = 0
		}
		return m, nil
	case "enter":
		return m.answer(m.cursor)
	default:
		if n, ok := choiceKey(key); ok {
			return m.answer(n)
		}
	}
	return m, nil
}

func (m model) answer(i int) (tea.Model, tea.Cmd) {
	if m.run.Q == nil || i < 0 || i >= len(m.run.Q.Choices) {
		return m, nil
	}
	if m.run.Answer(i) {
		m.note = fmt.Sprintf("Correct — streak %d", m.run.Streak)
		return m, m.startDraw()
	}
	m.note = ""
	m.phase = phaseOver
	return m, nil
}

func (m *model) toMenu() {
	m.gen++
	m.phase = phaseMenu
	m.note = ""
	m.err = ""
	m.run.NewRun()
	if m.cat != nil {
		for i, c := range m.cats {
			if c.Name == m.cat.Name {
				m.cursor = i
				break
			}
		}
	}
}

func (m *model) startDraw() tea.Cmd {
	if m.cat == nil {
		return nil
	}
	m.phase = phaseLoad
	m.gen++
	gen := m.gen
	cat := *m.cat
	avoid := copyMap(m.run.Used)
	client := m.client
	return func() tea.Msg {
		q, err := client.Draw(context.Background(), cat, avoid)
		return drawnMsg{gen: gen, q: q, err: err}
	}
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Quizgrok\n"
	}
	box := frameStyle.Render(m.body())
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

func (m model) inner() int {
	w := m.width - 10
	if w < 36 {
		w = 36
	}
	if w > 72 {
		w = 72
	}
	return w
}

func (m model) body() string {
	switch m.phase {
	case phaseLoad:
		return m.viewLoad()
	case phasePlay:
		return m.viewPlay()
	case phaseOver:
		return m.viewOver()
	case phaseTrouble:
		return m.viewTrouble()
	default:
		return m.viewMenu()
	}
}

func (m model) viewMenu() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("QUIZGROK"))
	b.WriteByte('\n')
	b.WriteString(mutedStyle.Render("One miss ends the run."))
	b.WriteString("\n\n")
	b.WriteString(spread(bodyStyle.Render("Best this sitting"), goldStyle.Render(fmt.Sprintf("%d", m.run.Best)), m.inner()))
	b.WriteString("\n\n")
	for i, c := range m.cats {
		line := fmt.Sprintf(" %2s  %s", menuKey(i), c.Name)
		if i == m.cursor {
			b.WriteString(selStyle.Width(m.inner()).Render(line))
		} else {
			b.WriteString(bodyStyle.Width(m.inner()).Render(line))
		}
		b.WriteByte('\n')
	}
	b.WriteByte('\n')
	b.WriteString(mutedStyle.Render("j/k move    enter start    q quit"))
	return b.String()
}

func (m model) viewLoad() string {
	var b strings.Builder
	b.WriteString(m.head())
	b.WriteString("\n\n")
	if m.note != "" {
		b.WriteString(okStyle.Render(m.note))
		b.WriteString("\n\n")
	}
	b.WriteString(bodyStyle.Render("Drawing a question from Grokipedia…"))
	b.WriteString("\n\n")
	b.WriteString(mutedStyle.Render("esc menu    q quit"))
	return b.String()
}

func (m model) viewPlay() string {
	q := m.run.Q
	if q == nil {
		return m.viewLoad()
	}
	var b strings.Builder
	b.WriteString(m.head())
	b.WriteString("\n\n")
	if m.note != "" {
		b.WriteString(okStyle.Render(m.note))
		b.WriteByte('\n')
	}
	b.WriteString(mutedStyle.Render("Which topic fits this fact?"))
	b.WriteString("\n\n")
	b.WriteString(bodyStyle.Width(m.inner()).Render(q.Prompt))
	b.WriteString("\n\n")
	letters := []string{"A", "B", "C", "D"}
	for i, choice := range q.Choices {
		letter := letters[i]
		row := fmt.Sprintf(" %s  %s", letter, choice)
		if i == m.cursor {
			b.WriteString(selStyle.Width(m.inner()).Render(row))
		} else {
			b.WriteString(letterStyle.Render(letter))
			b.WriteString(bodyStyle.Render("  " + choice))
		}
		b.WriteByte('\n')
	}
	b.WriteByte('\n')
	b.WriteString(mutedStyle.Render("a-d answer    j/k move    enter pick    esc menu"))
	return b.String()
}

func (m model) viewOver() string {
	var b strings.Builder
	b.WriteString(missStyle.Render("Miss."))
	b.WriteString("\n\n")
	if m.run.Q != nil {
		b.WriteString(bodyStyle.Render("That fact is about " + m.run.Q.Topic + "."))
		b.WriteByte('\n')
		if m.run.Picked >= 0 && m.run.Picked < len(m.run.Q.Choices) {
			b.WriteString(bodyStyle.Render("You answered " + m.run.Q.Choices[m.run.Picked] + "."))
			b.WriteByte('\n')
		}
	}
	b.WriteByte('\n')
	b.WriteString(spread(
		bodyStyle.Render(fmt.Sprintf("Streak %d", m.run.Streak)),
		goldStyle.Render(fmt.Sprintf("best %d", m.run.Best)),
		m.inner(),
	))
	b.WriteString("\n\n")
	b.WriteString(mutedStyle.Render("enter again    esc categories    q quit"))
	return b.String()
}

func (m model) viewTrouble() string {
	var b strings.Builder
	b.WriteString(m.head())
	b.WriteString("\n\n")
	b.WriteString(missStyle.Render("Couldn't draw a question."))
	b.WriteString("\n\n")
	msg := m.err
	if msg == "" {
		msg = "Grokipedia did not return a usable article."
	}
	b.WriteString(mutedStyle.Width(m.inner()).Render(msg))
	b.WriteString("\n\n")
	b.WriteString(mutedStyle.Render("enter retry    esc categories    q quit"))
	return b.String()
}

func (m model) head() string {
	name := "Quizgrok"
	if m.cat != nil {
		name = m.cat.Name
	}
	return spread(titleStyle.Render(name), goldStyle.Render(fmt.Sprintf("streak %d", m.run.Streak)), m.inner())
}

func spread(left, right string, width int) string {
	gap := width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 1 {
		return left + "  " + right
	}
	return left + strings.Repeat(" ", gap) + right
}

func menuKey(i int) string {
	if i == 9 {
		return "0"
	}
	return fmt.Sprintf("%d", i+1)
}

func menuNumber(key string) (int, bool) {
	if key == "0" {
		return 9, true
	}
	if len(key) != 1 || key[0] < '1' || key[0] > '9' {
		return 0, false
	}
	return int(key[0] - '1'), true
}

func choiceKey(key string) (int, bool) {
	switch key {
	case "a", "A", "1":
		return 0, true
	case "b", "B", "2":
		return 1, true
	case "c", "C", "3":
		return 2, true
	case "d", "D", "4":
		return 3, true
	default:
		return 0, false
	}
}

func copyMap(in map[string]bool) map[string]bool {
	out := make(map[string]bool, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

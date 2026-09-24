package grok

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/benoute/grokipedia-mcp/pkg/grokipedia"
	"github.com/caleb-parrot/grokquiz/internal/quiz"
)

func fact(title string) []grokipedia.SearchResult {
	return []grokipedia.SearchResult{{
		Title:   title,
		Slug:    title,
		Snippet: title + " was recorded in 1901 with a span of 42 units and a single well known form.",
	}}
}

func TestPickExactSkipsLookalikes(t *testing.T) {
	hit, ok := pickExact("Mercury", []grokipedia.SearchResult{
		{Title: "Mercury Technologies", Snippet: "a software company founded in 1999 with offices in three cities"},
		{Title: "Mercury", Snippet: "Mercury is the smallest planet in the Solar System and the closest to the Sun."},
	})
	if !ok || hit.Title != "Mercury" {
		t.Fatalf("hit %+v ok %v", hit, ok)
	}
}

func TestDrawSkipsRateLimit(t *testing.T) {
	n := 0
	c := &Client{search: func(_ context.Context, query string, _, _ int) ([]grokipedia.SearchResult, error) {
		n++
		if n < 3 {
			return nil, fmt.Errorf("API error: HTTP 429")
		}
		return fact(query), nil
	}}
	cat := quiz.Category{Name: "Space", Queries: []string{"Mars", "Venus", "Jupiter", "Saturn"}}
	q, err := c.Draw(context.Background(), cat, nil)
	if err != nil {
		t.Fatal(err)
	}
	if n < 3 {
		t.Fatalf("calls %d", n)
	}
	if strings.Contains(q.Prompt, q.Topic) {
		t.Fatalf("leaked %s in %s", q.Topic, q.Prompt)
	}
}

func TestDrawRateLimitFails(t *testing.T) {
	c := &Client{search: func(context.Context, string, int, int) ([]grokipedia.SearchResult, error) {
		return nil, fmt.Errorf("API error: HTTP 503")
	}}
	cat := quiz.Category{Name: "Space", Queries: []string{"Mars", "Venus", "Jupiter", "Saturn"}}
	_, err := c.Draw(context.Background(), cat, nil)
	if err == nil || !strings.Contains(err.Error(), "503") {
		t.Fatal(err)
	}
}

func TestDrawContinuesAfterNamedTopics(t *testing.T) {
	c := &Client{search: func(_ context.Context, query string, _, _ int) ([]grokipedia.SearchResult, error) {
		if query != "Space" {
			return []grokipedia.SearchResult{{Title: "Nope", Snippet: "not the article"}}, nil
		}
		return []grokipedia.SearchResult{
			{Title: "Europa", Slug: "Europa", Snippet: "Europa is an icy moon of Jupiter with a subsurface ocean about 100 kilometers deep."},
			{Title: "Titan", Slug: "Titan", Snippet: "Titan is the largest moon of Saturn and the only moon known to have a dense atmosphere."},
			{Title: "Ganymede", Slug: "Ganymede", Snippet: "Ganymede is the largest moon in the Solar System and the only moon with its own magnetic field."},
			{Title: "Callisto", Slug: "Callisto", Snippet: "Callisto is a moon of Jupiter with an ancient cratered surface and a possible subsurface ocean."},
		}, nil
	}}
	cat := quiz.Category{Name: "Space", Queries: []string{"Mars", "Venus", "Jupiter", "Saturn"}}
	avoid := map[string]bool{"mars": true, "venus": true, "jupiter": true, "saturn": true}
	q, err := c.Draw(context.Background(), cat, avoid)
	if err != nil {
		t.Fatal(err)
	}
	switch q.Topic {
	case "Europa", "Titan", "Ganymede", "Callisto":
	default:
		t.Fatalf("topic %q", q.Topic)
	}
}

func TestDrawSkipsAskedTopic(t *testing.T) {
	var asked []string
	c := &Client{search: func(_ context.Context, query string, _, _ int) ([]grokipedia.SearchResult, error) {
		asked = append(asked, query)
		if query == "Mars" {
			return fact(query), nil
		}
		return []grokipedia.SearchResult{{Title: "Nope", Snippet: "not the article"}}, nil
	}}
	cat := quiz.Category{Name: "Space", Queries: []string{"Mars", "Venus", "Jupiter", "Saturn"}}
	_, err := c.Draw(context.Background(), cat, map[string]bool{"mars": true})
	if err == nil {
		t.Fatal("expected failure when the only good topic was already asked")
	}
	for _, q := range asked {
		if q == "Mars" {
			t.Fatal("searched an asked topic")
		}
	}
}

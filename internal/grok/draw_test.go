package grok

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/benoute/grokipedia-mcp/pkg/grokipedia"
	"github.com/caleb-parrot/quizgrok/internal/quiz"
)

func batch() []grokipedia.SearchResult {
	titles := []string{"Europa", "Titan", "Ganymede", "Callisto"}
	out := make([]grokipedia.SearchResult, len(titles))
	for i, title := range titles {
		out[i] = grokipedia.SearchResult{
			Title:   title,
			Slug:    title,
			Snippet: title + " is a moon in astronomy, recorded in 1901 with a span of 42 units and a single well known form.",
		}
	}
	return out
}

func TestDrawSkipsRateLimit(t *testing.T) {
	n := 0
	c := &Client{search: func(context.Context, string, int, int) ([]grokipedia.SearchResult, error) {
		n++
		if n < 3 {
			return nil, fmt.Errorf("API error: HTTP 429")
		}
		return batch(), nil
	}}
	q, err := c.Draw(context.Background(), quiz.Category{Name: "Space", Queries: []string{"astronomy"}}, nil)
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
	_, err := c.Draw(context.Background(), quiz.Category{Name: "Space", Queries: []string{"astronomy"}}, nil)
	if err == nil || !strings.Contains(err.Error(), "503") {
		t.Fatal(err)
	}
}

func TestDrawUsesRandomSearchHits(t *testing.T) {
	var queries []string
	c := &Client{search: func(_ context.Context, query string, _, _ int) ([]grokipedia.SearchResult, error) {
		queries = append(queries, query)
		return batch(), nil
	}}
	q, err := c.Draw(context.Background(), quiz.Category{Name: "Space", Queries: []string{"astronomy"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(queries) != 1 || queries[0] != "astronomy" {
		t.Fatalf("queries %v", queries)
	}
	switch q.Topic {
	case "Europa", "Titan", "Ganymede", "Callisto":
	default:
		t.Fatalf("topic %q", q.Topic)
	}
	for _, choice := range q.Choices {
		switch choice {
		case "Europa", "Titan", "Ganymede", "Callisto":
		default:
			t.Fatalf("choice %q is not from the search", choice)
		}
	}
}

func TestUsableHitsRequiresWholePhrase(t *testing.T) {
	hits := usableHits([]grokipedia.SearchResult{
		{
			Title:   "Amberjack Hole",
			Slug:    "Amberjack_Hole",
			Snippet: "Amberjack Hole is a submarine blue hole, a type of karst sinkhole, located in the Gulf of Mexico.",
		},
		{
			Title:   "Black hole",
			Slug:    "Black_hole",
			Snippet: "A black hole is a region of spacetime where gravity is so strong that nothing can escape it.",
		},
	}, nil, "black hole")
	if len(hits) != 1 || hits[0].Title != "Black hole" {
		t.Fatalf("hits %+v", hits)
	}
}

func TestDrawSkipsAskedTopic(t *testing.T) {
	c := &Client{search: func(context.Context, string, int, int) ([]grokipedia.SearchResult, error) {
		hits := batch()
		hits = append(hits, grokipedia.SearchResult{
			Title:   "Mars",
			Slug:    "Mars",
			Snippet: "Mars is a planet in astronomy, recorded in 1901 with a span of 42 units and a single well known form.",
		})
		return hits, nil
	}}
	q, err := c.Draw(context.Background(), quiz.Category{Name: "Space", Queries: []string{"astronomy"}}, map[string]bool{"mars": true})
	if err != nil {
		t.Fatal(err)
	}
	if q.Topic == "Mars" {
		t.Fatal("asked topic was used")
	}
}

// Package grok draws one quiz question from a live Grokipedia search.
// Responses are not cached and are not written to disk.
package grok

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/benoute/grokipedia-mcp/pkg/grokipedia"
	"github.com/caleb-parrot/grokquiz/internal/quiz"
)

// SearchFunc is grokipedia.Search, or a stand-in in tests.
type SearchFunc func(ctx context.Context, query string, limit, offset int) ([]grokipedia.SearchResult, error)

// Client talks to Grokipedia through pkg/grokipedia.
type Client struct {
	search SearchFunc
}

// New returns a client that uses the library search call.
func New() *Client {
	return &Client{search: librarySearch}
}

func librarySearch(ctx context.Context, query string, limit, offset int) ([]grokipedia.SearchResult, error) {
	return grokipedia.Search(ctx, query,
		grokipedia.WithLimit(limit),
		grokipedia.WithOffset(offset),
	)
}

var errNoMatch = errors.New("no exact article")

// Draw searches until it can build one question for cat.
// avoid holds folded topic keys already answered in this run.
func (c *Client) Draw(ctx context.Context, cat quiz.Category, avoid map[string]bool) (quiz.Question, error) {
	if c.search == nil {
		return quiz.Question{}, errors.New("no search client")
	}
	rng := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(len(cat.Queries))+1))
	order := append([]string{}, cat.Queries...)
	rng.Shuffle(len(order), func(i, j int) {
		order[i], order[j] = order[j], order[i]
	})

	var last error
	tried := 0
	for _, seed := range order {
		if avoid[quiz.Fold(seed)] {
			continue
		}
		if tried == 6 {
			break
		}
		tried++
		hit, err := c.lookup(ctx, seed)
		if err != nil {
			last = err
			continue
		}
		q, err := quiz.Build(hit, distractors(cat, hit.Title, rng), rng)
		if err != nil {
			last = err
			continue
		}
		return q, nil
	}
	if last == nil {
		last = errNoMatch
	}
	return quiz.Question{}, fmt.Errorf("%s: %w", cat.Name, last)
}

func (c *Client) lookup(ctx context.Context, query string) (quiz.Hit, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	results, err := c.search(ctx, query, 5, 0)
	if err != nil {
		return quiz.Hit{}, err
	}
	hit, ok := pickExact(query, results)
	if !ok {
		return quiz.Hit{}, errNoMatch
	}
	return hit, nil
}

func pickExact(query string, results []grokipedia.SearchResult) (quiz.Hit, bool) {
	want := quiz.Fold(query)
	for _, r := range results {
		if quiz.Fold(r.Title) != want {
			continue
		}
		if strings.TrimSpace(r.Snippet) == "" {
			continue
		}
		return quiz.Hit{Title: r.Title, Slug: r.Slug, Snippet: r.Snippet}, true
	}
	return quiz.Hit{}, false
}

func distractors(cat quiz.Category, answer string, rng *rand.Rand) []string {
	ans := quiz.Fold(answer)
	pool := make([]string, 0, len(cat.Queries))
	for _, q := range cat.Queries {
		if quiz.Fold(q) == ans {
			continue
		}
		pool = append(pool, q)
	}
	rng.Shuffle(len(pool), func(i, j int) {
		pool[i], pool[j] = pool[j], pool[i]
	})
	if len(pool) > 3 {
		pool = pool[:3]
	}
	return pool
}

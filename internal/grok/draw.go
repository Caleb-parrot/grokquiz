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
	q, err := c.drawOpen(ctx, cat, avoid, rng)
	if err == nil {
		return q, nil
	}
	if last == nil {
		last = err
	}
	return quiz.Question{}, fmt.Errorf("%s: %w", cat.Name, last)
}

// drawOpen keeps a run going after the named topics have been asked.
// It searches the category at a random offset and builds a question from
// whatever articles come back. A long streak is allowed to repeat a topic
// rather than stop.
func (c *Client) drawOpen(ctx context.Context, cat quiz.Category, avoid map[string]bool, rng *rand.Rand) (quiz.Question, error) {
	var last error
	for attempt := 0; attempt < 4; attempt++ {
		q, err := c.oneOpen(ctx, cat, avoid, rng, rng.IntN(160))
		if err == nil {
			return q, nil
		}
		last = err
	}
	q, err := c.oneOpen(ctx, cat, nil, rng, rng.IntN(160))
	if err == nil {
		return q, nil
	}
	if last == nil {
		last = err
	}
	if last == nil {
		last = errNoMatch
	}
	return quiz.Question{}, last
}

func (c *Client) oneOpen(ctx context.Context, cat quiz.Category, avoid map[string]bool, rng *rand.Rand, offset int) (quiz.Question, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	results, err := c.search(ctx, cat.Name, 8, offset)
	if err != nil {
		return quiz.Question{}, err
	}
	hits := usableHits(results, avoid)
	if len(hits) == 0 {
		return quiz.Question{}, errNoMatch
	}
	rng.Shuffle(len(hits), func(i, j int) {
		hits[i], hits[j] = hits[j], hits[i]
	})
	var last error
	for i, answer := range hits {
		names := make([]string, 0, len(hits)-1+3)
		for j, h := range hits {
			if j == i {
				continue
			}
			names = append(names, h.Title)
		}
		names = append(names, distractors(cat, answer.Title, rng)...)
		q, err := quiz.Build(answer, names, rng)
		if err == nil {
			return q, nil
		}
		last = err
	}
	if last == nil {
		last = errNoMatch
	}
	return quiz.Question{}, last
}

func usableHits(results []grokipedia.SearchResult, avoid map[string]bool) []quiz.Hit {
	var hits []quiz.Hit
	seen := map[string]bool{}
	for _, r := range results {
		title := strings.TrimSpace(r.Title)
		if title == "" || strings.TrimSpace(r.Snippet) == "" || junkTitle(title) {
			continue
		}
		key := quiz.Fold(title)
		if key == "" || seen[key] {
			continue
		}
		if avoid[key] || (r.Slug != "" && avoid[quiz.Fold(r.Slug)]) {
			continue
		}
		seen[key] = true
		hits = append(hits, quiz.Hit{Title: title, Slug: r.Slug, Snippet: r.Snippet})
	}
	return hits
}

func junkTitle(title string) bool {
	t := strings.ToLower(title)
	for _, bad := range []string{
		"(film)", "(album)", "(band)", "(ep)", "(song)", "(novel)",
		"(tv", "inc.", "cryptocurrency",
	} {
		if strings.Contains(t, bad) {
			return true
		}
	}
	return false
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

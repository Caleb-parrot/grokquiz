// Package game keeps the in-memory streak for one sitting.
// Nothing here is written to disk.
package game

import "github.com/caleb-parrot/quizez/internal/quiz"

// Run is the current streak. Best lasts until the process exits.
type Run struct {
	Streak int
	Best   int
	Used   map[string]bool
	Q      *quiz.Question
	Picked int
}

// NewRun clears the streak and the topics already asked.
// Best is kept for this sitting.
func (r *Run) NewRun() {
	r.Streak = 0
	r.Used = map[string]bool{}
	r.Q = nil
	r.Picked = -1
}

// Accept stores the question currently on screen.
func (r *Run) Accept(q quiz.Question) {
	r.Q = &q
	r.Picked = -1
}

// Answer reports whether choice i is right.
// A miss leaves the streak where it is and does not record the topic.
func (r *Run) Answer(i int) bool {
	if r.Q == nil || i < 0 || i >= len(r.Q.Choices) {
		return false
	}
	r.Picked = i
	if i != r.Q.Answer {
		return false
	}
	r.Streak++
	if r.Streak > r.Best {
		r.Best = r.Streak
	}
	if r.Used == nil {
		r.Used = map[string]bool{}
	}
	r.Used[r.Q.Key] = true
	return true
}

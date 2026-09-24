package game

import (
	"testing"

	"github.com/caleb-parrot/quizgrok/internal/quiz"
)

func sample() quiz.Question {
	return quiz.Question{
		Prompt:  "This subject orbits the Sun.",
		Choices: []string{"Venus", "Mars", "Jupiter", "Saturn"},
		Answer:  1,
		Topic:   "Mars",
		Key:     "mars",
	}
}

func TestStreakUntilMiss(t *testing.T) {
	var r Run
	r.NewRun()
	r.Accept(sample())
	if !r.Answer(1) {
		t.Fatal("expected correct")
	}
	if r.Streak != 1 || r.Best != 1 {
		t.Fatalf("streak %d best %d", r.Streak, r.Best)
	}
	if !r.Used["mars"] {
		t.Fatal("topic not marked asked")
	}
	r.Accept(sample())
	if r.Answer(0) {
		t.Fatal("expected a miss")
	}
	if r.Streak != 1 || r.Best != 1 {
		t.Fatalf("after miss streak %d best %d", r.Streak, r.Best)
	}
	r.NewRun()
	if r.Streak != 0 || r.Best != 1 {
		t.Fatalf("new run streak %d best %d", r.Streak, r.Best)
	}
	if len(r.Used) != 0 {
		t.Fatal("used should clear")
	}
}

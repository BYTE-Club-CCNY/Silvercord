package cache_test

import (
	"testing"
	"time"

	"main/cache"
)

func TestUsernameKey(t *testing.T) {
	got := cache.UsernameKey("srv1", "usr1")
	want := "silvercord:username:srv1:usr1"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestScoreKey(t *testing.T) {
	got := cache.ScoreKey("srv1", "usr1", "3")
	want := "silvercord:score:srv1:usr1:3"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestLeaderboardKey(t *testing.T) {
	got := cache.LeaderboardKey("srv1", "3")
	want := "silvercord:leaderboard:srv1:3"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestProblemsKey(t *testing.T) {
	got := cache.ProblemsKey("srv1", "usr1")
	want := "silvercord:problems:srv1:usr1"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestTTLValues(t *testing.T) {
	if cache.UsernameTTL != 30*time.Minute {
		t.Fatalf("UsernameTTL = %v, want 30m", cache.UsernameTTL)
	}
	if cache.ScoreTTL != 10*time.Minute {
		t.Fatalf("ScoreTTL = %v, want 10m", cache.ScoreTTL)
	}
	if cache.LeaderboardTTL != 5*time.Minute {
		t.Fatalf("LeaderboardTTL = %v, want 5m", cache.LeaderboardTTL)
	}
	if cache.ProblemsTTL != 10*time.Minute {
		t.Fatalf("ProblemsTTL = %v, want 10m", cache.ProblemsTTL)
	}
}

// Key uniqueness — different inputs must not collide
func TestKeyUniqueness(t *testing.T) {
	keys := []string{
		cache.UsernameKey("s1", "u1"),
		cache.UsernameKey("s1", "u2"),
		cache.UsernameKey("s2", "u1"),
		cache.ScoreKey("s1", "u1", "1"),
		cache.ScoreKey("s1", "u1", "2"),
		cache.LeaderboardKey("s1", "1"),
		cache.LeaderboardKey("s2", "1"),
		cache.ProblemsKey("s1", "u1"),
		cache.ProblemsKey("s1", "u2"),
	}

	seen := make(map[string]bool)
	for _, k := range keys {
		if seen[k] {
			t.Fatalf("duplicate key: %q", k)
		}
		seen[k] = true
	}
}

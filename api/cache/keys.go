package cache

import (
	"fmt"
	"time"
)

const (
	UsernameTTL    = 30 * time.Minute
	ScoreTTL       = 10 * time.Minute
	LeaderboardTTL = 5 * time.Minute
	ProblemsTTL    = 10 * time.Minute
)

func UsernameKey(serverID, userID string) string {
	return fmt.Sprintf("silvercord:username:%s:%s", serverID, userID)
}

func ScoreKey(serverID, userID, season string) string {
	return fmt.Sprintf("silvercord:score:%s:%s:%s", serverID, userID, season)
}

func LeaderboardKey(serverID, season string) string {
	return fmt.Sprintf("silvercord:leaderboard:%s:%s", serverID, season)
}

func ProblemsKey(serverID, userID string) string {
	return fmt.Sprintf("silvercord:problems:%s:%s", serverID, userID)
}

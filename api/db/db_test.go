package db

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		t.Fatalf("failed to set dialect: %v", err)
	}
	if err := goose.Up(db, "migrations"); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	return db
}

func TestMigrationsRun(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Verify all 3 tables exist
	tables := []string{"leetboard_username", "leetboard_scores", "leetboard_problems"}
	for _, table := range tables {
		var name string
		err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		if err != nil {
			t.Errorf("expected table %s to exist, got error: %v", table, err)
		}
	}
}

func TestSeed(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	if err := Seed(db); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	// Verify usernames were seeded
	var userCount int
	db.QueryRow("SELECT COUNT(*) FROM leetboard_username").Scan(&userCount)
	if userCount != 3 {
		t.Errorf("expected 3 users, got %d", userCount)
	}

	// Verify scores were seeded
	var scoreCount int
	db.QueryRow("SELECT COUNT(*) FROM leetboard_scores").Scan(&scoreCount)
	if scoreCount != 3 {
		t.Errorf("expected 3 scores, got %d", scoreCount)
	}

	// Verify problems were seeded
	var problemCount int
	db.QueryRow("SELECT COUNT(*) FROM leetboard_problems").Scan(&problemCount)
	if problemCount != 2 {
		t.Errorf("expected 2 problems, got %d", problemCount)
	}
}

func TestSeedIdempotent(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	if err := Seed(db); err != nil {
		t.Fatalf("first seed failed: %v", err)
	}
	if err := Seed(db); err != nil {
		t.Fatalf("second seed failed: %v", err)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM leetboard_username").Scan(&count)
	if count != 3 {
		t.Errorf("expected 3 users after double seed, got %d", count)
	}
}

func TestUsernameUpsert(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, err := db.Exec("INSERT INTO leetboard_username (server_id, user_id, leetcode_username) VALUES (1, 1, 'old_name')")
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}

	_, err = db.Exec("INSERT OR REPLACE INTO leetboard_username (server_id, user_id, leetcode_username) VALUES (1, 1, 'new_name')")
	if err != nil {
		t.Fatalf("upsert failed: %v", err)
	}

	var username string
	db.QueryRow("SELECT leetcode_username FROM leetboard_username WHERE server_id=1 AND user_id=1").Scan(&username)
	if username != "new_name" {
		t.Errorf("expected 'new_name', got '%s'", username)
	}
}

func TestScoreUpsert(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, err := db.Exec("INSERT INTO leetboard_scores (server_id, user_id, season, score) VALUES (1, 1, 3, 5)")
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}

	_, err = db.Exec("INSERT INTO leetboard_scores (server_id, user_id, season, score) VALUES (1, 1, 3, 10) ON CONFLICT(server_id, user_id, season) DO UPDATE SET score=10")
	if err != nil {
		t.Fatalf("upsert failed: %v", err)
	}

	var score int
	db.QueryRow("SELECT score FROM leetboard_scores WHERE server_id=1 AND user_id=1 AND season=3").Scan(&score)
	if score != 10 {
		t.Errorf("expected score 10, got %d", score)
	}
}

func TestLeaderboardOrder(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	if err := Seed(db); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	rows, err := db.Query("SELECT user_id, score FROM leetboard_scores WHERE server_id=123456789 AND season=3 ORDER BY score DESC")
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	defer rows.Close()

	var prevScore int = 999
	for rows.Next() {
		var userID int64
		var score int
		rows.Scan(&userID, &score)
		if score > prevScore {
			t.Errorf("leaderboard not in descending order: %d came after %d", score, prevScore)
		}
		prevScore = score
	}
}

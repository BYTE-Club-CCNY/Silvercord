package db

import (
	"database/sql"
	"log"
)

func Seed(db *sql.DB) error {
	// Only seed if tables are empty
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM leetboard_username").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		log.Println("[seed] Tables already have data, skipping seed")
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Seed usernames
	_, err = tx.Exec(`
		INSERT INTO leetboard_username (server_id, user_id, leetcode_username) VALUES
		(123456789, 100000001, 'alice_codes'),
		(123456789, 100000002, 'bob_leetcode'),
		(123456789, 100000003, 'charlie_dev')
	`)
	if err != nil {
		return err
	}

	// Seed scores
	_, err = tx.Exec(`
		INSERT INTO leetboard_scores (server_id, user_id, season, score) VALUES
		(123456789, 100000001, 3, 15),
		(123456789, 100000002, 3, 10),
		(123456789, 100000003, 3, 5)
	`)
	if err != nil {
		return err
	}

	// Seed problems
	_, err = tx.Exec(`
		INSERT INTO leetboard_problems (server_id, user_id, link, problem) VALUES
		(123456789, 100000001, 'https://leetcode.com/problems/two-sum/', 'Two Sum'),
		(123456789, 100000002, 'https://leetcode.com/problems/add-two-numbers/', 'Add Two Numbers')
	`)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Println("[seed] Database seeded successfully")
	return nil
}

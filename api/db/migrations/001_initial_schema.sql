-- +goose Up
CREATE TABLE leetboard_username (
    server_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    leetcode_username TEXT NOT NULL,
    PRIMARY KEY (server_id, user_id)
);

CREATE TABLE leetboard_scores (
    server_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    season INTEGER NOT NULL,
    score INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (server_id, user_id, season)
);

CREATE TABLE leetboard_problems (
    server_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    link TEXT NOT NULL,
    problem TEXT NOT NULL,
    PRIMARY KEY (server_id, user_id, problem)
);

-- +goose Down
DROP TABLE leetboard_username;
DROP TABLE leetboard_scores;
DROP TABLE leetboard_problems;

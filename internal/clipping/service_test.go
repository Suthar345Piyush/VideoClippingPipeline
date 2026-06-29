// tests for clipping service

package clipping

import (
	"database/sql"
	"testing"

	"github.com/Suthar345Piyush/videoclippingpipeline/internal/database"
)

// testing schema

const testSchema = `

CREATE TABLE videos (
    id TEXT PRIMARY KEY,
		filename TEXT NOT NULL,
		original_path TEXT NOT NULL,
		duration REAL NOT NULL DEFAULT 0,
		filesize INTEGER NOT NULL DEFAULT 0,
		width INTEGER NOT NULL DEFAULT 0,
		height INTEGER NOT NULL DEFAULT 0,
		fps REAL NOT NULL DEFAULT 0,
		codec TEXT NOT NULL DEFAULT '',
		status TEXT NOT NULL DEFAULT 'pending',
		error_msg TEXT,
		created_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
		updated_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);


CREATE TABLE  clips (
	id TEXT PRIMARY KEY,
	video_id TEXT NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
	clip_path TEXT NOT NULL,
	start_time TEXT NOT NULL,
	end_time TEXT NOT NULL,
	duration REAL GENERATED ALWAYS AS (end_time - start_time) STORED,
	label TEXT NOT NULL DEFAULT '',
	status TEXT NOT NULL DEFAULT 'pending',
	error_msg TEXT,
	created_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
	updated_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);`

// setting up the test db

func setupDB(t *testing.T) *database.Queries {

	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open in-memory db: %v", err)
	}

	if _, err := db.Exec(testSchema); err != nil {
		t.Fatalf("applying the schema: %v", err)
	}

	t.Cleanup(func() { db.Close() })

	return database.New(db)

}

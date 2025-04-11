package schema

import (
	"database/sql"
	"fmt"
)

func CreateUserTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	nyx_code TEXT NOT NULL UNIQUE
);`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUserTable(db *sql.DB) error {
	query := `DROP TABLE users;`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("could not delete the user table")
	}
	return nil
}

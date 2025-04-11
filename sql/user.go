package sql

import (
	"database/sql"
)

type User struct {
	ID        string
	CreatedAt string
	UpdatedAt string
	NyxCode   string
}

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

func CreateUser(db *sql.DB, id string, created_at string, updated_at string, code string) (User, error) {
	query := `INSERT INTO users(id, created_at, updated_at, nyx_code)
			values(?, ?, ?, ?)
			RETURNING *`
	var user User
	err := db.QueryRow(query, id, created_at, updated_at, code).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.NyxCode,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil

}

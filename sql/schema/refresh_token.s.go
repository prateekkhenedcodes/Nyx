package schema

import "database/sql"

func CreateRefreshToken(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS refresh_tokens(
	token TEXT PRIMARY KEY, 
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL, 
	expires_at TEXT NOT NULL, 
	revoked_at TEXT,
	user_id TEXT NOT NULL, 
	FOREIGN KEY (user_id)
	REFERENCES users(id) ON DELETE CASCADE
);`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

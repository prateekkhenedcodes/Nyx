package schema

import "database/sql"

func CreateNyxServer(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS nyx_servers(
	server_id TEXT PRIMARY KEY, 
	created_at TEXT NOT NULL, 
	expires_at TEXT NOT NULL, 
	max_participants INTEGER, 
	active_session BOOLEAN,
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

package schema

import "database/sql"

func CreateNyxServerMessages(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS nyx_server_messages(
	message_id TEXT NOT NULL, 
	sender_alias TEXT NOT NULL, 
	send_at TEXT NOT NULL, 
	content TEXT NOT NULL, 
	server_id TEXT NOT NULL,
	FOREIGN KEY (server_id)
	REFERENCES  nyx_server(server_id) ON DELETE CASCADE
	);`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

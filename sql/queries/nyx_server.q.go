package queries

import "database/sql"

type NyxServer struct {
	ServerId        string
	CreatedAt       string
	ExpiresAt       string
	MaxParticipants int
	ActiveSession   bool
	UserId          string
}

func AddNyxServer(db *sql.DB,
	serverId string,
	createdAt string,
	expiresAt string,
	maxParti int,
	activeSession bool,
	userId string) (NyxServer, error) {
	query := `INSERT INTO nyx_servers(
	server_id, 
	created_at, 
	expires_at,
	max_participants, 
	active_session,
	user_id
	)
	VALUES(?, ?, ?, ?, ?, ?)
	RETURNING *`

	var nyxServer NyxServer
	err := db.QueryRow(query, serverId, createdAt, expiresAt, maxParti, activeSession, userId).Scan(
		&nyxServer.ServerId,
		&nyxServer.CreatedAt,
		&nyxServer.ExpiresAt,
		&nyxServer.MaxParticipants,
		&nyxServer.ActiveSession,
		&nyxServer.UserId,
	)
	if err != nil {
		return NyxServer{}, err
	}
	return nyxServer, nil
}

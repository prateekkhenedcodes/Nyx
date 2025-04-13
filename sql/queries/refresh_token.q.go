package queries

import (
	"database/sql"
)

type RefreshToken struct {
	Token     string
	CreatedAt string
	UpdatedAt string
	UserId    string
	ExpiresAt string
	RevokedAt string
}

func AddRefreshToken(db *sql.DB,
	token string,
	createdAt string,
	updatedAt string,
	userId string,
	expiresAt string,
	revokedAt string) (RefreshToken, error) {
	query := `INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at, revoked_at)
		VALUES(?, ?, ?, ?, ?, ?)
		RETURNING *`

	var refreshToken RefreshToken
	err := db.QueryRow(query, token, createdAt, updatedAt, userId, expiresAt, revokedAt).Scan(
		&refreshToken.Token,
		&refreshToken.CreatedAt,
		&refreshToken.UpdatedAt,
		&refreshToken.UserId,
		&refreshToken.ExpiresAt,
		&refreshToken.RevokedAt,
	)
	if err != nil {
		return RefreshToken{}, err
	}

	return refreshToken, nil

}

func GetUserFromRefreshToken(db *sql.DB, token string) (RefreshToken, error) {
	query := `SELECT * FROM refresh_tokens WHERE token = ?;`

	var refreshToken RefreshToken
	err := db.QueryRow(query, token).Scan(
		&refreshToken.Token,
		&refreshToken.CreatedAt,
		&refreshToken.UpdatedAt,
		&refreshToken.ExpiresAt,
		&refreshToken.RevokedAt,
		&refreshToken.UserId,
	)
	if err != nil {
		return RefreshToken{}, err
	}

	return refreshToken, nil
}

func RevokeRefreshToken(db *sql.DB, token string, updateT string, revokeT string) (RefreshToken, error) {
	query := `	UPDATE refresh_tokens
				SET updated_at = ?,  revoked_at = ?
				WHERE token = ?
				RETURNING *;
			`

	var refreshToken RefreshToken

	err := db.QueryRow(query, updateT, revokeT, token).Scan(
		&refreshToken.Token,
		&refreshToken.CreatedAt,
		&refreshToken.UpdatedAt,
		&refreshToken.ExpiresAt,
		&refreshToken.RevokedAt,
		&refreshToken.UserId,
	)
	if err != nil {
		return RefreshToken{}, err
	}
	return refreshToken, nil
}

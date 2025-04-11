package queries

import "database/sql"

type User struct {
	ID        string
	CreatedAt string
	UpdatedAt string
	NyxCode   string
}

func AddUser(db *sql.DB, id string, created_at string, updated_at string, code string) (User, error) {
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

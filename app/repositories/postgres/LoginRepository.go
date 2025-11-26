package postgres

import (
	"backendUAS/app/models/postgres"
	"backendUAS/databases"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(raw string, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw))
	return err == nil
}

func Authenticate(email string, password string) (*postgres.User, error) {

	var User postgres.User

	err := databases.DatabaseQuery.QueryRow(`
			SELECT u.id, u.username, u.email, u.full_name, u.password_hash, u.role_id, r.name  
			FROM users as u
			JOIN roles as r on u.role_id = r.id
			WHERE u.email = $1
		`, email).Scan(
		&User.ID, &User.Username, &User.Email, &User.FullName, &User.PasswordHash, &User.RoleID, &User.RoleName,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("email tidak ditemukan")
	}

	if !CheckPassword(password, User.PasswordHash) {
		return nil, errors.New("password salah")
	}

	if err != nil {
		return nil, err
	}

	return &User, err

}

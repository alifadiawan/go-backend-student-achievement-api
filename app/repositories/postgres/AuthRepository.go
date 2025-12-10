package postgres

import (
	models "backendUAS/app/models/postgres"
	"backendUAS/databases"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(raw string, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw))
	return err == nil
}

func GetProfile(userId string) (*models.User, error) {

	var User models.User

	err := databases.DatabaseQuery.QueryRow(`
		SELECT 
			u.id, u.username, u.email, u.full_name, u.role_id, r.name  
		FROM 
			users as u
		JOIN 
			roles as r on u.role_id = r.id
		WHERE 
			u.id = $1 
	`, userId).Scan(
		&User.ID, &User.Username, &User.Email, &User.FullName, &User.RoleID, &User.RoleName,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("tidak ditemukan")
	}

	return &User, err

}

func Authenticate(email string, password string) (*models.LoginResponse, error) {
	query := `
        SELECT 
			u.id,
			st.id AS student_id,
			st.student_id AS nim,
			u.username,
			u.email,
			u.full_name,
			u.password_hash,
			r.name AS role_name,
			p.name AS permission_name,
			l.id AS lecturer_id
		FROM users AS u
		JOIN roles AS r ON u.role_id = r.id
		JOIN role_permissions AS rp ON r.id = rp.role_id
		JOIN permissions AS p ON rp.permission_id = p.id
		LEFT JOIN students AS st ON u.id = st.user_id
		LEFT JOIN lecturers AS l ON u.id = l.user_id
        WHERE u.email = $1
   `

	rows, err := databases.DatabaseQuery.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resp models.LoginResponse
	permMap := make(map[string]bool)
	var passwordHash string
	var firstRow = true

	for rows.Next() {
		var id, username, emailDB, fullName, roleName, permName, lecturerID string
		var studentID, nim sql.NullString

		err := rows.Scan(
			&id,
			&studentID,
			&nim,
			&username,
			&emailDB,
			&fullName,
			&passwordHash,
			&roleName,
			&permName,
			&lecturerID,
		)

		if err != nil {
			return nil, err
		}

		if firstRow {
			resp.ID = id
			resp.Email = emailDB
			resp.Username = username
			resp.FullName = fullName
			resp.Role = roleName

			if studentID.Valid {
				resp.StudentID = &studentID.String
			}
			if nim.Valid {
				resp.NIM = &nim.String
			}
			if lecturerID != "" {
				resp.LecturerID = &lecturerID
			}

			firstRow = false
		}

	}

	if resp.ID == "" {
		return nil, errors.New("email tidak ditemukan")
	}

	if !CheckPassword(password, passwordHash) {
		return nil, errors.New("password salah")
	}

	for perm := range permMap {
		resp.Permissions = append(resp.Permissions, perm)
	}

	return &resp, nil
}

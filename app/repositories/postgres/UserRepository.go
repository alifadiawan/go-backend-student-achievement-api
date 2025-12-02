package postgres

import (
	models "backendUAS/app/models/postgres"
	"backendUAS/databases"
	"backendUAS/utils"
)

func GetAllUserRepository() ([]models.User, error) {

	Query, err := databases.DatabaseQuery.Query(`
		SELECT 
			u.id, 
			u.username, 
			u.email, 
			u.full_name, 
			u.role_id, 
			r.name,
			u.is_active,
			u.created_at,
			u.updated_at 
		FROM users as u
		JOIN roles as r on r.id = u.role_id
	`)
	if err != nil {
		return nil, err
	}
	defer Query.Close()

	var Users []models.User

	for Query.Next(){
		var item models.User

		foreach := Query.Scan(
			&item.ID,
			&item.Username,
			&item.Email,
			&item.FullName,
			&item.RoleID,
			&item.RoleName,
			&item.IsActive,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if foreach != err {
			return nil, err
		}

		Users = append(Users, item)
	}

	return Users, nil

}

func GetUsersByIdRepository (UserID string) (models.User, error) {

	var User models.User
	
	err := databases.DatabaseQuery.QueryRow(`
		SELECT 
			u.id, 
			u.username, 
			u.email, 
			u.full_name, 
			u.role_id, 
			r.name,
			u.is_active,
			u.created_at,
			u.updated_at 
		FROM users as u
		JOIN roles as r on r.id = u.role_id
		WHERE u.id = $1
		`, UserID).Scan(
			&User.ID,
			&User.Username,
			&User.Email,
			&User.FullName,
			&User.RoleID,
			&User.RoleName,
			&User.IsActive,
			&User.CreatedAt,
			&User.UpdatedAt,	
		)
	 if err != nil {
	        return models.User{}, err
	 }
	
	return User, nil

}


func StoreUserRepository(request models.UserRequest)(*models.User, error){
	var User models.User
	var roleID string
	
	defaultPassword := "password"
	
	passwordHash, err := utils.HashPassword(defaultPassword)
	if err != nil {
		return nil, err
	}
	
	// cek role di table
	err = databases.DatabaseQuery.QueryRow(`
		SELECT id
		FROM roles
		WHERE name = 'mahasiswa'
	`).Scan(&roleID)
	
	_, err = databases.DatabaseQuery.Exec(`
		INSERT INTO users (username, email, full_name, password_hash, role_id)
		VALUES ($1, $2, $3, $4, $5)
	`, 
		request.Username,
		request.Email,
		request.FullName,
		passwordHash,
		roleID,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &User, err
	
}
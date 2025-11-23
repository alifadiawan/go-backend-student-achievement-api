package postgres

import (
    "time"

    "github.com/google/uuid"
)

const (
    RoleAdmin    = "Admin"
    RoleMahasiswa  = "Mahasiswa"
    RoleDosen = "Dosen Wali"
)

type Role struct {
    ID          uuid.UUID `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
}
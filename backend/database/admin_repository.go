package database

import (
	"database/sql"
)

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository() *AdminRepository {
	return &AdminRepository{db: DB} // DB — твой глобальный *sql.DB
}

func (r *AdminRepository) IsAdmin(username string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(1) FROM admins WHERE username = ?`, username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

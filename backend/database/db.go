package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"renault-backend/models"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB инициализирует SQLite базу данных
func InitDB() error {
	// Создаем директорию для базы данных, если её нет
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("error creating data directory: %v", err)
	}

	// Открываем базу данных
	dbPath := filepath.Join(dataDir, "renault.db")
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	// Проверяем подключение
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	log.Printf("Successfully connected to SQLite database: %s", dbPath)

	// Создаем таблицу пользователей
	err = createUsersTable()
	if err != nil {
		return err
	}

	// Создаем таблицы для автомобилей
	err = createCarsTables()
	if err != nil {
		return err
	}

	if err := createAdminsTable(); err != nil {
		return err
	}

	if err := createCartTables(); err != nil {
		return err
	}

	// Заполняем данными автомобилей
	err = SeedCarsData()
	if err != nil {
		return err
	}

	return nil
}

func createCartTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS cart_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    car_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1
)`
	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating admins table: %v", err)
	}
	log.Println("Admins table created or already exists")
	return nil
}

func createAdminsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS admins (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL
	)`
	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating admins table: %v", err)
	}
	log.Println("Admins table created or already exists")
	return nil
}

func createUsersTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		IsAdmin bool DEFAULT FALSE
	)
	`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating users table: %v", err)
	}

	log.Println("Users table created or already exists")
	return nil
}

// UserRepository для работы с пользователями
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: DB}
}

func (r *UserRepository) CreateUser(username, email, password string, isAdmin bool) error {
	query := `INSERT INTO users (username, email, password, IsAdmin) VALUES (?, ?, ?, ?)`

	adminValue := 0
	if isAdmin {
		adminValue = 1
	}

	_, err := r.db.Exec(query, username, email, password, adminValue)
	return err
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, email, password, IsAdmin, created_at FROM users WHERE username = ?`
	row := r.db.QueryRow(query, username)

	var user models.User
	var isAdminInt int

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &isAdminInt, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	user.IsAdmin = isAdminInt == 1
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, username, email, password, created_at FROM users WHERE email = ?`
	row := r.db.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// GetAllUsers возвращает всех пользователей (для отладки)
func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	query := `SELECT id, username, email, created_at FROM users ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func createCarsTables() error {
	// Таблица автомобилей
	query := `
    CREATE TABLE IF NOT EXISTS cars (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        model TEXT UNIQUE NOT NULL,
        title TEXT NOT NULL,
        price TEXT NOT NULL,
        image TEXT NOT NULL,
        description TEXT NOT NULL,
        category TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating cars table: %v", err)
	}

	// Таблица технических характеристик
	query = `
    CREATE TABLE IF NOT EXISTS car_specs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        car_id INTEGER NOT NULL,
        name TEXT NOT NULL,
        value TEXT NOT NULL,
        spec_type TEXT NOT NULL,
        FOREIGN KEY (car_id) REFERENCES cars (id) ON DELETE CASCADE
    )`

	_, err = DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating car_specs table: %v", err)
	}

	// Таблица комплектаций
	query = `
    CREATE TABLE IF NOT EXISTS car_equipment (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        car_id INTEGER NOT NULL,
        name TEXT NOT NULL,
        value TEXT NOT NULL,
        FOREIGN KEY (car_id) REFERENCES cars (id) ON DELETE CASCADE
    )`

	_, err = DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating car_equipment table: %v", err)
	}

	// Таблица особенностей
	query = `
    CREATE TABLE IF NOT EXISTS car_features (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        car_id INTEGER NOT NULL,
        feature TEXT NOT NULL,
        FOREIGN KEY (car_id) REFERENCES cars (id) ON DELETE CASCADE
    )`

	_, err = DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating car_features table: %v", err)
	}

	log.Println("Cars tables created or already exist")
	return nil
}

// DeleteUser удаляет пользователя (для отладки)
func (r *UserRepository) DeleteUser(username string) error {
	query := `DELETE FROM users WHERE username = ?`
	_, err := r.db.Exec(query, username)
	return err
}

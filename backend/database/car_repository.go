package database

import (
	"database/sql"
	"renault-backend/models"
)

type CarRepository struct {
	db *sql.DB
}

func NewCarRepository() *CarRepository {
	return &CarRepository{db: DB}
}

// CreateCar создает автомобиль со всеми деталями
func (r *CarRepository) CreateCar(car *models.Car, details *models.CarDetails) error {
	// Начинаем транзакцию
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Вставляем основной автомобиль
	query := `INSERT INTO cars (model, title, price, category, image, description) 
              VALUES (?, ?, ?, ?, ?, ?)`

	result, err := tx.Exec(query, car.Model, car.Title, car.Price,
		car.Category, car.Image, car.Description)
	if err != nil {
		return err
	}

	// Получаем ID созданного автомобиля
	carID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Добавляем технические характеристики
	for _, spec := range details.TechSpecs {
		query = `INSERT INTO car_specs (car_id, name, value, spec_type) 
                 VALUES (?, ?, ?, 'tech')`
		_, err = tx.Exec(query, carID, spec.Name, spec.Value)
		if err != nil {
			return err
		}
	}

	// Добавляем комплектацию
	for _, eq := range details.Equipment {
		query = `INSERT INTO car_specs (car_id, name, value, spec_type) 
                 VALUES (?, ?, ?, 'equipment')`
		_, err = tx.Exec(query, carID, eq.Name, eq.Value)
		if err != nil {
			return err
		}
	}

	// Добавляем особенности
	for _, feature := range details.Features {
		query = `INSERT INTO car_features (car_id, feature) VALUES (?, ?)`
		_, err = tx.Exec(query, carID, feature)
		if err != nil {
			return err
		}
	}

	// Коммитим транзакцию
	return tx.Commit()
}

// GetCarWithDetails возвращает автомобиль со всеми деталями
func (r *CarRepository) GetCarWithDetails(model string) (*models.Car, *models.CarDetails, error) {
	// Получаем основной автомобиль
	query := `SELECT id, model, title, price, category, image, description, created_at 
              FROM cars WHERE model = ?`

	var car models.Car
	err := r.db.QueryRow(query, model).Scan(
		&car.ID, &car.Model, &car.Title, &car.Price, &car.Category,
		&car.Image, &car.Description, &car.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	// Получаем технические характеристики
	var details models.CarDetails

	query = `SELECT name, value FROM car_specs 
             WHERE car_id = ? AND spec_type = 'tech'`
	rows, err := r.db.Query(query, car.ID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var spec models.CarSpec
		if err := rows.Scan(&spec.Name, &spec.Value); err != nil {
			return nil, nil, err
		}
		details.TechSpecs = append(details.TechSpecs, spec)
	}

	// Получаем комплектацию
	query = `SELECT name, value FROM car_specs 
             WHERE car_id = ? AND spec_type = 'equipment'`
	rows, err = r.db.Query(query, car.ID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var spec models.CarSpec
		if err := rows.Scan(&spec.Name, &spec.Value); err != nil {
			return nil, nil, err
		}
		details.Equipment = append(details.Equipment, spec)
	}

	// Получаем особенности
	query = `SELECT feature FROM car_features WHERE car_id = ?`
	rows, err = r.db.Query(query, car.ID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var feature string
		if err := rows.Scan(&feature); err != nil {
			return nil, nil, err
		}
		details.Features = append(details.Features, feature)
	}

	return &car, &details, nil
}

// GetAllCars возвращает все автомобили
func (r *CarRepository) GetAllCars() ([]models.Car, error) {
	query := `SELECT id, model, title, price, category, image, description, created_at 
              FROM cars ORDER BY category, title`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []models.Car
	for rows.Next() {
		var car models.Car
		err := rows.Scan(
			&car.ID, &car.Model, &car.Title, &car.Price, &car.Category,
			&car.Image, &car.Description, &car.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}

	return cars, nil
}

func (r *CarRepository) GetCarByID(id int) (*models.Car, error) {
	row := r.db.QueryRow(`
        SELECT id, title, image, price, description, category
        FROM cars
        WHERE id = ?`,
		id,
	)

	var car models.Car
	err := row.Scan(
		&car.ID,
		&car.Title,
		&car.Image,
		&car.Price,
		&car.Description,
		&car.Category,
	)
	if err != nil {
		return nil, err
	}

	return &car, nil
}

// GetCarsByCategory возвращает автомобили по категории
func (r *CarRepository) GetCarsByCategory(category string) ([]models.Car, error) {
	query := `SELECT id, model, title, price, category, image, description, created_at 
              FROM cars WHERE category = ? ORDER BY title`

	rows, err := r.db.Query(query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []models.Car
	for rows.Next() {
		var car models.Car
		err := rows.Scan(
			&car.ID, &car.Model, &car.Title, &car.Price, &car.Category,
			&car.Image, &car.Description, &car.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}

	return cars, nil
}

// GetCategories возвращает список всех категорий
func (r *CarRepository) GetCategories() ([]string, error) {
	query := `SELECT DISTINCT category FROM cars ORDER BY category`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

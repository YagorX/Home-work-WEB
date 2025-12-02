package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"renault-backend/models"
)

type CarRepository struct {
	db *sql.DB
}

func NewCarRepository() *CarRepository {
	return &CarRepository{db: DB}
}

// CreateCar создает новую запись об автомобиле
func (r *CarRepository) CreateCar(car *models.Car, details *models.CarDetails) error {
	// Преобразуем массивы в JSON
	techSpecsJSON, err := json.Marshal(details.TechSpecs)
	if err != nil {
		return fmt.Errorf("error marshaling tech specs: %v", err)
	}

	equipmentJSON, err := json.Marshal(details.Equipment)
	if err != nil {
		return fmt.Errorf("error marshaling equipment: %v", err)
	}

	featuresJSON, err := json.Marshal(details.Features)
	if err != nil {
		return fmt.Errorf("error marshaling features: %v", err)
	}

	query := `
    INSERT INTO cars (model, title, price, category, image, description, 
                     tech_specs_json, equipment_json, features_json)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

	_, err = r.db.Exec(query, car.Model, car.Title, car.Price, car.Category,
		car.Image, car.Description, techSpecsJSON, equipmentJSON, featuresJSON)

	return err
}

// GetCarByModel возвращает автомобиль по модели
func (r *CarRepository) GetCarByModel(model string) (*models.Car, *models.CarDetails, error) {
	query := `
    SELECT id, model, title, price, category, image, description,
           tech_specs_json, equipment_json, features_json, created_at
    FROM cars WHERE model = ?
    `

	row := r.db.QueryRow(query, model)

	var car models.Car
	var techSpecsJSON, equipmentJSON, featuresJSON string
	var details models.CarDetails

	err := row.Scan(&car.ID, &car.Model, &car.Title, &car.Price, &car.Category,
		&car.Image, &car.Description, &techSpecsJSON, &equipmentJSON,
		&featuresJSON, &car.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	// Преобразуем JSON обратно в структуры
	if err := json.Unmarshal([]byte(techSpecsJSON), &details.TechSpecs); err != nil {
		return nil, nil, fmt.Errorf("error unmarshaling tech specs: %v", err)
	}

	if err := json.Unmarshal([]byte(equipmentJSON), &details.Equipment); err != nil {
		return nil, nil, fmt.Errorf("error unmarshaling equipment: %v", err)
	}

	if err := json.Unmarshal([]byte(featuresJSON), &details.Features); err != nil {
		return nil, nil, fmt.Errorf("error unmarshaling features: %v", err)
	}

	return &car, &details, nil
}

// GetAllCars возвращает все автомобили
func (r *CarRepository) GetAllCars() ([]models.Car, error) {
	query := `SELECT id, model, title, price, category, image, created_at FROM cars ORDER BY category, title`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []models.Car
	for rows.Next() {
		var car models.Car
		err := rows.Scan(&car.ID, &car.Model, &car.Title, &car.Price,
			&car.Category, &car.Image, &car.CreatedAt)
		if err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}

	return cars, nil
}

// GetCarsByCategory возвращает автомобили по категории
func (r *CarRepository) GetCarsByCategory(category string) ([]models.Car, error) {
	query := `SELECT id, model, title, price, category, image, created_at 
              FROM cars WHERE category = ? ORDER BY title`
	rows, err := r.db.Query(query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []models.Car
	for rows.Next() {
		var car models.Car
		err := rows.Scan(&car.ID, &car.Model, &car.Title, &car.Price,
			&car.Category, &car.Image, &car.CreatedAt)
		if err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}

	return cars, nil
}

// DeleteCar удаляет автомобиль
func (r *CarRepository) DeleteCar(model string) error {
	query := `DELETE FROM cars WHERE model = ?`
	_, err := r.db.Exec(query, model)
	return err
}

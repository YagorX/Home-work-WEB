package handlers

import (
	"encoding/json"
	"net/http"
	"renault-backend/database"
)

type CarHandler struct {
	repo *database.CarRepository
}

func NewCarHandler() *CarHandler {
	return &CarHandler{
		repo: database.NewCarRepository(),
	}
}

// GetAllCars возвращает все автомобили
func (h *CarHandler) GetAllCars(w http.ResponseWriter, r *http.Request) {
	cars, err := h.repo.GetAllCars()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

// GetCarByModel возвращает подробную информацию об автомобиле
func (h *CarHandler) GetCarByModel(w http.ResponseWriter, r *http.Request) {
	model := r.URL.Query().Get("model")
	if model == "" {
		http.Error(w, "Model parameter is required", http.StatusBadRequest)
		return
	}

	car, details, err := h.repo.GetCarByModel(model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if car == nil {
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"car":     car,
		"details": details,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCarsByCategory возвращает автомобили по категории
func (h *CarHandler) GetCarsByCategory(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	if category == "" {
		http.Error(w, "Category parameter is required", http.StatusBadRequest)
		return
	}

	cars, err := h.repo.GetCarsByCategory(category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

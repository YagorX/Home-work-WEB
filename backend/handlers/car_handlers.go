package handlers

import (
	"encoding/json"
	"net/http"
	"renault-backend/database"

	"github.com/gorilla/mux"
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
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, cars)
}

// GetCarByModel возвращает автомобиль по модели
func (h *CarHandler) GetCarByModel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	model := vars["model"]

	car, details, err := h.repo.GetCarWithDetails(model)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if car == nil {
		respondWithError(w, http.StatusNotFound, "Car not found")
		return
	}

	response := map[string]interface{}{
		"id":          car.ID,
		"model":       car.Model,
		"title":       car.Title,
		"price":       car.Price,
		"category":    car.Category,
		"image":       car.Image,
		"description": car.Description,
		"techSpecs":   details.TechSpecs,
		"equipment":   details.Equipment,
		"features":    details.Features,
	}

	respondWithJSON(w, http.StatusOK, response)
}

// GetCarsByCategory возвращает автомобили по категории
func (h *CarHandler) GetCarsByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]

	cars, err := h.repo.GetCarsByCategory(category)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, cars)
}

// GetCategories возвращает все категории
func (h *CarHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.repo.GetCategories()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, categories)
}

// Вспомогательные функции для работы с JSON

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error marshaling response"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

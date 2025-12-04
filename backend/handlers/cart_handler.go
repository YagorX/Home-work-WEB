package handlers

import (
	"encoding/json"
	"net/http"
	"renault-backend/database"
	"strconv"
	"strings"
)

type CartHandler struct {
	cartRepo *database.CartRepository
	carRepo  *database.CarRepository
}

func NewCartHandler() *CartHandler {
	return &CartHandler{
		cartRepo: database.NewCartRepository(),
		carRepo:  database.NewCarRepository(),
	}
}

// =============== УТИЛИТА: получить userID из заголовка ===============
func getUserID(r *http.Request) (int, error) {
	header := r.Header.Get("X-User-Id")
	return strconv.Atoi(header)
}

// ========================= GET /api/cart =============================
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "NO USER ID", http.StatusUnauthorized)
		return
	}

	items, err := h.cartRepo.GetCart(userID)
	if err != nil {
		http.Error(w, "DB ERROR", http.StatusInternalServerError)
		return
	}

	// ДЕЛАЕМ НЕ nil, а пустой слайс
	response := make([]database.CartItemResponse, 0)

	for _, item := range items {
		car, err := h.carRepo.GetCarByID(item.CarID)
		if err != nil {
			continue
		}

		response = append(response, database.CartItemResponse{
			ID:       item.ID,
			CarID:    car.ID,
			Title:    car.Title,
			Image:    car.Image,
			Price:    car.Price,
			Quantity: item.Quantity,
		})
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response) // [] даже если пустой
}

// ========================= POST /api/cart =============================
type AddToCartDTO struct {
	CarID    int `json:"carId"`
	Quantity int `json:"quantity"`
}

func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "NO USER ID", http.StatusUnauthorized)
		return
	}

	var req AddToCartDTO
	json.NewDecoder(r.Body).Decode(&req)

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	// Проверяем, есть ли товар в корзине
	existing, err := h.cartRepo.FindItem(userID, req.CarID)
	if err != nil {
		http.Error(w, "DB ERROR", http.StatusInternalServerError)
		return
	}

	if existing != nil {
		// обновляем количество
		h.cartRepo.UpdateQuantity(existing.ID, existing.Quantity+req.Quantity)
	} else {
		// добавляем новую запись
		h.cartRepo.AddItem(userID, req.CarID, req.Quantity)
	}

	w.WriteHeader(http.StatusOK)
}

// ================= PATCH /api/cart/{id} =============================
type UpdateQuantityDTO struct {
	Quantity int `json:"quantity"`
}

func (h *CartHandler) UpdateQuantity(w http.ResponseWriter, r *http.Request) {
	_, err := getUserID(r)
	if err != nil {
		http.Error(w, "NO USER ID", http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/cart/")
	itemID, _ := strconv.Atoi(idStr)

	var req UpdateQuantityDTO
	json.NewDecoder(r.Body).Decode(&req)

	if req.Quantity <= 0 {
		h.cartRepo.DeleteItem(itemID)
		w.WriteHeader(http.StatusOK)
		return
	}

	h.cartRepo.UpdateQuantity(itemID, req.Quantity)
	w.WriteHeader(http.StatusOK)
}

// ================= DELETE /api/cart/{id} =============================
func (h *CartHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/cart/")
	itemID, _ := strconv.Atoi(idStr)

	h.cartRepo.DeleteItem(itemID)
	w.WriteHeader(http.StatusOK)
}

// ================= DELETE /api/cart =============================
func (h *CartHandler) Clear(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "NO USER ID", http.StatusUnauthorized)
		return
	}

	h.cartRepo.ClearCart(userID)
	w.WriteHeader(http.StatusOK)
}

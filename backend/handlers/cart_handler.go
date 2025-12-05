package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"renault-backend/database"

	"github.com/gorilla/mux"
)

type CartHandler struct {
	db *sql.DB
}

func NewCartHandler() *CartHandler {
	h := &CartHandler{db: database.DB}
	if err := h.initCartTable(); err != nil {
		fmt.Println("initCartTable error:", err)
	}
	return h
}

// создаём таблицу, если её ещё нет
func (h *CartHandler) initCartTable() error {
	_, err := h.db.Exec(`
        CREATE TABLE IF NOT EXISTS cart_items (
            id        INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id   TEXT    NOT NULL,
            car_id    TEXT    NOT NULL,
            quantity  INTEGER NOT NULL DEFAULT 1,
            UNIQUE(user_id, car_id)
        );
    `)
	return err
}

type CartItem struct {
	ID       int    `json:"id"`
	UserID   string `json:"userId"`
	CarID    string `json:"carId"`
	Quantity int    `json:"quantity"`
}

type addToCartRequest struct {
	CarID    string `json:"carId"`
	Quantity int    `json:"quantity"`
}

type updateQuantityRequest struct {
	Quantity int `json:"quantity"`
}

func getUserID(r *http.Request) string {
	return r.Header.Get("X-User-Id")
}

// ---------- GET /api/cart ----------
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		http.Error(w, "missing X-User-Id", http.StatusBadRequest)
		return
	}

	rows, err := h.db.Query(`
        SELECT id, user_id, car_id, quantity
        FROM cart_items
        WHERE user_id = ?
        ORDER BY id
    `, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("db error GetCart: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []CartItem
	for rows.Next() {
		var it CartItem
		if err := rows.Scan(&it.ID, &it.UserID, &it.CarID, &it.Quantity); err != nil {
			http.Error(w, fmt.Sprintf("scan error GetCart: %v", err), http.StatusInternalServerError)
			return
		}
		items = append(items, it)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(items)
}

// ---------- POST /api/cart ----------
func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		http.Error(w, "missing X-User-Id", http.StatusBadRequest)
		return
	}

	var req addToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("invalid json: %v", err), http.StatusBadRequest)
		return
	}
	if req.CarID == "" {
		http.Error(w, "carId is required", http.StatusBadRequest)
		return
	}
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	// upsert без ON CONFLICT: сначала UPDATE, если не затронуло строк — INSERT
	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, fmt.Sprintf("db begin error: %v", err), http.StatusInternalServerError)
		return
	}

	// пробуем увеличить quantity
	res, err := tx.Exec(`
        UPDATE cart_items
        SET quantity = quantity + ?
        WHERE user_id = ? AND car_id = ?
    `, req.Quantity, userID, req.CarID)
	if err != nil {
		tx.Rollback()
		http.Error(w, fmt.Sprintf("db update error: %v", err), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		http.Error(w, fmt.Sprintf("rows affected error: %v", err), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		// записи не было — вставляем новую
		_, err = tx.Exec(`
            INSERT INTO cart_items (user_id, car_id, quantity)
            VALUES (?, ?, ?)
        `, userID, req.CarID, req.Quantity)
		if err != nil {
			tx.Rollback()
			http.Error(w, fmt.Sprintf("db insert error: %v", err), http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, fmt.Sprintf("db commit error: %v", err), http.StatusInternalServerError)
		return
	}

	// вернём актуальную корзину
	h.GetCart(w, r)
}

// ---------- PATCH /api/cart/{id} ----------
func (h *CartHandler) UpdateQuantity(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		http.Error(w, "missing X-User-Id", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	carID := vars["id"]
	if carID == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	var req updateQuantityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("invalid json: %v", err), http.StatusBadRequest)
		return
	}

	if req.Quantity <= 0 {
		_, err := h.db.Exec(`
            DELETE FROM cart_items
            WHERE user_id = ? AND car_id = ?
        `, userID, carID)
		if err != nil {
			http.Error(w, fmt.Sprintf("db delete error: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		_, err := h.db.Exec(`
            UPDATE cart_items
            SET quantity = ?
            WHERE user_id = ? AND car_id = ?
        `, req.Quantity, userID, carID)
		if err != nil {
			http.Error(w, fmt.Sprintf("db update error: %v", err), http.StatusInternalServerError)
			return
		}
	}

	h.GetCart(w, r)
}

// ---------- DELETE /api/cart/{id} ----------
func (h *CartHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		http.Error(w, "missing X-User-Id", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	carID := vars["id"]
	if carID == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	_, err := h.db.Exec(`
        DELETE FROM cart_items
        WHERE user_id = ? AND car_id = ?
    `, userID, carID)
	if err != nil {
		http.Error(w, fmt.Sprintf("db delete error: %v", err), http.StatusInternalServerError)
		return
	}

	h.GetCart(w, r)
}

// ---------- DELETE /api/cart ----------
func (h *CartHandler) Clear(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		http.Error(w, "missing X-User-Id", http.StatusBadRequest)
		return
	}

	_, err := h.db.Exec(`
        DELETE FROM cart_items
        WHERE user_id = ?
    `, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("db clear error: %v", err), http.StatusInternalServerError)
		return
	}

	h.GetCart(w, r)
}

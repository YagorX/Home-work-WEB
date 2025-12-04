package database

import "database/sql"

type CartItem struct {
	ID       int `json:"id"`
	UserID   int `json:"user_id"`
	CarID    int `json:"car_id"`
	Quantity int `json:"quantity"`
}

type CartItemResponse struct {
	ID       int    `json:"id"`
	CarID    int    `json:"car_id"`
	Title    string `json:"title"`
	Image    string `json:"image"`
	Price    string `json:"price"`
	Quantity int    `json:"quantity"`
}

type CartRepository struct {
	db *sql.DB
}

func NewCartRepository() *CartRepository {
	return &CartRepository{db: DB}
}

// Получить корзину пользователя
func (r *CartRepository) GetCart(userID int) ([]CartItem, error) {
	rows, err := r.db.Query(`
        SELECT id, user_id, car_id, quantity 
        FROM cart_items 
        WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CartItem
	for rows.Next() {
		var item CartItem
		if err := rows.Scan(&item.ID, &item.UserID, &item.CarID, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// Найти запись по userId + carId
func (r *CartRepository) FindItem(userID, carID int) (*CartItem, error) {
	row := r.db.QueryRow(`
        SELECT id, user_id, car_id, quantity 
        FROM cart_items 
        WHERE user_id = ? AND car_id = ?`,
		userID, carID)

	var item CartItem
	err := row.Scan(&item.ID, &item.UserID, &item.CarID, &item.Quantity)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &item, err
}

// Добавить новую запись
func (r *CartRepository) AddItem(userID, carID, quantity int) error {
	_, err := r.db.Exec(`
        INSERT INTO cart_items (user_id, car_id, quantity) 
        VALUES (?, ?, ?)`,
		userID, carID, quantity)
	return err
}

// Обновить количество
func (r *CartRepository) UpdateQuantity(id, quantity int) error {
	_, err := r.db.Exec(`
        UPDATE cart_items SET quantity = ? WHERE id = ?`,
		quantity, id)
	return err
}

// Удалить 1 запись
func (r *CartRepository) DeleteItem(id int) error {
	_, err := r.db.Exec(`DELETE FROM cart_items WHERE id = ?`, id)
	return err
}

// Очистить корзину
func (r *CartRepository) ClearCart(userID int) error {
	_, err := r.db.Exec(`DELETE FROM cart_items WHERE user_id = ?`, userID)
	return err
}

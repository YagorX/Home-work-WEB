// backend/main.go
package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"renault-backend/database"
	"renault-backend/handlers"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

const (
	JWT_SECRET = "your_very_strong_jwt_secret_key_change_this_in_production_123!"
	PORT       = "8080"
)

// ----- Модели каталога -----

type Car struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Model       string   `json:"model"` // дублируем title, чтобы фронт не ломался
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Image       string   `json:"image"`  // главное изображение (превью)
	Images      []string `json:"images"` // все изображения для галереи
	Price       int      `json:"price"`
	Features    []string `json:"features"`
	TechSpecs   []Spec   `json:"techSpecs"`
	Equipment   []Spec   `json:"equipment"`
}

type Spec struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// отдельная БД под каталог автомобилей
var carDB *sql.DB

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-User-Id")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// preflight OPTIONS request — отвечаем сразу
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// ---------- БД пользователей / auth (твоя старая логика) ----------
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to connect to users DB: %v", err)
	}
	defer database.DB.Close()

	userRepo := database.NewUserRepository()
	authHandler := handlers.NewAuthHandler(userRepo, JWT_SECRET)

	// ---------- БД каталога автомобилей ----------
	var err error
	// путь поправь, если бинарник запускается не из корня проекта
	carDB, err = sql.Open("sqlite3", "cars.db")
	if err != nil {
		log.Fatal("Ошибка открытия БД каталога:", err)
	}
	defer carDB.Close()

	if err := carDB.Ping(); err != nil {
		log.Fatal("Ошибка подключения к БД каталога:", err)
	}

	if err := createCarTables(); err != nil {
		log.Fatal("Ошибка создания таблиц каталога:", err)
	}

	if err := seedCarData(); err != nil {
		log.Fatal("Ошибка начального заполнения каталога:", err)
	}

	// ---------- Роутер ----------
	router := mux.NewRouter()

	// подроутер /api
	api := router.PathPrefix("/api").Subrouter()

	// Публичные маршруты (auth и прочее)
	api.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	api.HandleFunc("/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/login", authHandler.Login).Methods("POST")
	api.HandleFunc("/validate-password", authHandler.ValidatePassword).Methods("POST")
	api.HandleFunc("/password-rules", authHandler.PasswordRules).Methods("GET")

	// Отладочные маршруты (как было)
	api.HandleFunc("/users", authHandler.GetAllUsers).Methods("GET")

	// Каталог автомобилей — новые хендлеры на carDB
	api.HandleFunc("/cars", getAllCarsHandler).Methods("GET")
	api.HandleFunc("/cars/{id}", getCarByIDHandler).Methods("GET")
	// если нужно будет фильтровать по категории:
	// api.HandleFunc("/cars/category/{category}", getCarsByCategoryHandler).Methods("GET")

	// Каталог автомобилей — публичные GET
	api.HandleFunc("/cars", getAllCarsHandler).Methods("GET")
	api.HandleFunc("/cars/{id}", getCarByIDHandler).Methods("GET")
	// api.HandleFunc("/cars/category/{category}", getCarsByCategoryHandler).Methods("GET")

	// ----- АДМИНСКИЕ РОУТЫ ДЛЯ КАТАЛОГА -----
	admin := api.PathPrefix("/admin").Subrouter()

	// защищаем все маршруты /api/admin/...
	admin.Use(JWTAdminMiddleware)

	admin.HandleFunc("/cars", createCarHandler).Methods("POST")
	admin.HandleFunc("/cars/{id}", updateCarHandler).Methods("PUT")
	admin.HandleFunc("/cars/{id}", deleteCarHandler).Methods("DELETE")

	cartHandler := handlers.NewCartHandler()

	// /api/cart (список, добавление, очистка)
	api.HandleFunc("/cart", cartHandler.GetCart).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/cart", cartHandler.AddToCart).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/cart", cartHandler.Clear).Methods(http.MethodDelete, http.MethodOptions)

	// /api/cart/{id} (обновление количества, удаление позиции)
	api.HandleFunc("/cart/{id}", cartHandler.UpdateQuantity).Methods(http.MethodPatch, http.MethodOptions)
	api.HandleFunc("/cart/{id}", cartHandler.DeleteItem).Methods(http.MethodDelete, http.MethodOptions)

	// ---------- CORS ----------
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // пока можно так, потом ограничишь
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With", "X-User-Id"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400,
	})

	addr := ":" + PORT
	log.Printf("🚗 Renault Backend Server starting on http://localhost%s", addr)
	log.Printf("📡 API endpoints:")
	log.Printf("  🔍 GET  http://localhost%s/api/health", addr)
	log.Printf("  📝 POST http://localhost%s/api/register", addr)
	log.Printf("  🔑 POST http://localhost%s/api/login", addr)
	log.Printf("  📊 POST http://localhost%s/api/validate-password", addr)
	log.Printf("  📋 GET  http://localhost%s/api/password-rules", addr)
	log.Printf("  👥 GET  http://localhost%s/api/users", addr)
	log.Println("")
	log.Println("🔒 Правила паролей:")
	log.Println("  - Минимум 8 символов")
	log.Println("  - Хотя бы одна заглавная и строчная буква")
	log.Println("  - Хотя бы одна цифра")
	log.Println("  - Хотя бы один специальный символ")
	log.Println("  - Запрещены простые пароли и последовательности")

	handler := corsHandler.Handler(router)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// ---------- Работа с БД каталога ----------

func createCarTables() error {
	_, err := carDB.Exec(`
        CREATE TABLE IF NOT EXISTS cars (
            id TEXT PRIMARY KEY,
            title TEXT NOT NULL,
            description TEXT,
            category TEXT,
            image TEXT,
            base_price INTEGER
        );

        CREATE TABLE IF NOT EXISTS car_features (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            car_id TEXT NOT NULL,
            name TEXT NOT NULL,
            FOREIGN KEY (car_id) REFERENCES cars(id) ON DELETE CASCADE
        );

        CREATE TABLE IF NOT EXISTS car_images (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            car_id TEXT NOT NULL,
            image_path TEXT NOT NULL,
            FOREIGN KEY (car_id) REFERENCES cars(id) ON DELETE CASCADE
        );

    `)

	return err
}

func createCarHandler(w http.ResponseWriter, r *http.Request) {
	var c Car
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// ID либо приходит с фронта, либо можно сгенерировать (slug/uuid)
	if c.ID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	tx, err := carDB.Begin()
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// вставляем запись в cars
	_, err = tx.Exec(`
        INSERT INTO cars (id, title, description, category, image, base_price)
        VALUES (?, ?, ?, ?, ?, ?)
    `, c.ID, c.Title, c.Description, c.Category, c.Image, c.Price)
	if err != nil {
		http.Error(w, "db error: insert car", http.StatusInternalServerError)
		return
	}

	// features
	if len(c.Features) > 0 {
		stmt, err := tx.Prepare(`INSERT INTO car_features (car_id, name) VALUES (?, ?)`)
		if err != nil {
			http.Error(w, "db error: prepare features", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		for _, f := range c.Features {
			if _, err := stmt.Exec(c.ID, f); err != nil {
				http.Error(w, "db error: insert feature", http.StatusInternalServerError)
				return
			}
		}
	}

	// images
	if len(c.Images) > 0 {
		stmt, err := tx.Prepare(`INSERT INTO car_images (car_id, image_path) VALUES (?, ?)`)
		if err != nil {
			http.Error(w, "db error: prepare images", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		for _, img := range c.Images {
			if _, err := stmt.Exec(c.ID, img); err != nil {
				http.Error(w, "db error: insert image", http.StatusInternalServerError)
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "db error: commit", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "created", "id": c.ID})
}

func updateCarHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var c Car
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// на всякий случай принудительно проставим id
	c.ID = id

	tx, err := carDB.Begin()
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// обновляем основную запись
	_, err = tx.Exec(`
        UPDATE cars
        SET title = ?, description = ?, category = ?, image = ?, base_price = ?
        WHERE id = ?
    `, c.Title, c.Description, c.Category, c.Image, c.Price, c.ID)
	if err != nil {
		http.Error(w, "db error: update car", http.StatusInternalServerError)
		return
	}

	// проще всего – удалить старые features/images и записать новые
	if _, err := tx.Exec(`DELETE FROM car_features WHERE car_id = ?`, c.ID); err != nil {
		http.Error(w, "db error: clear features", http.StatusInternalServerError)
		return
	}
	if _, err := tx.Exec(`DELETE FROM car_images WHERE car_id = ?`, c.ID); err != nil {
		http.Error(w, "db error: clear images", http.StatusInternalServerError)
		return
	}

	if len(c.Features) > 0 {
		stmt, err := tx.Prepare(`INSERT INTO car_features (car_id, name) VALUES (?, ?)`)
		if err != nil {
			http.Error(w, "db error: prepare features", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()
		for _, f := range c.Features {
			if _, err := stmt.Exec(c.ID, f); err != nil {
				http.Error(w, "db error: insert feature", http.StatusInternalServerError)
				return
			}
		}
	}

	if len(c.Images) > 0 {
		stmt, err := tx.Prepare(`INSERT INTO car_images (car_id, image_path) VALUES (?, ?)`)
		if err != nil {
			http.Error(w, "db error: prepare images", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()
		for _, img := range c.Images {
			if _, err := stmt.Exec(c.ID, img); err != nil {
				http.Error(w, "db error: insert image", http.StatusInternalServerError)
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "db error: commit", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func deleteCarHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	res, err := carDB.Exec(`DELETE FROM cars WHERE id = ?`, id)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

func seedCarData() error {
	var count int
	if err := carDB.QueryRow("SELECT COUNT(*) FROM cars").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		// уже есть данные — ничего не делаем
		return nil
	}

	tx, err := carDB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// ----- cars -----
	carsInsert := `
		INSERT INTO cars (id, title, description, category, image, base_price)
		VALUES (?, ?, ?, ?, ?, ?);
	`
	cars := []struct {
		ID, Title, Desc, Category, Image string
		Price                            int
	}{
		{"logan", "Renault Logan", "Надежный седан для города и трассы. Идеальное сочетание цены и качества.",
			"Легковые", "images/renault_logan.jpeg", 950000},
		{"sandero", "Renault Sandero", "Компактный хэтчбек с просторным салоном и экономичным двигателем.",
			"Легковые", "images/renault_sander.jpg", 890000},
		{"stepway", "Renault Sandero Stepway", "Хэтчбек в кросс-кузове с увеличенным клиренсом и стильным дизайном.",
			"Легковые", "images/renault_sander_stepway.jpeg", 1100000},

		{"duster", "Renault Duster", "Легендарный внедорожник с полным приводом. Покоритель любых дорог.",
			"Кроссоверы", "images/duster.jpeg", 1450000},
		{"kaptur", "Renault Kaptur", "Стильный компактный кроссовер с передовыми технологиями безопасности.",
			"Кроссоверы", "images/kapture.jpeg", 1350000},
		{"arkana", "Renault Arkana", "Элегантное кросс-купе с динамичным характером и просторным салоном.",
			"Кроссоверы", "images/arkana.jpeg", 1650000},

		{"loganvan", "Renault Logan Van", "Коммерческая версия Logan с увеличенным багажным отделением.",
			"Коммерческие", "images/van.jpeg", 1000000},
		{"kangoo", "Renault Kangoo", "Компактный коммерческий автомобиль с отличной маневренностью.",
			"Коммерческие", "images/kangoo.jpeg", 1300000},
		{"trafic", "Renault Trafic", "Универсальный коммерческий автомобиль для перевозки грузов.",
			"Коммерческие", "images/trafic.jpg", 1800000},

		{"zoe", "Renault ZOE", "Компактный электромобиль для города с впечатляющим запасом хода.",
			"Электромобили", "images/zoe.jpeg", 2200000},
		{"megane", "Renault Megane E-Tech", "Современный электрокроссовер с технологиями нового поколения.",
			"Электромобили", "images/megane e.jpg", 3500000},
		{"captur", "Renault Captur E-Tech", "Гибридный кроссовер с экономичным расходом и отличной динамикой.",
			"Гибриды", "images/captur e.jpg", 1900000},
	}

	for _, c := range cars {
		if _, err := tx.Exec(carsInsert, c.ID, c.Title, c.Desc, c.Category, c.Image, c.Price); err != nil {
			return err
		}
	}

	// ----- car_features -----
	featuresInsert := `INSERT INTO car_features (car_id, name) VALUES (?, ?);`
	features := map[string][]string{
		"logan": {
			"Расход: 6.1 л/100км",
			"Мощность: 82 л.с.",
			"Объем багажника: 510 л",
		},
		"sandero": {
			"Расход: 5.8 л/100км",
			"Мощность: 75 л.с.",
			"5-ступенчатая МКПП",
		},
		"stepway": {
			"Клиренс: 195 мм",
			"Мощность: 90 л.с.",
			"Защита бампера",
		},
		"duster": {
			"Полный привод 4x4",
			"Мощность: 114 л.с.",
			"Клиренс: 210 мм",
		},
		"kaptur": {
			"Система ESP",
			"Мощность: 113 л.с.",
			"Мультимедиа R-Link",
		},
		"arkana": {
			"Купе-форма",
			"Мощность: 150 л.с.",
			"Вариатор X-Tronic",
		},
		"loganvan": {
			"Объем багажника: 800 л",
			"Грузоподъемность: 500 кг",
			"Низкий расход топлива",
		},
		"kangoo": {
			"Объем: 3-4.6 м³",
			"Грузоподъемность: 650 кг",
			"Сдвижные двери",
		},
		"trafic": {
			"Объем: 5.2-8.6 м³",
			"Грузоподъемность: 1-1.5 т",
			"Дизельный двигатель",
		},
		"zoe": {
			"Запас хода: 395 км",
			"Мощность: 135 л.с.",
			"Быстрая зарядка за 30 мин",
		},
		"megane": {
			"Запас хода: 470 км",
			"Мощность: 220 л.с.",
			"Цифровая панель 12,3\"",
		},
		"captur": {
			"Гибридная система",
			"Расход: 4.5 л/100км",
			"Электро-привод на малых скоростях",
		},
	}

	imagesInsert := `INSERT INTO car_images (car_id, image_path) VALUES (?, ?);`

	carImages := map[string][]string{
		"logan": {
			"images/renault_logan.jpeg",
			"images/renault_logan_2.jpg",
			"images/renaul_logan_3.jpg",
		},
		"sandero": {
			"images/renault_sander.jpg",
			"images/renault_sandero2.jpg",
		},
		"stepway": {
			"images/renault_sander_stepway.jpeg",
			"images/renault_sandero_stepway2.jpg",
		},
	}

	for carID, imgs := range carImages {
		for _, path := range imgs {
			if _, err := tx.Exec(imagesInsert, carID, path); err != nil {
				return err
			}
		}
	}

	for carID, list := range features {
		for _, f := range list {
			if _, err := tx.Exec(featuresInsert, carID, f); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func getCarImages(carID string) ([]string, error) {
	rows, err := carDB.Query(`SELECT image_path FROM car_images WHERE car_id = ? ORDER BY id`, carID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []string
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, err
		}
		images = append(images, path)
	}
	return images, nil
}

// ---------- HTTP-хендлеры каталога ----------

func getAllCarsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := carDB.Query(`SELECT id, title, description, category, image, base_price FROM cars`)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cars []Car

	for rows.Next() {
		var c Car
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.Category, &c.Image, &c.Price); err != nil {
			http.Error(w, "scan error", http.StatusInternalServerError)
			return
		}
		c.Model = c.Title

		// подгружаем features
		featRows, err := carDB.Query(`SELECT name FROM car_features WHERE car_id = ?`, c.ID)
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		for featRows.Next() {
			var name string
			if err := featRows.Scan(&name); err != nil {
				http.Error(w, "scan error", http.StatusInternalServerError)
				featRows.Close()
				return
			}
			c.Features = append(c.Features, name)
		}
		featRows.Close()

		// подгружаем изображения
		imgs, err := getCarImages(c.ID)
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		// если в таблице нет записей, хотя бы главное изображение
		if len(imgs) == 0 && c.Image != "" {
			imgs = []string{c.Image}
		}
		c.Images = imgs

		cars = append(cars, c)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(cars)
}

func getCarByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var c Car
	err := carDB.QueryRow(
		`SELECT id, title, description, category, image, base_price FROM cars WHERE id = ?`,
		id,
	).Scan(&c.ID, &c.Title, &c.Description, &c.Category, &c.Image, &c.Price)
	if err == sql.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	c.Model = c.Title

	featRows, err := carDB.Query(`SELECT name FROM car_features WHERE car_id = ?`, c.ID)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	for featRows.Next() {
		var name string
		if err := featRows.Scan(&name); err != nil {
			http.Error(w, "scan error", http.StatusInternalServerError)
			featRows.Close()
			return
		}
		c.Features = append(c.Features, name)
	}
	featRows.Close()

	imgs, err := getCarImages(c.ID)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	if len(imgs) == 0 && c.Image != "" {
		imgs = []string{c.Image}
	}
	c.Images = imgs

	c.TechSpecs = []Spec{}
	c.Equipment = []Spec{}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(c)
}

// опционально — фильтр по категории, если понадобится
func getCarsByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]

	rows, err := carDB.Query(
		`SELECT id, title, description, category, image, base_price FROM cars WHERE category = ?`,
		category,
	)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cars []Car

	for rows.Next() {
		var c Car
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.Category, &c.Image, &c.Price); err != nil {
			http.Error(w, "scan error", http.StatusInternalServerError)
			return
		}
		c.Model = c.Title

		featRows, err := carDB.Query(`SELECT name FROM car_features WHERE car_id = ?`, c.ID)
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		for featRows.Next() {
			var name string
			if err := featRows.Scan(&name); err != nil {
				http.Error(w, "scan error", http.StatusInternalServerError)
				featRows.Close()
				return
			}
			c.Features = append(c.Features, name)
		}
		featRows.Close()

		cars = append(cars, c)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(cars)
}

func JWTAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(JWT_SECRET), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		isAdmin, _ := claims["is_admin"].(bool)
		if !isAdmin {
			http.Error(w, "forbidden: admin only", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

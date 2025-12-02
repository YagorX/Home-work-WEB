package main

import (
	"log"
	"net/http"
	"renault-backend/database"
	"renault-backend/handlers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const (
	JWT_SECRET = "your_very_strong_jwt_secret_key_change_this_in_production_123!"
	PORT       = "8080"
)

func main() {
	// Инициализация базы данных SQLite
	err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.DB.Close()

	// Инициализация репозитория и обработчиков
	userRepo := database.NewUserRepository()
	authHandler := handlers.NewAuthHandler(userRepo, JWT_SECRET)

	// Настройка роутера
	router := mux.NewRouter()

	carHandler := handlers.NewCarHandler()

	// Маршруты API
	api := router.PathPrefix("/api").Subrouter()

	// Публичные маршруты
	api.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	api.HandleFunc("/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/login", authHandler.Login).Methods("POST")
	api.HandleFunc("/validate-password", authHandler.ValidatePassword).Methods("POST")
	api.HandleFunc("/password-rules", authHandler.PasswordRules).Methods("GET")

	// Отладочные маршруты (в продакшене убрать или защитить)
	api.HandleFunc("/users", authHandler.GetAllUsers).Methods("GET")

	api.HandleFunc("/cars", carHandler.GetAllCars).Methods("GET")
	api.HandleFunc("/cars/{model}", carHandler.GetCarByModel).Methods("GET")
	api.HandleFunc("/cars/category/{category}", carHandler.GetCarsByCategory).Methods("GET")

	// Настройка CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // В продакшене заменить на конкретные домены
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400,
	})

	// Запуск сервера
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

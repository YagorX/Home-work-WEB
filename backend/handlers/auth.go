package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"renault-backend/database"
	"renault-backend/models"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type AuthHandler struct {
	userRepo           *database.UserRepository
	jwtSecret          string
	passwordValidation models.PasswordValidation
}

func NewAuthHandler(userRepo *database.UserRepository, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userRepo:           userRepo,
		jwtSecret:          jwtSecret,
		passwordValidation: models.DefaultPasswordValidation,
	}
}

// Register обрабатывает регистрацию пользователя
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Неверный формат данных", http.StatusBadRequest, nil)
		return
	}

	// Валидация данных
	errors := h.validateRegistration(req)
	if len(errors) > 0 {
		sendError(w, "Ошибка валидации", http.StatusBadRequest, errors)
		return
	}

	// Проверка, существует ли пользователь с таким именем
	existingUser, err := h.userRepo.GetUserByUsername(req.Username)
	if err != nil && err != sql.ErrNoRows {
		sendError(w, "Ошибка базы данных", http.StatusInternalServerError, nil)
		return
	}
	if existingUser != nil {
		sendError(w, "Пользователь с таким именем уже существует", http.StatusConflict, nil)
		return
	}

	// Проверка, существует ли пользователь с таким email
	existingUser, err = h.userRepo.GetUserByEmail(req.Email)
	if err != nil && err != sql.ErrNoRows {
		sendError(w, "Ошибка базы данных", http.StatusInternalServerError, nil)
		return
	}
	if existingUser != nil {
		sendError(w, "Пользователь с таким email уже существует", http.StatusConflict, nil)
		return
	}

	// Создаем пользователя (пока без ID)
	var user models.User
	user.Username = strings.TrimSpace(req.Username)
	user.Email = strings.TrimSpace(req.Email)

	// Хешируем пароль
	if err := user.HashPassword(req.Password); err != nil {
		sendError(w, "Ошибка при обработке пароля", http.StatusInternalServerError, nil)
		return
	}

	isAdmin := h.isAdmin(user.Username)

	// Сохраняем в БД с правильным флагом
	if err := h.userRepo.CreateUser(user.Username, user.Email, user.Password, isAdmin); err != nil {
		sendError(w, "Ошибка при создании пользователя", http.StatusInternalServerError, nil)
		return
	}

	println(user.Username, "username")
	// ЕЩЁ РАЗ читаем пользователя из БД, чтобы получить ID
	createdUser, err := h.userRepo.GetUserByUsername(user.Username)
	println(createdUser.Email, "email is register user")
	if err != nil {
		sendError(w, "Ошибка при получении данных пользователя", http.StatusInternalServerError, nil)
		return
	}

	// Генерируем JWT токен c признаком is_admin
	token, err := h.generateToken(createdUser.Username, isAdmin)
	if err != nil {
		sendError(w, "Ошибка при генерации токена", http.StatusInternalServerError, nil)
		return
	}

	response := models.AuthResponse{
		Success: true,
		Message: "Регистрация успешно завершена",
		Token:   token,
		User: &models.User{
			ID:       createdUser.ID,
			Username: createdUser.Username,
			Email:    createdUser.Email,
			IsAdmin:  isAdmin,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) isAdmin(username string) bool {
	// используем глобальную DB из пакета database
	var count int
	err := database.DB.QueryRow(
		"SELECT COUNT(*) FROM admins WHERE username = ?",
		username,
	).Scan(&count)
	println(count, "true or false on table admins")
	if err != nil {
		return false
	}
	return count > 0
}

// Login обрабатывает вход пользователя
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Неверный формат данных", http.StatusBadRequest, nil)
		return
	}

	// Валидация
	validationErrors := h.validateLogin(req)
	if len(validationErrors) > 0 {
		sendError(w, "Ошибка валидации", http.StatusBadRequest, validationErrors)
		return
	}

	// Получаем пользователя из БД
	user, err := h.userRepo.GetUserByUsername(req.Username)
	if err != nil && err != sql.ErrNoRows {
		sendError(w, "Ошибка базы данных", http.StatusInternalServerError, nil)
		return
	}
	if user == nil {
		// Для безопасности не говорим, что пользователь не существует
		sendError(w, "Неверное имя пользователя или пароль", http.StatusUnauthorized, nil)
		return
	}

	// Проверяем пароль
	if err := user.CheckPassword(req.Password); err != nil {
		sendError(w, "Неверное имя пользователя или пароль", http.StatusUnauthorized, nil)
		return
	}

	isAdmin := h.isAdmin(user.Username)

	// Генерируем JWT токен c признаком is_admin
	token, err := h.generateToken(user.Username, isAdmin)
	if err != nil {
		sendError(w, "Ошибка при генерации токена", http.StatusInternalServerError, nil)
		return
	}

	response := models.AuthResponse{
		Success: true,
		Message: "Вход выполнен успешно",
		Token:   token,
		User: &models.User{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			IsAdmin:  isAdmin, // <-- вот тут важно
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ValidatePassword проверяет пароль (публичный endpoint)
func (h *AuthHandler) ValidatePassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Неверный формат данных", http.StatusBadRequest, nil)
		return
	}

	isValid, errors := models.ValidatePassword(req.Password, h.passwordValidation)

	response := struct {
		Valid  bool     `json:"valid"`
		Errors []string `json:"errors,omitempty"`
		Score  int      `json:"score"`
	}{
		Valid:  isValid,
		Errors: errors,
		Score:  calculatePasswordScore(req.Password),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ValidateRegistration валидация регистрационных данных
func (h *AuthHandler) validateRegistration(req models.RegisterRequest) []string {
	var errors []string

	// Валидация имени пользователя
	if valid, usernameErrors := models.ValidateUsername(req.Username); !valid {
		errors = append(errors, usernameErrors...)
	}

	// Валидация email
	if !models.ValidateEmail(req.Email) {
		errors = append(errors, "Некорректный email адрес")
	}

	// Валидация пароля
	if valid, passwordErrors := models.ValidatePassword(req.Password, h.passwordValidation); !valid {
		errors = append(errors, passwordErrors...)
	}

	// Проверка совпадения паролей
	if req.Password != req.ConfirmPassword {
		errors = append(errors, "Пароли не совпадают")
	}

	return errors
}

// ValidateLogin валидация данных входа
func (h *AuthHandler) validateLogin(req models.LoginRequest) []string {
	var errors []string

	if strings.TrimSpace(req.Username) == "" {
		errors = append(errors, "Имя пользователя обязательно")
	}

	if strings.TrimSpace(req.Password) == "" {
		errors = append(errors, "Пароль обязателен")
	}

	return errors
}

func (h *AuthHandler) generateToken(username string, isAdmin bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"is_admin": isAdmin, // <-- claim
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
		"type":     "access",
	})

	return token.SignedString([]byte(h.jwtSecret))
}

// SendError отправляет ошибку в формате JSON
func sendError(w http.ResponseWriter, message string, statusCode int, errors []string) {
	response := models.ErrorResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// CalculatePasswordScore рассчитывает сложность пароля
func calculatePasswordScore(password string) int {
	score := 0

	// Базовая оценка за длину
	if len(password) >= 8 {
		score += 10
	}
	if len(password) >= 12 {
		score += 10
	}
	if len(password) >= 16 {
		score += 10
	}

	// Проверка наличия различных типов символов
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		default:
			hasSpecial = true
		}
	}

	if hasUpper {
		score += 10
	}
	if hasLower {
		score += 10
	}
	if hasNumber {
		score += 10
	}
	if hasSpecial {
		score += 20
	}

	// Штраф за повторяющиеся символы
	if hasRepeatingCharacters(password) {
		score -= 15
	}

	// Штраф за последовательности
	if hasSequentialCharacters(password) {
		score -= 20
	}

	// Ограничение оценки от 0 до 100
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

func hasRepeatingCharacters(password string) bool {
	if len(password) < 3 {
		return false
	}
	for i := 0; i < len(password)-2; i++ {
		if password[i] == password[i+1] && password[i] == password[i+2] {
			return true
		}
	}
	return false
}

func hasSequentialCharacters(password string) bool {
	if len(password) < 3 {
		return false
	}
	for i := 0; i < len(password)-2; i++ {
		c1 := password[i]
		c2 := password[i+1]
		c3 := password[i+2]

		// Проверка числовых последовательностей
		if '0' <= c1 && c1 <= '9' && '0' <= c2 && c2 <= '9' && '0' <= c3 && c3 <= '9' {
			if c2 == c1+1 && c3 == c2+1 {
				return true
			}
			if c2 == c1-1 && c3 == c2-1 {
				return true
			}
		}

		// Проверка буквенных последовательностей (только для строчных)
		if 'a' <= c1 && c1 <= 'z' && 'a' <= c2 && c2 <= 'z' && 'a' <= c3 && c3 <= 'z' {
			if c2 == c1+1 && c3 == c2+1 {
				return true
			}
			if c2 == c1-1 && c3 == c2-1 {
				return true
			}
		}
	}
	return false
}

// GetAllUsers возвращает всех пользователей (для отладки)
func (h *AuthHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userRepo.GetAllUsers()
	if err != nil {
		sendError(w, "Ошибка при получении пользователей", http.StatusInternalServerError, nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// HealthCheck проверка работоспособности сервера
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":  "ok",
		"service": "Renault Backend API",
		"time":    time.Now().Format(time.RFC3339),
		"version": "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// PasswordRules возвращает правила для пароля
func (h *AuthHandler) PasswordRules(w http.ResponseWriter, r *http.Request) {
	rules := map[string]interface{}{
		"min_length":      h.passwordValidation.MinLength,
		"require_upper":   h.passwordValidation.RequireUpper,
		"require_lower":   h.passwordValidation.RequireLower,
		"require_number":  h.passwordValidation.RequireNumber,
		"require_special": h.passwordValidation.RequireSpecial,
		"rules": []string{
			"Минимум 8 символов",
			"Хотя бы одна заглавная буква",
			"Хотя бы одна строчная буква",
			"Хотя бы одна цифра",
			"Хотя бы один специальный символ (!@#$%^&* и т.д.)",
			"Не использовать простые пароли (password, 123456 и т.д.)",
			"Не использовать повторяющиеся символы (aaa, 111)",
			"Не использовать последовательности (abc, 123)",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rules)
}

package models

import (
	"regexp"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	User    *User  `json:"user,omitempty"`
}

type ErrorResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}

type PasswordValidation struct {
	MinLength       int
	RequireUpper    bool
	RequireLower    bool
	RequireNumber   bool
	RequireSpecial  bool
	CommonPasswords []string
}

var DefaultPasswordValidation = PasswordValidation{
	MinLength:      8,
	RequireUpper:   true,
	RequireLower:   true,
	RequireNumber:  true,
	RequireSpecial: true,
	CommonPasswords: []string{
		"password", "123456", "12345678", "123456789", "1234567890",
		"qwerty", "abc123", "password1", "admin", "letmein",
		"welcome", "monkey", "dragon", "baseball", "football",
		"hello", "master", "sunshine", "password123", "superman",
	},
}

// HashPassword хеширует пароль
func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// CheckPassword проверяет пароль
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// ValidatePassword проверяет сложность пароля
func ValidatePassword(password string, validation PasswordValidation) (bool, []string) {
	var errors []string

	// Проверка длины
	if len(password) < validation.MinLength {
		errors = append(errors,
			"Пароль должен содержать минимум 8 символов")
	}

	// Проверка наличия различных типов символов
	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if validation.RequireUpper && !hasUpper {
		errors = append(errors,
			"Пароль должен содержать хотя бы одну заглавную букву")
	}

	if validation.RequireLower && !hasLower {
		errors = append(errors,
			"Пароль должен содержать хотя бы одну строчную букву")
	}

	if validation.RequireNumber && !hasNumber {
		errors = append(errors,
			"Пароль должен содержать хотя бы одну цифру")
	}

	if validation.RequireSpecial && !hasSpecial {
		errors = append(errors,
			"Пароль должен содержать хотя бы один специальный символ (!@#$%^&*)")
	}

	// Проверка на распространенные пароли
	lowerPassword := strings.ToLower(password)
	for _, common := range validation.CommonPasswords {
		if strings.Contains(lowerPassword, common) {
			errors = append(errors,
				"Пароль слишком простой и распространенный")
			break
		}
	}

	// Проверка на последовательности
	if hasSequentialCharacters(password) {
		errors = append(errors,
			"Пароль содержит слишком простые последовательности символов")
	}

	// Проверка на повторяющиеся символы
	if hasRepeatingCharacters(password) {
		errors = append(errors,
			"Пароль содержит слишком много повторяющихся символов")
	}

	return len(errors) == 0, errors
}

// Проверка на последовательные символы
func hasSequentialCharacters(password string) bool {
	if len(password) < 3 {
		return false
	}

	for i := 0; i < len(password)-2; i++ {
		// Проверка числовых последовательностей
		if isDigit(password[i]) && isDigit(password[i+1]) && isDigit(password[i+2]) {
			c1 := int(password[i] - '0')
			c2 := int(password[i+1] - '0')
			c3 := int(password[i+2] - '0')

			// Проверка возрастающей последовательности
			if c2 == c1+1 && c3 == c2+1 {
				return true
			}
			// Проверка убывающей последовательности
			if c2 == c1-1 && c3 == c2-1 {
				return true
			}
		}

		// Проверка буквенных последовательностей
		if isLetter(password[i]) && isLetter(password[i+1]) && isLetter(password[i+2]) {
			c1 := strings.ToLower(string(password[i]))[0]
			c2 := strings.ToLower(string(password[i+1]))[0]
			c3 := strings.ToLower(string(password[i+2]))[0]

			// Проверка алфавитной последовательности
			if c2 == c1+1 && c3 == c2+1 {
				return true
			}
			// Проверка обратной последовательности
			if c2 == c1-1 && c3 == c2-1 {
				return true
			}
		}
	}

	return false
}

// Проверка на повторяющиеся символы
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

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// ValidateEmail проверяет email
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidateUsername проверяет имя пользователя
func ValidateUsername(username string) (bool, []string) {
	var errors []string

	if len(username) < 3 {
		errors = append(errors,
			"Имя пользователя должно содержать минимум 3 символа")
	}

	if len(username) > 20 {
		errors = append(errors,
			"Имя пользователя не должно превышать 20 символов")
	}

	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	if !usernameRegex.MatchString(username) {
		errors = append(errors,
			"Имя пользователя может содержать только буквы, цифры, точки, дефисы и подчеркивания")
	}

	return len(errors) == 0, errors
}

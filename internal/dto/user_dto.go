package dto

// UserRegisterRequest DTO для регистрации
// @Description Данные, необходимые для регистрации пользователя
type UserRegisterRequest struct {
	// Email пользователя (обязательное поле)
	// Required: true
	// Example: user@example.com
	Email string `json:"email" binding:"required,email"`

	// Username пользователя (обязательное поле)
	// Required: true
	// Example: }{0ТТ@БЬ)Ч
	Username string `json:"username" binding:"required"`

	// Пароль пользователя (минимум 6 символов, обязательное поле)
	// Required: true
	// Example: mysecurepassword
	Password string `json:"password" binding:"required,min=6"`
}

// UserLoginRequest DTO для входа
// @Description Данные для авторизации пользователя
type UserLoginRequest struct {
	// Email пользователя (обязательное поле)
	// Required: true
	// Example: user@example.com
	Email string `json:"email" binding:"required,email"`

	// Пароль (обязательное поле)
	// Required: true
	// Example: mysecurepassword
	Password string `json:"password" binding:"required"`
}

// TokenRefreshRequest DTO для обновления токена
// @Description Запрос на обновление access-токена с использованием refresh-токена
type TokenRefreshRequest struct {
	// Refresh-токен (обязательное поле)
	// Required: true
	// Example: eyJhbGciOiJIUzI1NiIsInR...
	Token string `json:"token" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserResponse struct {
	ID    string `json:"access_token"`
	Email string `json:"refresh_token"`
	Role  string `json:"role"`
}

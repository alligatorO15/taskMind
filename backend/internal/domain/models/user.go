package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User - модель пользователя системы TaskMind
// содержит учетные данные и метаинфу
type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username     string             `json:"username" bson:"username"`
	Email        string             `json:"email" bson:"email"`
	PasswordHash string             `json:"-" bson:"password_hash"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

// UserRegisterRequest - запрос на регистрацию нового пользователя
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=128"`
}

// UserLoginRequest - запрос на логин пользователя
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// TokenPair - пара токенов доступа: access и refresh
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenRequest - запрос на обновление access-токена
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// UserResponse - ответ с данными пользователя (без пароля)
type UserResponse struct {
	ID        primitive.ObjectID `json:"id"`
	Username  string             `json:"username"`
	Email     string             `json:"email"`
	CreatedAt time.Time          `json:"created_at"`
}

// ToResponse - преобразует модель пользователя в безопасный ответ без пароля
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Email,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

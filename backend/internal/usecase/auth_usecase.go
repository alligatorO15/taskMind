package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alligatorO15/taskMind/backend/pkg/apperror"
	"github.com/alligatorO15/taskmind-backend/internal/config"
	"github.com/alligatorO15/taskmind-backend/internal/domain/models"
	"github.com/alligatorO15/taskmind-backend/internal/domain/repository"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// AuthUseCase - бизнес-логика аутентификации и авторизации пользователей
// Отвечает за регистрацию, вход, генерацию и обновление JWT-токенов
type AuthUseCase struct {
	userRepo repository.UserRepository
	jwtCfg   config.JWTConfig
}

// NewAuthUseCase - создает экземпляр AuthUseCase с необходимыми зависимостями
func NewAuthUseCase(userRepo repository.UserRepository, jwtCfg config.JWTConfig) *AuthUseCase {
	return &AuthUseCase{
		userRepo: userRepo,
		jwtCfg:   jwtCfg,
	}
}

// Register - регистрирует нового пользователя
// Хеширует пароль, сохраняет пользователя и возвраащет пару токенов.
func (uc *AuthUseCase) Register(ctx context.Context, req models.UserRegisterRequest) (*models.TokenPair, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, apperror.ErrAlreadyExists) {
			return nil, fmt.Errorf("пользователь с таким email или именем уже существует: %w", err)
		}
		return nil, err
	}

	return uc.generateTokenPair(user.ID)
}

// Login - авторизует пользователя по email и паролю
// Проверяет учетные данные и возвращает пару токенов при успехе
func (uc *AuthUseCase) Login(ctx context.Context, req models.UserLoginRequest) (*models.TokenPair, error) {
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperror.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperror.ErrInvalidCredentials
	}

	return uc.generateTokenPair(user.ID)
}

// RefreshToken - обновляет access-токенм используя валидный refersh-токен
// Проверяет подпись refersh-токена и генерирует новую пару
func (uc *AuthUseCase) RefreshToken(ctx context.Context, refreshToken string) (*models.TokenPair, error) {
	claims, err := uc.parseToken(refreshToken, uc.jwtCfg.RefreshSecret)
	if err != nil {
		return nil, apperror.ErrUnauthorized
	}

	userID, err := primitive.ObjectIDFromHex(claims["user_id"].(string))
	if err != nil {
		return nil, apperror.ErrUnauthorized
	}

	// Проверка что юзер еще существует
	_, err = uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, apperror.ErrUnauthorized
	}

	return uc.generateTokenPair(userID)
}

// GetUserByID - возвращает пользователя по идентификатору
func (uc *AuthUseCase) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	return uc.userRepo.FindByID(ctx, id)
}

// generateTokenPair - генерирует пару JWT-токенов (access + refresh) для пользователя
func (uc *AuthUseCase) generateTokenPair(userID primitive.ObjectID) (*models.TokenPair, error) {
	accessToken, err := uc.generateToken(userID, uc.jwtCfg.AccessSecret, uc.jwtCfg.AccessTTL)
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.generateToken(userID, uc.jwtCfg.RefreshSecret, uc.jwtCfg.RefreshTTL)
	if err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// generateToken - создает JWT-токен с указанным временем жизни и секретным ключом
func (uc *AuthUseCase) generateToken(userId primitive.ObjectID, secret string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId.Hex(),
		"exp":     time.Now().Add(ttl).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// parseToken - разбирает и валидирует JWT-токен, возвращает claims при успехе
func (uc *AuthUseCase) parseToken(tokenString string, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperror.ErrUnauthorized
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, apperror.ErrUnauthorized
	}

	return claims, nil
}

// ParseAccessToken - публичный метод для парсинга access-токена
// Используется в middleware аутентификации
// немного нарушает чистую архитектуру, в больших проектаз надо:
// - вынести jwt-логику за интерфейс TokenService в domain/ports/token_service.go
// - реализацию в инфраструктурный слой
// - AuthUseCase struct - должно зависеть только от интерфейса а не конкретного конфига (ports.TokenService)
func (uc *AuthUseCase) ParseAccessToken(tokenString string) (primitive.ObjectID, error) {
	claims, err := uc.parseToken(tokenString, uc.jwtCfg.AccessSecret)
	if err != nil {
		return primitive.NilObjectID, apperror.ErrUnauthorized
	}
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return primitive.NilObjectID, apperror.ErrUnauthorized
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return primitive.NilObjectID, apperror.ErrUnauthorized
	}

	return userID, nil
}

package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/alligatorO15/taskMind/backend/pkg/apperror"
	"github.com/alligatorO15/taskmind-backend/internal/config"
	"github.com/alligatorO15/taskmind-backend/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// mockUserRepo — мок-реализация UserRepository для тестирования
type mockUserRepo struct {
	users map[string]*models.User // email -> user
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users: make(map[string]*models.User),
	}
}

func (m *mockUserRepo) Create(ctx context.Context, user *models.User) error {
	if _, exists := m.users[user.Email]; exists {
		return apperror.ErrAlreadyExists
	}
	user.ID = primitive.NewObjectID()
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, apperror.ErrNotFound
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	u, ok := m.users[email]
	if !ok {
		return nil, apperror.ErrNotFound
	}
	return u, nil
}

func (m *mockUserRepo) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	for _, u := range m.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, apperror.ErrNotFound
}

// testJWTConfig — тестовая конфигурация JWT
func testJWTConfig() config.JWTConfig {
	return config.JWTConfig{
		AccessSecret:  "test-access-secret",
		RefreshSecret: "test-refresh-secret",
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    7 * 24 * time.Hour,
	}
}

// TestRegister_Success — проверяет успешную регистрацию нового пользователя
func TestRegister_Success(t *testing.T) {
	repo := newMockUserRepo()
	uc := NewAuthUseCase(repo, testJWTConfig())

	req := models.UserRegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	tokenPair, err := uc.Register(context.Background(), req)
	if err != nil {
		t.Fatalf("Ошибка регистрации: %v", err)
	}

	if tokenPair.AccessToken == "" {
		t.Error("Access-токен не должен быть пустым")
	}
	if tokenPair.RefreshToken == "" {
		t.Error("Refresh-токен не должен быть пустым")
	}

	// Проверяем, что пользователь сохранён
	user, exists := repo.users["test@example.com"]
	if !exists {
		t.Fatal("Пользователь не был сохранён в репозитории")
	}

	// Проверяем, что пароль захеширован
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("password123"))
	if err != nil {
		t.Error("Пароль неверно захеширован")
	}
}

// TestRegister_DuplicateEmail — проверяет, что повторная регистрация с тем же email отклоняется
func TestRegister_DuplicateEmail(t *testing.T) {
	repo := newMockUserRepo()
	uc := NewAuthUseCase(repo, testJWTConfig())

	req := models.UserRegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	_, err := uc.Register(context.Background(), req)
	if err != nil {
		t.Fatalf("Первая регистрация не должна была вернуть ошибку: %v", err)
	}

	// Повторная регистрация с тем же email
	req.Username = "testuser2"
	_, err = uc.Register(context.Background(), req)
	if err != apperror.ErrAlreadyExists {
		t.Errorf("Ожидалась ошибка ErrAlreadyExists, получили: %v", err)
	}
}

// TestLogin_Success — проверяет успешную авторизацию
func TestLogin_Success(t *testing.T) {
	repo := newMockUserRepo()
	uc := NewAuthUseCase(repo, testJWTConfig())

	// Сначала регистрируем пользователя
	regReq := models.UserRegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	_, err := uc.Register(context.Background(), regReq)
	if err != nil {
		t.Fatalf("Ошибка регистрации: %v", err)
	}

	// Авторизуемся
	loginReq := models.UserLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	tokenPair, err := uc.Login(context.Background(), loginReq)
	if err != nil {
		t.Fatalf("Ошибка авторизации: %v", err)
	}

	if tokenPair.AccessToken == "" {
		t.Error("Access-токен не должен быть пустым")
	}
}

// TestLogin_WrongPassword — проверяет отклонение при неверном пароле
func TestLogin_WrongPassword(t *testing.T) {
	repo := newMockUserRepo()
	uc := NewAuthUseCase(repo, testJWTConfig())

	regReq := models.UserRegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	_, _ = uc.Register(context.Background(), regReq)

	loginReq := models.UserLoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	_, err := uc.Login(context.Background(), loginReq)
	if err != apperror.ErrInvalidCredentials {
		t.Errorf("Ожидалась ошибка ErrInvalidCredentials, получили: %v", err)
	}
}

// TestLogin_UserNotFound — проверяет отклонение при несуществующем email
func TestLogin_UserNotFound(t *testing.T) {
	repo := newMockUserRepo()
	uc := NewAuthUseCase(repo, testJWTConfig())

	loginReq := models.UserLoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	_, err := uc.Login(context.Background(), loginReq)
	if err != apperror.ErrInvalidCredentials {
		t.Errorf("Ожидалась ошибка ErrInvalidCredentials, получили: %v", err)
	}
}

// TestParseAccessToken — проверяет генерацию и парсинг access-токена
func TestParseAccessToken(t *testing.T) {
	repo := newMockUserRepo()
	uc := NewAuthUseCase(repo, testJWTConfig())

	regReq := models.UserRegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	tokenPair, err := uc.Register(context.Background(), regReq)
	if err != nil {
		t.Fatalf("Ошибка регистрации: %v", err)
	}

	// Парсим access-токен
	userID, err := uc.ParseAccessToken(tokenPair.AccessToken)
	if err != nil {
		t.Fatalf("Ошибка парсинга токена: %v", err)
	}

	if userID.IsZero() {
		t.Error("UserID не должен быть нулевым")
	}

	// Проверяем, что ID совпадает
	user, _ := repo.FindByEmail(context.Background(), "test@example.com")
	if userID != user.ID {
		t.Errorf("UserID не совпадает: ожидался %s, получили %s", user.ID.Hex(), userID.Hex())
	}
}

// TestRefreshToken_Success — проверяет обновление токенов
func TestRefreshToken_Success(t *testing.T) {
	repo := newMockUserRepo()
	uc := NewAuthUseCase(repo, testJWTConfig())

	regReq := models.UserRegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	tokenPair, err := uc.Register(context.Background(), regReq)
	if err != nil {
		t.Fatalf("Ошибка регистрации: %v", err)
	}

	// Обновляем токены
	newPair, err := uc.RefreshToken(context.Background(), tokenPair.RefreshToken)
	if err != nil {
		t.Fatalf("Ошибка обновления токена: %v", err)
	}

	if newPair.AccessToken == "" {
		t.Error("Новый access-токен не должен быть пустым")
	}
	if newPair.RefreshToken == "" {
		t.Error("Новый refresh-токен не должен быть пустым")
	}
}

// TestRefreshToken_InvalidToken — проверяет отклонение невалидного refresh-токена
func TestRefreshToken_InvalidToken(t *testing.T) {
	repo := newMockUserRepo()
	uc := NewAuthUseCase(repo, testJWTConfig())

	_, err := uc.RefreshToken(context.Background(), "invalid-token")
	if err != apperror.ErrUnauthorized {
		t.Errorf("Ожидалась ошибка ErrUnauthorized, получили: %v", err)
	}
}

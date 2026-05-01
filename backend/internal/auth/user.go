package auth

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"vue-api/backend/internal/id"
)

const (
	GlobalRoleManager = "manager"
	GlobalRoleUser    = "user"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrEmailAlreadyInUse    = errors.New("email already in use")
	ErrUsernameAlreadyInUse = errors.New("username already in use")
)

type User struct {
	ID           string
	Email        string
	Username     string
	PasswordHash string
	GlobalRole   string
	TokenVersion int
	Active       bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CreateUserParams struct {
	Email        string
	Username     string
	PasswordHash string
	GlobalRole   string
}

type UserRepository interface {
	CountUsers(ctx context.Context) (int, error)
	CreateUser(ctx context.Context, params CreateUserParams) (User, error)
	FindUserByEmail(ctx context.Context, email string) (User, error)
	FindUserByUsername(ctx context.Context, username string) (User, error)
	FindUserByID(ctx context.Context, id string) (User, error)
	UpdateUsername(ctx context.Context, userID string, username string) (User, error)
}

type MemoryUserRepository struct {
	mu         sync.RWMutex
	users      map[string]User
	byEmail    map[string]string
	byUsername map[string]string
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users:      make(map[string]User),
		byEmail:    make(map[string]string),
		byUsername: make(map[string]string),
	}
}

func (repo *MemoryUserRepository) CountUsers(ctx context.Context) (int, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	return len(repo.users), nil
}

func (repo *MemoryUserRepository) CreateUser(ctx context.Context, params CreateUserParams) (User, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	email := NormalizeEmail(params.Email)
	username := NormalizeUsername(params.Username)
	if _, exists := repo.byEmail[email]; exists {
		return User{}, ErrEmailAlreadyInUse
	}
	if _, exists := repo.byUsername[username]; exists {
		return User{}, ErrUsernameAlreadyInUse
	}

	now := time.Now().UTC()
	user := User{
		ID:           newUserID(),
		Email:        email,
		Username:     username,
		PasswordHash: params.PasswordHash,
		GlobalRole:   params.GlobalRole,
		TokenVersion: 1,
		Active:       true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	repo.users[user.ID] = user
	repo.byEmail[email] = user.ID
	repo.byUsername[username] = user.ID

	return user, nil
}

func (repo *MemoryUserRepository) FindUserByEmail(ctx context.Context, email string) (User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	id, exists := repo.byEmail[NormalizeEmail(email)]
	if !exists {
		return User{}, ErrUserNotFound
	}

	return repo.users[id], nil
}

func (repo *MemoryUserRepository) FindUserByUsername(ctx context.Context, username string) (User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	id, exists := repo.byUsername[NormalizeUsername(username)]
	if !exists {
		return User{}, ErrUserNotFound
	}

	return repo.users[id], nil
}

func (repo *MemoryUserRepository) FindUserByID(ctx context.Context, id string) (User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	user, exists := repo.users[id]
	if !exists {
		return User{}, ErrUserNotFound
	}

	return user, nil
}

func (repo *MemoryUserRepository) UpdateUsername(ctx context.Context, userID string, username string) (User, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	user, exists := repo.users[userID]
	if !exists {
		return User{}, ErrUserNotFound
	}

	nextUsername := NormalizeUsername(username)
	if existingID, exists := repo.byUsername[nextUsername]; exists && existingID != userID {
		return User{}, ErrUsernameAlreadyInUse
	}

	delete(repo.byUsername, user.Username)
	user.Username = nextUsername
	user.UpdatedAt = time.Now().UTC()
	repo.users[user.ID] = user
	repo.byUsername[nextUsername] = user.ID

	return user, nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func newUserID() string {
	return id.NewUUIDV7()
}

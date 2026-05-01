package auth

import (
	"context"
	"errors"
	"strings"
)

type BootstrapConfig struct {
	Enabled  bool
	Email    string
	Username string
	Password string
}

const (
	BootstrapManagerDisabled             = "disabled"
	BootstrapManagerSkippedExistingUsers = "skipped_existing_users"
	BootstrapManagerSeeded               = "seeded"
)

type BootstrapResult struct {
	Status   string
	UserID   string
	Email    string
	Username string
}

func BootstrapManager(ctx context.Context, repo UserRepository, hasher PasswordHasher, cfg BootstrapConfig) (BootstrapResult, error) {
	if !cfg.Enabled {
		return BootstrapResult{Status: BootstrapManagerDisabled}, nil
	}

	count, err := repo.CountUsers(ctx)
	if err != nil {
		return BootstrapResult{}, err
	}
	if count > 0 {
		return BootstrapResult{Status: BootstrapManagerSkippedExistingUsers}, nil
	}

	email := NormalizeEmail(cfg.Email)
	if email == "" || !strings.Contains(email, "@") {
		return BootstrapResult{}, errors.New("BOOTSTRAP_MANAGER_EMAIL must be a valid email")
	}
	username := NormalizeUsername(cfg.Username)
	if username == "" {
		return BootstrapResult{}, errors.New("BOOTSTRAP_MANAGER_USERNAME is required when BOOTSTRAP_MANAGER_ENABLED is true")
	}
	if len(cfg.Password) < 12 {
		return BootstrapResult{}, errors.New("BOOTSTRAP_MANAGER_PASSWORD must be at least 12 characters")
	}

	hash, err := hasher.Hash(cfg.Password)
	if err != nil {
		return BootstrapResult{}, err
	}

	user, err := repo.CreateUser(ctx, CreateUserParams{
		Email:        email,
		Username:     username,
		PasswordHash: hash,
		GlobalRole:   GlobalRoleManager,
	})
	if err != nil {
		return BootstrapResult{}, err
	}

	return BootstrapResult{
		Status:   BootstrapManagerSeeded,
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
	}, nil
}

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

func BootstrapManager(ctx context.Context, repo UserRepository, hasher PasswordHasher, cfg BootstrapConfig) error {
	if !cfg.Enabled {
		return nil
	}

	count, err := repo.CountUsers(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	email := NormalizeEmail(cfg.Email)
	if email == "" || !strings.Contains(email, "@") {
		return errors.New("BOOTSTRAP_MANAGER_EMAIL must be a valid email")
	}
	username := NormalizeUsername(cfg.Username)
	if username == "" {
		return errors.New("BOOTSTRAP_MANAGER_USERNAME is required when BOOTSTRAP_MANAGER_ENABLED is true")
	}
	if len(cfg.Password) < 12 {
		return errors.New("BOOTSTRAP_MANAGER_PASSWORD must be at least 12 characters")
	}

	hash, err := hasher.Hash(cfg.Password)
	if err != nil {
		return err
	}

	_, err = repo.CreateUser(ctx, CreateUserParams{
		Email:        email,
		Username:     username,
		PasswordHash: hash,
		GlobalRole:   GlobalRoleManager,
	})
	return err
}

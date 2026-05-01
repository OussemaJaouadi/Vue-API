package gormstorage

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"vue-api/backend/internal/auth"
	"vue-api/backend/internal/id"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) CountUsers(ctx context.Context) (int, error) {
	var count int64
	if err := repo.db.WithContext(ctx).Model(&UserModel{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (repo *UserRepository) CreateUser(ctx context.Context, params auth.CreateUserParams) (auth.User, error) {
	model := UserModel{
		ID:           id.NewUUIDV7(),
		Email:        auth.NormalizeEmail(params.Email),
		Username:     auth.NormalizeUsername(params.Username),
		PasswordHash: params.PasswordHash,
		GlobalRole:   params.GlobalRole,
		TokenVersion: 1,
		Active:       true,
	}

	if err := repo.db.WithContext(ctx).Create(&model).Error; err != nil {
		return auth.User{}, mapCreateUserError(err)
	}

	return toAuthUser(model), nil
}

func (repo *UserRepository) FindUserByEmail(ctx context.Context, email string) (auth.User, error) {
	var model UserModel
	err := repo.db.WithContext(ctx).
		Where("email = ?", auth.NormalizeEmail(email)).
		First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return auth.User{}, auth.ErrUserNotFound
	}
	if err != nil {
		return auth.User{}, err
	}

	return toAuthUser(model), nil
}

func (repo *UserRepository) FindUserByUsername(ctx context.Context, username string) (auth.User, error) {
	var model UserModel
	err := repo.db.WithContext(ctx).
		Where("username = ?", auth.NormalizeUsername(username)).
		First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return auth.User{}, auth.ErrUserNotFound
	}
	if err != nil {
		return auth.User{}, err
	}

	return toAuthUser(model), nil
}

func (repo *UserRepository) FindUserByID(ctx context.Context, userID string) (auth.User, error) {
	var model UserModel
	err := repo.db.WithContext(ctx).
		Where("id = ?", userID).
		First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return auth.User{}, auth.ErrUserNotFound
	}
	if err != nil {
		return auth.User{}, err
	}

	return toAuthUser(model), nil
}

func (repo *UserRepository) UpdateUsername(ctx context.Context, userID string, username string) (auth.User, error) {
	var model UserModel
	err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", userID).First(&model).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return auth.ErrUserNotFound
			}
			return err
		}

		model.Username = auth.NormalizeUsername(username)
		if err := tx.Save(&model).Error; err != nil {
			return mapCreateUserError(err)
		}

		return nil
	})
	if err != nil {
		return auth.User{}, err
	}

	return toAuthUser(model), nil
}

func mapCreateUserError(err error) error {
	message := strings.ToLower(err.Error())
	if strings.Contains(message, "idx_users_email") || strings.Contains(message, "users.email") {
		return auth.ErrEmailAlreadyInUse
	}
	if strings.Contains(message, "idx_users_username") || strings.Contains(message, "users.username") {
		return auth.ErrUsernameAlreadyInUse
	}

	return err
}

func toAuthUser(model UserModel) auth.User {
	return auth.User{
		ID:           model.ID,
		Email:        model.Email,
		Username:     model.Username,
		PasswordHash: model.PasswordHash,
		GlobalRole:   model.GlobalRole,
		TokenVersion: model.TokenVersion,
		Active:       model.Active,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}
}

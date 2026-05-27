package gormstorage

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"vue-api/backend/internal/environment"
	"vue-api/backend/internal/id"
)

type EnvironmentRepository struct {
	db *gorm.DB
}

func NewEnvironmentRepository(db *gorm.DB) *EnvironmentRepository {
	return &EnvironmentRepository{db: db}
}

func (repo *EnvironmentRepository) ListEnvironments(ctx context.Context, workspaceID string) ([]environment.Environment, error) {
	var models []EnvironmentModel
	if err := repo.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Order("sort_order ASC, name ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	envs := make([]environment.Environment, 0, len(models))
	for _, m := range models {
		envs = append(envs, toDomainEnvironment(m))
	}

	return envs, nil
}

func (repo *EnvironmentRepository) ListVariables(ctx context.Context, environmentID string) ([]environment.Variable, error) {
	var models []VariableModel
	if err := repo.db.WithContext(ctx).Where("environment_id = ?", environmentID).Order("sort_order ASC, key ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	vars := make([]environment.Variable, 0, len(models))
	for _, m := range models {
		vars = append(vars, toDomainVariable(m))
	}

	return vars, nil
}

func (repo *EnvironmentRepository) CreateEnvironment(ctx context.Context, params environment.CreateEnvironmentParams) (environment.Environment, error) {
	var count int64
	if err := repo.db.WithContext(ctx).Model(&EnvironmentModel{}).Where("workspace_id = ? AND name = ?", params.WorkspaceID, params.Name).Count(&count).Error; err != nil {
		return environment.Environment{}, err
	}
	if count > 0 {
		return environment.Environment{}, environment.ErrEnvironmentNameTaken
	}

	visibility := params.Visibility
	if visibility == "" {
		visibility = "project"
	}

	model := EnvironmentModel{
		ID:              id.NewUUIDV7(),
		WorkspaceID:     params.WorkspaceID,
		Name:            params.Name,
		Visibility:      visibility,
		SortOrder:       0,
		CreatedByUserID: params.CreatedByUserID,
	}

	if err := repo.db.WithContext(ctx).Create(&model).Error; err != nil {
		return environment.Environment{}, err
	}

	return toDomainEnvironment(model), nil
}

func (repo *EnvironmentRepository) CreateVariable(ctx context.Context, params environment.CreateVariableParams) (environment.Variable, error) {
	var count int64
	if err := repo.db.WithContext(ctx).Model(&VariableModel{}).Where("environment_id = ? AND key = ?", params.EnvironmentID, params.Key).Count(&count).Error; err != nil {
		return environment.Variable{}, err
	}
	if count > 0 {
		return environment.Variable{}, environment.ErrVariableKeyTaken
	}

	model := VariableModel{
		ID:            id.NewUUIDV7(),
		EnvironmentID: params.EnvironmentID,
		Key:           params.Key,
		Value:         params.Value,
		Secret:        params.Secret,
		SortOrder:     0,
	}

	if err := repo.db.WithContext(ctx).Create(&model).Error; err != nil {
		return environment.Variable{}, err
	}

	return toDomainVariable(model), nil
}

func (repo *EnvironmentRepository) UpdateEnvironment(ctx context.Context, envID string, params environment.UpdateEnvironmentParams) (environment.Environment, error) {
	var model EnvironmentModel
	if err := repo.db.WithContext(ctx).Where("id = ?", envID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return environment.Environment{}, environment.ErrEnvironmentNotFound
		}
		return environment.Environment{}, err
	}

	if params.Name != nil {
		model.Name = *params.Name
	}
	if params.Visibility != nil {
		model.Visibility = *params.Visibility
	}

	if err := repo.db.WithContext(ctx).Save(&model).Error; err != nil {
		return environment.Environment{}, err
	}

	return toDomainEnvironment(model), nil
}

func (repo *EnvironmentRepository) UpdateVariable(ctx context.Context, varID string, params environment.UpdateVariableParams) (environment.Variable, error) {
	var model VariableModel
	if err := repo.db.WithContext(ctx).Where("id = ?", varID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return environment.Variable{}, environment.ErrVariableNotFound
		}
		return environment.Variable{}, err
	}

	if params.Key != nil {
		model.Key = *params.Key
	}
	if params.Value != nil {
		model.Value = *params.Value
	}
	if params.Secret != nil {
		model.Secret = *params.Secret
	}

	if err := repo.db.WithContext(ctx).Save(&model).Error; err != nil {
		return environment.Variable{}, err
	}

	return toDomainVariable(model), nil
}

func (repo *EnvironmentRepository) DeleteEnvironment(ctx context.Context, envID string) error {
	result := repo.db.WithContext(ctx).Where("id = ?", envID).Delete(&EnvironmentModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return environment.ErrEnvironmentNotFound
	}

	if err := repo.db.WithContext(ctx).Where("environment_id = ?", envID).Delete(&VariableModel{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *EnvironmentRepository) DeleteVariable(ctx context.Context, varID string) error {
	result := repo.db.WithContext(ctx).Where("id = ?", varID).Delete(&VariableModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return environment.ErrVariableNotFound
	}

	return nil
}

func toDomainEnvironment(model EnvironmentModel) environment.Environment {
	return environment.Environment{
		ID:              model.ID,
		WorkspaceID:     model.WorkspaceID,
		Name:            model.Name,
		Visibility:      model.Visibility,
		SortOrder:       model.SortOrder,
		CreatedByUserID: model.CreatedByUserID,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
	}
}

func toDomainVariable(model VariableModel) environment.Variable {
	return environment.Variable{
		ID:            model.ID,
		EnvironmentID: model.EnvironmentID,
		Key:           model.Key,
		Value:         model.Value,
		Secret:        model.Secret,
		SortOrder:     model.SortOrder,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}
}

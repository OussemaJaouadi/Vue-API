package gormstorage

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"vue-api/backend/internal/id"
	"vue-api/backend/internal/workspace"
)

type WorkspaceRepository struct {
	db *gorm.DB
}

func NewWorkspaceRepository(db *gorm.DB) *WorkspaceRepository {
	return &WorkspaceRepository{db: db}
}

func (r *WorkspaceRepository) ListByUser(ctx context.Context, userID string) ([]workspace.Workspace, error) {
	var models []WorkspaceModel
	if err := r.db.WithContext(ctx).
		Joins("JOIN workspace_memberships ON workspace_memberships.workspace_id = workspaces.id").
		Where("workspace_memberships.user_id = ?", userID).
		Order("workspaces.created_at ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]workspace.Workspace, len(models))
	for i, m := range models {
		result[i] = toWorkspace(m)
	}
	return result, nil
}

func (r *WorkspaceRepository) Create(ctx context.Context, params workspace.CreateWorkspaceParams) (workspace.Workspace, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&WorkspaceModel{}).
		Where("name = ? AND created_by_user_id = ?", params.Name, params.CreatedByUserID).
		Count(&count).Error; err != nil {
		return workspace.Workspace{}, err
	}
	if count > 0 {
		return workspace.Workspace{}, workspace.ErrWorkspaceNameTaken
	}

	now := time.Now().UTC()
	model := WorkspaceModel{
		ID:              id.NewUUIDV7(),
		Name:            params.Name,
		CreatedByUserID: params.CreatedByUserID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return workspace.Workspace{}, err
	}
	return toWorkspace(model), nil
}

func (r *WorkspaceRepository) FindByID(ctx context.Context, id string) (workspace.Workspace, error) {
	var model WorkspaceModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return workspace.Workspace{}, workspace.ErrWorkspaceNotFound
		}
		return workspace.Workspace{}, err
	}
	return toWorkspace(model), nil
}

func (r *WorkspaceRepository) Update(ctx context.Context, id string, params workspace.UpdateWorkspaceParams) (workspace.Workspace, error) {
	var model WorkspaceModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return workspace.Workspace{}, workspace.ErrWorkspaceNotFound
		}
		return workspace.Workspace{}, err
	}

	if params.Name != nil {
		var count int64
		if err := r.db.WithContext(ctx).Model(&WorkspaceModel{}).
			Where("name = ? AND created_by_user_id = ? AND id != ?", *params.Name, model.CreatedByUserID, id).
			Count(&count).Error; err != nil {
			return workspace.Workspace{}, err
		}
		if count > 0 {
			return workspace.Workspace{}, workspace.ErrWorkspaceNameTaken
		}
		model.Name = *params.Name
	}

	model.UpdatedAt = time.Now().UTC()
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return workspace.Workspace{}, err
	}
	return toWorkspace(model), nil
}

func (r *WorkspaceRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&WorkspaceModel{}).Where("id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return workspace.ErrWorkspaceNotFound
		}

		var environmentIDs []string
		if err := tx.Model(&EnvironmentModel{}).
			Where("workspace_id = ?", id).
			Pluck("id", &environmentIDs).Error; err != nil {
			return err
		}
		if len(environmentIDs) > 0 {
			if err := tx.Where("environment_id IN ?", environmentIDs).Delete(&VariableModel{}).Error; err != nil {
				return err
			}
		}

		deletes := []struct {
			query string
			value any
			model any
		}{
			{query: "workspace_id = ?", value: id, model: &ResourceGrantModel{}},
			{query: "workspace_id = ?", value: id, model: &RequestModel{}},
			{query: "workspace_id = ?", value: id, model: &FolderModel{}},
			{query: "workspace_id = ?", value: id, model: &EnvironmentModel{}},
			{query: "workspace_id = ?", value: id, model: &WorkspaceMembershipModel{}},
			{query: "id = ?", value: id, model: &WorkspaceModel{}},
		}

		for _, deletion := range deletes {
			if err := tx.Where(deletion.query, deletion.value).Delete(deletion.model).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func toWorkspace(m WorkspaceModel) workspace.Workspace {
	return workspace.Workspace{
		ID:              m.ID,
		Name:            m.Name,
		CreatedByUserID: m.CreatedByUserID,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

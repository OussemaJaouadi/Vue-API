package gormstorage

import (
	"context"
	"time"

	"gorm.io/gorm"

	"vue-api/backend/internal/id"
	"vue-api/backend/internal/workspace"
)

type GrantRepository struct {
	db *gorm.DB
}

func NewGrantRepository(db *gorm.DB) *GrantRepository {
	return &GrantRepository{db: db}
}

func (r *GrantRepository) ListByUserAndWorkspace(ctx context.Context, userID string, workspaceID string) ([]workspace.ResourceGrant, error) {
	var models []ResourceGrantModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND workspace_id = ?", userID, workspaceID).
		Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]workspace.ResourceGrant, len(models))
	for i, m := range models {
		result[i] = toGrant(m)
	}
	return result, nil
}

func (r *GrantRepository) Set(ctx context.Context, userID string, workspaceID string, grants []workspace.GrantInput) ([]workspace.ResourceGrant, error) {
	var result []workspace.ResourceGrant
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? AND workspace_id = ?", userID, workspaceID).
			Delete(&ResourceGrantModel{}).Error; err != nil {
			return err
		}

		now := time.Now().UTC()
		result = make([]workspace.ResourceGrant, 0, len(grants))
		for _, input := range grants {
			model := ResourceGrantModel{
				ID:           id.NewUUIDV7(),
				WorkspaceID:  workspaceID,
				UserID:       userID,
				ResourceType: input.ResourceType,
				ResourceID:   input.ResourceID,
				AccessLevel:  input.AccessLevel,
				CreatedAt:    now,
				UpdatedAt:    now,
			}
			if err := tx.Create(&model).Error; err != nil {
				return err
			}
			result = append(result, toGrant(model))
		}
		return nil
	})
	return result, err
}

func toGrant(m ResourceGrantModel) workspace.ResourceGrant {
	return workspace.ResourceGrant{
		ID:           m.ID,
		WorkspaceID:  m.WorkspaceID,
		UserID:       m.UserID,
		ResourceType: m.ResourceType,
		ResourceID:   m.ResourceID,
		AccessLevel:  m.AccessLevel,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

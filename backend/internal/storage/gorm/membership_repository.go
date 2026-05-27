package gormstorage

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"vue-api/backend/internal/id"
	"vue-api/backend/internal/workspace"
)

type MembershipRepository struct {
	db *gorm.DB
}

func NewMembershipRepository(db *gorm.DB) *MembershipRepository {
	return &MembershipRepository{db: db}
}

func (r *MembershipRepository) ListByWorkspace(ctx context.Context, workspaceID string) ([]workspace.Membership, error) {
	var models []WorkspaceMembershipModel
	if err := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Order("created_at ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]workspace.Membership, len(models))
	for i, m := range models {
		result[i] = toMembership(m)
	}
	return result, nil
}

func (r *MembershipRepository) FindByUserAndWorkspace(ctx context.Context, userID string, workspaceID string) (workspace.Membership, error) {
	var model WorkspaceMembershipModel
	if err := r.db.WithContext(ctx).Where("user_id = ? AND workspace_id = ?", userID, workspaceID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return workspace.Membership{}, workspace.ErrMembershipNotFound
		}
		return workspace.Membership{}, err
	}
	return toMembership(model), nil
}

func (r *MembershipRepository) Create(ctx context.Context, params workspace.CreateMembershipParams) (workspace.Membership, error) {
	var count int64
	r.db.WithContext(ctx).Model(&WorkspaceMembershipModel{}).
		Where("user_id = ? AND workspace_id = ?", params.UserID, params.WorkspaceID).
		Count(&count)
	if count > 0 {
		return workspace.Membership{}, workspace.ErrAlreadyMember
	}

	now := time.Now().UTC()
	model := WorkspaceMembershipModel{
		ID:              id.NewUUIDV7(),
		WorkspaceID:     params.WorkspaceID,
		UserID:          params.UserID,
		Role:            params.Role,
		CreatedByUserID: params.CreatedByUserID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return workspace.Membership{}, err
	}
	return toMembership(model), nil
}

func (r *MembershipRepository) UpdateRole(ctx context.Context, id string, params workspace.UpdateMembershipParams) (workspace.Membership, error) {
	var model WorkspaceMembershipModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return workspace.Membership{}, workspace.ErrMembershipNotFound
		}
		return workspace.Membership{}, err
	}
	if params.Role != nil {
		model.Role = *params.Role
	}
	model.UpdatedAt = time.Now().UTC()
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return workspace.Membership{}, err
	}
	return toMembership(model), nil
}

func (r *MembershipRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&WorkspaceMembershipModel{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return workspace.ErrMembershipNotFound
	}
	return nil
}

func (r *MembershipRepository) CountByWorkspace(ctx context.Context, workspaceID string) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&WorkspaceMembershipModel{}).Where("workspace_id = ?", workspaceID).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func toMembership(m WorkspaceMembershipModel) workspace.Membership {
	return workspace.Membership{
		ID:              m.ID,
		WorkspaceID:     m.WorkspaceID,
		UserID:          m.UserID,
		Role:            m.Role,
		CreatedByUserID: m.CreatedByUserID,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

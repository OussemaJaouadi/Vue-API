package gormstorage

import "time"

type WorkspaceMembershipModel struct {
	ID              string `gorm:"primaryKey;column:id"`
	WorkspaceID     string `gorm:"column:workspace_id;not null"`
	UserID          string `gorm:"column:user_id;not null"`
	Role            string `gorm:"column:role;not null"`
	CreatedByUserID string `gorm:"column:created_by_user_id;not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (WorkspaceMembershipModel) TableName() string {
	return "workspace_memberships"
}

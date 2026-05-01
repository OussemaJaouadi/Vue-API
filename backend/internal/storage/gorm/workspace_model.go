package gormstorage

import "time"

type WorkspaceModel struct {
	ID              string `gorm:"primaryKey;column:id"`
	Name            string `gorm:"column:name;not null"`
	CreatedByUserID string `gorm:"column:created_by_user_id;not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (WorkspaceModel) TableName() string {
	return "workspaces"
}

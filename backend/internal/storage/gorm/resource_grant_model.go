package gormstorage

import "time"

type ResourceGrantModel struct {
	ID            string `gorm:"primaryKey;column:id"`
	WorkspaceID   string `gorm:"column:workspace_id;not null;index"`
	UserID        string `gorm:"column:user_id;not null;index"`
	ResourceType  string `gorm:"column:resource_type;not null"`
	ResourceID    string `gorm:"column:resource_id;not null"`
	AccessLevel   string `gorm:"column:access_level;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (ResourceGrantModel) TableName() string {
	return "resource_grants"
}

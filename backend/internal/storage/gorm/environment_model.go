package gormstorage

import "time"

type EnvironmentModel struct {
	ID              string `gorm:"primaryKey;column:id"`
	WorkspaceID     string `gorm:"column:workspace_id;not null;uniqueIndex:idx_environments_workspace_name"`
	Name            string `gorm:"column:name;not null;uniqueIndex:idx_environments_workspace_name"`
	Visibility      string `gorm:"column:visibility;not null;default:project"`
	SortOrder       int    `gorm:"column:sort_order;not null;default:0"`
	CreatedByUserID string `gorm:"column:created_by_user_id;not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (EnvironmentModel) TableName() string { return "environments" }

type VariableModel struct {
	ID            string `gorm:"primaryKey;column:id"`
	EnvironmentID string `gorm:"column:environment_id;not null;uniqueIndex:idx_variables_env_key"`
	Key           string `gorm:"column:key;not null;uniqueIndex:idx_variables_env_key"`
	Value         string `gorm:"column:value;not null;default:''"`
	Secret        bool   `gorm:"column:secret;not null;default:false"`
	SortOrder     int    `gorm:"column:sort_order;not null;default:0"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (VariableModel) TableName() string { return "environment_variables" }

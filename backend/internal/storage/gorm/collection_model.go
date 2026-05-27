package gormstorage

import "time"

type FolderModel struct {
	ID          string `gorm:"primaryKey;column:id"`
	WorkspaceID string `gorm:"column:workspace_id;not null;uniqueIndex:idx_collections_workspace_name"`
	Name        string `gorm:"column:name;not null;uniqueIndex:idx_collections_workspace_name"`
	Icon        string `gorm:"column:icon;not null;default:PhGlobe"`
	SortOrder   int    `gorm:"column:sort_order;not null;default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (FolderModel) TableName() string { return "collections" }

type RequestModel struct {
	ID           string  `gorm:"primaryKey;column:id"`
	CollectionID *string `gorm:"column:collection_id"`
	WorkspaceID  string  `gorm:"column:workspace_id;not null"`
	Method       string  `gorm:"column:method;not null;default:GET"`
	Name         string  `gorm:"column:name;not null"`
	Path         string  `gorm:"column:path;not null"`
	SortOrder    int     `gorm:"column:sort_order;not null;default:0"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (RequestModel) TableName() string { return "requests" }

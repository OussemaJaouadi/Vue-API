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
	ID              string  `gorm:"primaryKey;column:id"`
	CollectionID    *string `gorm:"column:collection_id"`
	WorkspaceID     string  `gorm:"column:workspace_id;not null"`
	Method          string  `gorm:"column:method;not null;default:GET"`
	Name            string  `gorm:"column:name;not null"`
	Path            string  `gorm:"column:path;not null"`
	QueryParamsJSON string  `gorm:"column:query_params_json;not null;default:'[]'"`
	HeadersJSON     string  `gorm:"column:headers_json;not null;default:'[]'"`
	Body            string  `gorm:"column:body;not null;default:''"`
	BodyLanguage    string  `gorm:"column:body_language;not null;default:json"`
	AuthConfigJSON  string  `gorm:"column:auth_config_json;not null;default:'{}'"`
	SortOrder       int     `gorm:"column:sort_order;not null;default:0"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (RequestModel) TableName() string { return "requests" }

type CollectionEnvironmentPolicyModel struct {
	ID                 string `gorm:"primaryKey;column:id"`
	WorkspaceID        string `gorm:"column:workspace_id;not null;index:idx_collection_env_policies_workspace_collection"`
	CollectionID       string `gorm:"column:collection_id;not null;uniqueIndex:idx_collection_env_policies_collection_environment;index:idx_collection_env_policies_workspace_collection"`
	EnvironmentID      string `gorm:"column:environment_id;not null;uniqueIndex:idx_collection_env_policies_collection_environment"`
	DefaultEnvironment bool   `gorm:"column:default_environment;not null;default:false"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (CollectionEnvironmentPolicyModel) TableName() string {
	return "collection_environment_policies"
}

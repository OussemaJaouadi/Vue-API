package collection

import (
	"context"
	"errors"
	"time"
)

type Folder struct {
	ID          string
	WorkspaceID string
	Name        string
	Icon        string
	SortOrder   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Request struct {
	ID              string
	CollectionID    *string
	WorkspaceID     string
	Method          string
	Name            string
	Path            string
	QueryParamsJSON string
	HeadersJSON     string
	Body            string
	BodyLanguage    string
	AuthConfigJSON  string
	SortOrder       int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type CreateFolderParams struct {
	WorkspaceID string
	Name        string
	Icon        string
}

type UpdateFolderParams struct {
	Name *string
	Icon *string
}

type CreateRequestParams struct {
	CollectionID    *string
	WorkspaceID     string
	Method          string
	Name            string
	Path            string
	QueryParamsJSON string
	HeadersJSON     string
	Body            string
	BodyLanguage    string
	AuthConfigJSON  string
}

type UpdateRequestParams struct {
	Method          *string
	Name            *string
	Path            *string
	QueryParamsJSON *string
	HeadersJSON     *string
	Body            *string
	BodyLanguage    *string
	AuthConfigJSON  *string
}

type RequestOrderGroup struct {
	CollectionID *string
	RequestIDs   []string
}

type FolderWithRequests struct {
	Folder   Folder
	Requests []Request
}

var (
	ErrFolderNotFound  = errors.New("folder not found")
	ErrRequestNotFound = errors.New("request not found")
	ErrFolderNameTaken = errors.New("folder name already exists")
)

type Repository interface {
	ListFolders(ctx context.Context, workspaceID string) ([]Folder, error)
	ListRequests(ctx context.Context, workspaceID string, collectionID *string) ([]Request, error)
	ListRootRequests(ctx context.Context, workspaceID string) ([]Request, error)
	CreateFolder(ctx context.Context, params CreateFolderParams) (Folder, error)
	CreateRequest(ctx context.Context, params CreateRequestParams) (Request, error)
	UpdateFolder(ctx context.Context, id string, params UpdateFolderParams) (Folder, error)
	UpdateRequest(ctx context.Context, id string, params UpdateRequestParams) (Request, error)
	ReorderRequests(ctx context.Context, workspaceID string, groups []RequestOrderGroup) error
	DeleteFolder(ctx context.Context, id string) error
	DeleteRequest(ctx context.Context, id string) error
}

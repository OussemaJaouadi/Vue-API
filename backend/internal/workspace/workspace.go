package workspace

import (
	"context"
	"errors"
	"time"
)

type Workspace struct {
	ID              string
	Name            string
	CreatedByUserID string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Membership struct {
	ID              string
	WorkspaceID     string
	UserID          string
	Role            string
	CreatedByUserID string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ResourceGrant struct {
	ID            string
	WorkspaceID   string
	UserID        string
	ResourceType  string
	ResourceID    string
	AccessLevel   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CreateWorkspaceParams struct {
	Name            string
	CreatedByUserID string
}

type UpdateWorkspaceParams struct {
	Name *string
}

type CreateMembershipParams struct {
	WorkspaceID     string
	UserID          string
	Role            string
	CreatedByUserID string
}

type UpdateMembershipParams struct {
	Role *string
}

type GrantInput struct {
	ResourceType string
	ResourceID   string
	AccessLevel  string
}

var (
	ErrWorkspaceNotFound  = errors.New("workspace not found")
	ErrWorkspaceNameTaken = errors.New("workspace name already taken")
	ErrMembershipNotFound = errors.New("membership not found")
	ErrAlreadyMember      = errors.New("user is already a member of this workspace")
	ErrGrantNotFound      = errors.New("grant not found")
)

type WorkspaceRepository interface {
	ListByUser(ctx context.Context, userID string) ([]Workspace, error)
	Create(ctx context.Context, params CreateWorkspaceParams) (Workspace, error)
	FindByID(ctx context.Context, id string) (Workspace, error)
	Update(ctx context.Context, id string, params UpdateWorkspaceParams) (Workspace, error)
}

type MembershipRepository interface {
	ListByWorkspace(ctx context.Context, workspaceID string) ([]Membership, error)
	FindByUserAndWorkspace(ctx context.Context, userID string, workspaceID string) (Membership, error)
	Create(ctx context.Context, params CreateMembershipParams) (Membership, error)
	UpdateRole(ctx context.Context, id string, params UpdateMembershipParams) (Membership, error)
	Delete(ctx context.Context, id string) error
	CountByWorkspace(ctx context.Context, workspaceID string) (int, error)
}

type GrantRepository interface {
	ListByUserAndWorkspace(ctx context.Context, userID string, workspaceID string) ([]ResourceGrant, error)
	Set(ctx context.Context, userID string, workspaceID string, grants []GrantInput) ([]ResourceGrant, error)
}

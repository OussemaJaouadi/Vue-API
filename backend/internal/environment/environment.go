package environment

import (
	"context"
	"errors"
	"time"
)

type Environment struct {
	ID              string
	WorkspaceID     string
	Name            string
	Visibility      string
	SortOrder       int
	CreatedByUserID string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Variable struct {
	ID            string
	EnvironmentID string
	Key           string
	Value         string
	Secret        bool
	SortOrder     int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CreateEnvironmentParams struct {
	WorkspaceID     string
	Name            string
	Visibility      string
	CreatedByUserID string
}

type UpdateEnvironmentParams struct {
	Name       *string
	Visibility *string
}

type CreateVariableParams struct {
	EnvironmentID string
	Key           string
	Value         string
	Secret        bool
}

type UpdateVariableParams struct {
	Key    *string
	Value  *string
	Secret *bool
}

type EnvironmentWithVariables struct {
	Environment Environment
	Variables   []Variable
}

var (
	ErrEnvironmentNotFound     = errors.New("environment not found")
	ErrEnvironmentNameTaken    = errors.New("environment name already exists")
	ErrVariableNotFound        = errors.New("variable not found")
	ErrVariableKeyTaken        = errors.New("variable key already exists in this environment")
)

type Repository interface {
	ListEnvironments(ctx context.Context, workspaceID string) ([]Environment, error)
	ListVariables(ctx context.Context, environmentID string) ([]Variable, error)
	CreateEnvironment(ctx context.Context, params CreateEnvironmentParams) (Environment, error)
	CreateVariable(ctx context.Context, params CreateVariableParams) (Variable, error)
	UpdateEnvironment(ctx context.Context, id string, params UpdateEnvironmentParams) (Environment, error)
	UpdateVariable(ctx context.Context, id string, params UpdateVariableParams) (Variable, error)
	DeleteEnvironment(ctx context.Context, id string) error
	DeleteVariable(ctx context.Context, id string) error
}

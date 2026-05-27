package environment_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"vue-api/backend/internal/environment"
)

func TestErrorSentinels_AreDistinct(t *testing.T) {
	assert.ErrorIs(t, environment.ErrEnvironmentNotFound, environment.ErrEnvironmentNotFound)
	assert.ErrorIs(t, environment.ErrEnvironmentNameTaken, environment.ErrEnvironmentNameTaken)
	assert.ErrorIs(t, environment.ErrVariableNotFound, environment.ErrVariableNotFound)
	assert.ErrorIs(t, environment.ErrVariableKeyTaken, environment.ErrVariableKeyTaken)

	assert.False(t, errors.Is(environment.ErrEnvironmentNotFound, environment.ErrEnvironmentNameTaken))
	assert.False(t, errors.Is(environment.ErrEnvironmentNotFound, environment.ErrVariableNotFound))
	assert.False(t, errors.Is(environment.ErrEnvironmentNameTaken, environment.ErrVariableKeyTaken))
	assert.False(t, errors.Is(environment.ErrVariableNotFound, environment.ErrVariableKeyTaken))
}

func TestCreateEnvironmentParams_HoldsValues(t *testing.T) {
	params := environment.CreateEnvironmentParams{
		WorkspaceID:     "ws1",
		Name:            "Production",
		Visibility:      "project",
		CreatedByUserID: "user1",
	}

	assert.Equal(t, "ws1", params.WorkspaceID)
	assert.Equal(t, "Production", params.Name)
	assert.Equal(t, "project", params.Visibility)
	assert.Equal(t, "user1", params.CreatedByUserID)
}

func TestUpdateEnvironmentParams_PartialUpdates(t *testing.T) {
	name := "Staging"
	vis := "personal"
	params := environment.UpdateEnvironmentParams{
		Name:       &name,
		Visibility: &vis,
	}

	assert.Equal(t, "Staging", *params.Name)
	assert.Equal(t, "personal", *params.Visibility)
}

func TestUpdateEnvironmentParams_NilFields(t *testing.T) {
	params := environment.UpdateEnvironmentParams{}

	assert.Nil(t, params.Name)
	assert.Nil(t, params.Visibility)
}

func TestCreateVariableParams_HoldsValues(t *testing.T) {
	params := environment.CreateVariableParams{
		EnvironmentID: "env1",
		Key:           "API_KEY",
		Value:         "secret123",
		Secret:        true,
	}

	assert.Equal(t, "env1", params.EnvironmentID)
	assert.Equal(t, "API_KEY", params.Key)
	assert.Equal(t, "secret123", params.Value)
	assert.True(t, params.Secret)
}

func TestUpdateVariableParams_PartialUpdates(t *testing.T) {
	key := "NEW_KEY"
	value := "new_val"
	secret := false
	params := environment.UpdateVariableParams{
		Key:    &key,
		Value:  &value,
		Secret: &secret,
	}

	assert.Equal(t, "NEW_KEY", *params.Key)
	assert.Equal(t, "new_val", *params.Value)
	assert.False(t, *params.Secret)
}

func TestUpdateVariableParams_NilFields(t *testing.T) {
	params := environment.UpdateVariableParams{}

	assert.Nil(t, params.Key)
	assert.Nil(t, params.Value)
	assert.Nil(t, params.Secret)
}

func TestEnvironmentWithVariables_HoldsValues(t *testing.T) {
	ev := environment.EnvironmentWithVariables{
		Environment: environment.Environment{
			ID:   "env1",
			Name: "Prod",
		},
		Variables: []environment.Variable{
			{ID: "v1", Key: "KEY", Value: "val"},
		},
	}

	assert.Equal(t, "env1", ev.Environment.ID)
	assert.Len(t, ev.Variables, 1)
	assert.Equal(t, "KEY", ev.Variables[0].Key)
}

func TestEnvironment_DefaultZeroValues(t *testing.T) {
	var e environment.Environment

	assert.Empty(t, e.ID)
	assert.Empty(t, e.WorkspaceID)
	assert.Empty(t, e.Name)
	assert.Empty(t, e.Visibility)
	assert.Equal(t, 0, e.SortOrder)
	assert.Empty(t, e.CreatedByUserID)
	assert.True(t, e.CreatedAt.IsZero())
	assert.True(t, e.UpdatedAt.IsZero())
}

func TestVariable_DefaultZeroValues(t *testing.T) {
	var v environment.Variable

	assert.Empty(t, v.ID)
	assert.Empty(t, v.EnvironmentID)
	assert.Empty(t, v.Key)
	assert.Empty(t, v.Value)
	assert.False(t, v.Secret)
	assert.Equal(t, 0, v.SortOrder)
	assert.True(t, v.CreatedAt.IsZero())
	assert.True(t, v.UpdatedAt.IsZero())
}

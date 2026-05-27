package gormstorage_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vue-api/backend/internal/environment"
	gormstorage "vue-api/backend/internal/storage/gorm"
)

func newEnvRepo(t *testing.T) *gormstorage.EnvironmentRepository {
	t.Helper()
	db := openTestDB(t)
	require.NoError(t, gormstorage.Migrate(db))
	return gormstorage.NewEnvironmentRepository(db)
}

func TestCreateEnvironment_ReturnsEnvWithFields(t *testing.T) {
	repo := newEnvRepo(t)

	env, err := repo.CreateEnvironment(context.Background(), environment.CreateEnvironmentParams{
		WorkspaceID:     "ws1",
		Name:            "Production",
		Visibility:      "project",
		CreatedByUserID: "user1",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, env.ID)
	assert.Equal(t, "ws1", env.WorkspaceID)
	assert.Equal(t, "Production", env.Name)
	assert.Equal(t, "project", env.Visibility)
	assert.Equal(t, "user1", env.CreatedByUserID)
	assert.False(t, env.CreatedAt.IsZero())
	assert.False(t, env.UpdatedAt.IsZero())
}

func TestCreateEnvironment_DefaultVisibility(t *testing.T) {
	repo := newEnvRepo(t)

	env, err := repo.CreateEnvironment(context.Background(), environment.CreateEnvironmentParams{
		WorkspaceID:     "ws1",
		Name:            "Dev",
		CreatedByUserID: "user1",
	})

	require.NoError(t, err)
	assert.Equal(t, "project", env.Visibility)
}

func TestCreateEnvironment_DuplicateNameReturnsError(t *testing.T) {
	repo := newEnvRepo(t)

	_, err := repo.CreateEnvironment(context.Background(), environment.CreateEnvironmentParams{
		WorkspaceID: "ws1", Name: "Staging", CreatedByUserID: "user1",
	})
	require.NoError(t, err)

	_, err = repo.CreateEnvironment(context.Background(), environment.CreateEnvironmentParams{
		WorkspaceID: "ws1", Name: "Staging", CreatedByUserID: "user1",
	})
	require.ErrorIs(t, err, environment.ErrEnvironmentNameTaken)
}

func TestCreateEnvironment_SameNameDifferentWorkspaceIsOK(t *testing.T) {
	repo := newEnvRepo(t)

	_, err := repo.CreateEnvironment(context.Background(), environment.CreateEnvironmentParams{
		WorkspaceID: "ws1", Name: "Staging", CreatedByUserID: "user1",
	})
	require.NoError(t, err)

	_, err = repo.CreateEnvironment(context.Background(), environment.CreateEnvironmentParams{
		WorkspaceID: "ws2", Name: "Staging", CreatedByUserID: "user2",
	})
	require.NoError(t, err)
}

func TestListEnvironments_ReturnsAllForWorkspace(t *testing.T) {
	repo := newEnvRepo(t)

	_, _ = repo.CreateEnvironment(context.Background(), environment.CreateEnvironmentParams{
		WorkspaceID: "ws1", Name: "Prod", CreatedByUserID: "user1",
	})
	_, _ = repo.CreateEnvironment(context.Background(), environment.CreateEnvironmentParams{
		WorkspaceID: "ws1", Name: "Dev", CreatedByUserID: "user1",
	})

	envs, err := repo.ListEnvironments(context.Background(), "ws1")

	require.NoError(t, err)
	assert.Len(t, envs, 2)
}

func TestListEnvironments_FiltersByWorkspace(t *testing.T) {
	repo := newEnvRepo(t)

	_, _ = repo.CreateEnvironment(context.Background(), environment.CreateEnvironmentParams{
		WorkspaceID: "ws1", Name: "Prod", CreatedByUserID: "user1",
	})
	_, _ = repo.CreateEnvironment(context.Background(), environment.CreateEnvironmentParams{
		WorkspaceID: "ws2", Name: "Other", CreatedByUserID: "user2",
	})

	envs, err := repo.ListEnvironments(context.Background(), "ws1")
	require.NoError(t, err)
	assert.Len(t, envs, 1)
	assert.Equal(t, "Prod", envs[0].Name)
}

func TestListEnvironments_EmptyWorkspace(t *testing.T) {
	repo := newEnvRepo(t)

	envs, err := repo.ListEnvironments(context.Background(), "empty")
	require.NoError(t, err)
	assert.Empty(t, envs)
}

func TestUpdateEnvironment_NameAndVisibility(t *testing.T) {
	repo := newEnvRepo(t)

	env, _ := repo.CreateEnvironment(context.Background(), environment.CreateEnvironmentParams{
		WorkspaceID: "ws1", Name: "Old", CreatedByUserID: "user1",
	})

	newName := "Renamed"
	newVis := "personal"
	updated, err := repo.UpdateEnvironment(context.Background(), env.ID, environment.UpdateEnvironmentParams{
		Name:       &newName,
		Visibility: &newVis,
	})

	require.NoError(t, err)
	assert.Equal(t, "Renamed", updated.Name)
	assert.Equal(t, "personal", updated.Visibility)
}

func TestUpdateEnvironment_NotFound(t *testing.T) {
	repo := newEnvRepo(t)

	name := "Nope"
	_, err := repo.UpdateEnvironment(context.Background(), "nonexistent", environment.UpdateEnvironmentParams{
		Name: &name,
	})

	require.ErrorIs(t, err, environment.ErrEnvironmentNotFound)
}

func TestDeleteEnvironment_RemovesAndCascadesVariables(t *testing.T) {
	repo := newEnvRepo(t)

	env, _ := repo.CreateEnvironment(context.Background(), environment.CreateEnvironmentParams{
		WorkspaceID: "ws1", Name: "Temp", CreatedByUserID: "user1",
	})
	variable, _ := repo.CreateVariable(context.Background(), environment.CreateVariableParams{
		EnvironmentID: env.ID,
		Key:           "API_KEY",
		Value:         "secret",
		Secret:        true,
	})

	err := repo.DeleteEnvironment(context.Background(), env.ID)
	require.NoError(t, err)

	envs, _ := repo.ListEnvironments(context.Background(), "ws1")
	assert.Empty(t, envs)

	err = repo.DeleteVariable(context.Background(), variable.ID)
	require.ErrorIs(t, err, environment.ErrVariableNotFound)
}

func TestDeleteEnvironment_NotFound(t *testing.T) {
	repo := newEnvRepo(t)

	err := repo.DeleteEnvironment(context.Background(), "nonexistent")
	require.ErrorIs(t, err, environment.ErrEnvironmentNotFound)
}

func TestCreateVariable_ReturnsVariableWithFields(t *testing.T) {
	repo := newEnvRepo(t)

	env, _ := repo.CreateEnvironment(context.Background(), environment.CreateEnvironmentParams{
		WorkspaceID: "ws1", Name: "Test", CreatedByUserID: "user1",
	})

	v, err := repo.CreateVariable(context.Background(), environment.CreateVariableParams{
		EnvironmentID: env.ID,
		Key:           "DB_HOST",
		Value:         "localhost",
		Secret:        false,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, v.ID)
	assert.Equal(t, env.ID, v.EnvironmentID)
	assert.Equal(t, "DB_HOST", v.Key)
	assert.Equal(t, "localhost", v.Value)
	assert.False(t, v.Secret)
}

func TestCreateVariable_Secret(t *testing.T) {
	repo := newEnvRepo(t)

	env, _ := repo.CreateEnvironment(context.Background(), environment.CreateEnvironmentParams{
		WorkspaceID: "ws1", Name: "Secrets", CreatedByUserID: "user1",
	})

	v, err := repo.CreateVariable(context.Background(), environment.CreateVariableParams{
		EnvironmentID: env.ID,
		Key:           "PASSWORD",
		Value:         "hunter2",
		Secret:        true,
	})

	require.NoError(t, err)
	assert.True(t, v.Secret)
}

func TestCreateVariable_DuplicateKeyReturnsError(t *testing.T) {
	repo := newEnvRepo(t)

	env, _ := repo.CreateEnvironment(context.Background(), environment.CreateEnvironmentParams{
		WorkspaceID: "ws1", Name: "Test", CreatedByUserID: "user1",
	})

	_, err := repo.CreateVariable(context.Background(), environment.CreateVariableParams{
		EnvironmentID: env.ID, Key: "KEY", Value: "val1",
	})
	require.NoError(t, err)

	_, err = repo.CreateVariable(context.Background(), environment.CreateVariableParams{
		EnvironmentID: env.ID, Key: "KEY", Value: "val2",
	})
	require.ErrorIs(t, err, environment.ErrVariableKeyTaken)
}

func TestCreateVariable_SameKeyDifferentEnvIsOK(t *testing.T) {
	repo := newEnvRepo(t)

	env1, _ := repo.CreateEnvironment(context.Background(), environment.CreateEnvironmentParams{
		WorkspaceID: "ws1", Name: "Dev", CreatedByUserID: "user1",
	})
	env2, _ := repo.CreateEnvironment(context.Background(), environment.CreateEnvironmentParams{
		WorkspaceID: "ws1", Name: "Prod", CreatedByUserID: "user1",
	})

	_, err := repo.CreateVariable(context.Background(), environment.CreateVariableParams{
		EnvironmentID: env1.ID, Key: "DB_URL", Value: "dev",
	})
	require.NoError(t, err)

	_, err = repo.CreateVariable(context.Background(), environment.CreateVariableParams{
		EnvironmentID: env2.ID, Key: "DB_URL", Value: "prod",
	})
	require.NoError(t, err)
}

func TestListVariables_ReturnsAllForEnvironment(t *testing.T) {
	repo := newEnvRepo(t)

	env, _ := repo.CreateEnvironment(context.Background(), environment.CreateEnvironmentParams{
		WorkspaceID: "ws1", Name: "Test", CreatedByUserID: "user1",
	})

	_, _ = repo.CreateVariable(context.Background(), environment.CreateVariableParams{
		EnvironmentID: env.ID, Key: "A", Value: "1",
	})
	_, _ = repo.CreateVariable(context.Background(), environment.CreateVariableParams{
		EnvironmentID: env.ID, Key: "B", Value: "2",
	})

	vars, err := repo.ListVariables(context.Background(), env.ID)

	require.NoError(t, err)
	assert.Len(t, vars, 2)
}

func TestListVariables_EmptyEnvironment(t *testing.T) {
	repo := newEnvRepo(t)

	env, _ := repo.CreateEnvironment(context.Background(), environment.CreateEnvironmentParams{
		WorkspaceID: "ws1", Name: "Empty", CreatedByUserID: "user1",
	})

	vars, err := repo.ListVariables(context.Background(), env.ID)
	require.NoError(t, err)
	assert.Empty(t, vars)
}

func TestUpdateVariable_KeyValueSecret(t *testing.T) {
	repo := newEnvRepo(t)

	env, _ := repo.CreateEnvironment(context.Background(), environment.CreateEnvironmentParams{
		WorkspaceID: "ws1", Name: "Test", CreatedByUserID: "user1",
	})
	v, _ := repo.CreateVariable(context.Background(), environment.CreateVariableParams{
		EnvironmentID: env.ID, Key: "OLD", Value: "old", Secret: false,
	})

	newKey := "NEW"
	newVal := "new"
	newSecret := true
	updated, err := repo.UpdateVariable(context.Background(), v.ID, environment.UpdateVariableParams{
		Key:    &newKey,
		Value:  &newVal,
		Secret: &newSecret,
	})

	require.NoError(t, err)
	assert.Equal(t, "NEW", updated.Key)
	assert.Equal(t, "new", updated.Value)
	assert.True(t, updated.Secret)
}

func TestUpdateVariable_PartialUpdate(t *testing.T) {
	repo := newEnvRepo(t)

	env, _ := repo.CreateEnvironment(context.Background(), environment.CreateEnvironmentParams{
		WorkspaceID: "ws1", Name: "Test", CreatedByUserID: "user1",
	})
	v, _ := repo.CreateVariable(context.Background(), environment.CreateVariableParams{
		EnvironmentID: env.ID, Key: "KEY", Value: "val", Secret: false,
	})

	newVal := "updated"
	updated, err := repo.UpdateVariable(context.Background(), v.ID, environment.UpdateVariableParams{
		Value: &newVal,
	})

	require.NoError(t, err)
	assert.Equal(t, "KEY", updated.Key)
	assert.Equal(t, "updated", updated.Value)
	assert.False(t, updated.Secret)
}

func TestUpdateVariable_NotFound(t *testing.T) {
	repo := newEnvRepo(t)

	val := "nope"
	_, err := repo.UpdateVariable(context.Background(), "nonexistent", environment.UpdateVariableParams{
		Value: &val,
	})

	require.ErrorIs(t, err, environment.ErrVariableNotFound)
}

func TestDeleteVariable_RemovesVariable(t *testing.T) {
	repo := newEnvRepo(t)

	env, _ := repo.CreateEnvironment(context.Background(), environment.CreateEnvironmentParams{
		WorkspaceID: "ws1", Name: "Test", CreatedByUserID: "user1",
	})
	v, _ := repo.CreateVariable(context.Background(), environment.CreateVariableParams{
		EnvironmentID: env.ID, Key: "TEMP", Value: "x",
	})

	err := repo.DeleteVariable(context.Background(), v.ID)
	require.NoError(t, err)

	vars, _ := repo.ListVariables(context.Background(), env.ID)
	assert.Empty(t, vars)
}

func TestDeleteVariable_NotFound(t *testing.T) {
	repo := newEnvRepo(t)

	err := repo.DeleteVariable(context.Background(), "nonexistent")
	require.ErrorIs(t, err, environment.ErrVariableNotFound)
}

func TestCreateVariable_DefaultSortOrder(t *testing.T) {
	repo := newEnvRepo(t)

	env, _ := repo.CreateEnvironment(context.Background(), environment.CreateEnvironmentParams{
		WorkspaceID: "ws1", Name: "Test", CreatedByUserID: "user1",
	})

	v, err := repo.CreateVariable(context.Background(), environment.CreateVariableParams{
		EnvironmentID: env.ID, Key: "K", Value: "v",
	})

	require.NoError(t, err)
	assert.Equal(t, 0, v.SortOrder)
}

func TestCreateEnvironment_DefaultSortOrder(t *testing.T) {
	repo := newEnvRepo(t)

	env, err := repo.CreateEnvironment(context.Background(), environment.CreateEnvironmentParams{
		WorkspaceID: "ws1", Name: "Test", CreatedByUserID: "user1",
	})

	require.NoError(t, err)
	assert.Equal(t, 0, env.SortOrder)
}

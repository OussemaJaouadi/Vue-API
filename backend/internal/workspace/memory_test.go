package workspace_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vue-api/backend/internal/workspace"
)

func TestMemoryWorkspaceRepository_CreateAndFindByID(t *testing.T) {
	repo := workspace.NewMemoryWorkspaceRepository()
	ctx := context.Background()

	ws, err := repo.Create(ctx, workspace.CreateWorkspaceParams{
		Name:            "My Workspace",
		CreatedByUserID: "user1",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, ws.ID)
	assert.Equal(t, "My Workspace", ws.Name)

	found, err := repo.FindByID(ctx, ws.ID)
	require.NoError(t, err)
	assert.Equal(t, ws.Name, found.Name)
}

func TestMemoryWorkspaceRepository_Create_DuplicateName(t *testing.T) {
	repo := workspace.NewMemoryWorkspaceRepository()
	ctx := context.Background()

	_, err := repo.Create(ctx, workspace.CreateWorkspaceParams{
		Name: "My Workspace", CreatedByUserID: "user1",
	})
	require.NoError(t, err)

	_, err = repo.Create(ctx, workspace.CreateWorkspaceParams{
		Name: "My Workspace", CreatedByUserID: "user1",
	})
	assert.ErrorIs(t, err, workspace.ErrWorkspaceNameTaken)
}

func TestMemoryWorkspaceRepository_FindByID_NotFound(t *testing.T) {
	repo := workspace.NewMemoryWorkspaceRepository()
	_, err := repo.FindByID(context.Background(), "nonexistent")
	assert.ErrorIs(t, err, workspace.ErrWorkspaceNotFound)
}

func TestMemoryWorkspaceRepository_ListByUser(t *testing.T) {
	repo := workspace.NewMemoryWorkspaceRepository()
	ctx := context.Background()

	repo.Create(ctx, workspace.CreateWorkspaceParams{Name: "A", CreatedByUserID: "user1"})
	repo.Create(ctx, workspace.CreateWorkspaceParams{Name: "B", CreatedByUserID: "user2"})
	repo.Create(ctx, workspace.CreateWorkspaceParams{Name: "C", CreatedByUserID: "user1"})

	workspaces, err := repo.ListByUser(ctx, "user1")
	require.NoError(t, err)
	assert.Len(t, workspaces, 2)
}

func TestMemoryWorkspaceRepository_Update(t *testing.T) {
	repo := workspace.NewMemoryWorkspaceRepository()
	ctx := context.Background()

	ws, _ := repo.Create(ctx, workspace.CreateWorkspaceParams{Name: "Old", CreatedByUserID: "user1"})

	name := "New"
	updated, err := repo.Update(ctx, ws.ID, workspace.UpdateWorkspaceParams{Name: &name})
	require.NoError(t, err)
	assert.Equal(t, "New", updated.Name)
}

func TestMemoryWorkspaceRepository_Update_NotFound(t *testing.T) {
	repo := workspace.NewMemoryWorkspaceRepository()
	name := "X"
	_, err := repo.Update(context.Background(), "nope", workspace.UpdateWorkspaceParams{Name: &name})
	assert.ErrorIs(t, err, workspace.ErrWorkspaceNotFound)
}

func TestMemoryWorkspaceRepository_ListByUser_NoWorkspaces(t *testing.T) {
	repo := workspace.NewMemoryWorkspaceRepository()
	workspaces, err := repo.ListByUser(context.Background(), "user1")
	require.NoError(t, err)
	assert.Empty(t, workspaces)
}

func TestMemoryWorkspaceRepository_DifferentUserSameName(t *testing.T) {
	repo := workspace.NewMemoryWorkspaceRepository()
	ctx := context.Background()

	_, err := repo.Create(ctx, workspace.CreateWorkspaceParams{Name: "My WS", CreatedByUserID: "user1"})
	require.NoError(t, err)

	_, err = repo.Create(ctx, workspace.CreateWorkspaceParams{Name: "My WS", CreatedByUserID: "user2"})
	require.NoError(t, err)

	workspaces, _ := repo.ListByUser(ctx, "user1")
	assert.Len(t, workspaces, 1)
	workspaces2, _ := repo.ListByUser(ctx, "user2")
	assert.Len(t, workspaces2, 1)
}

func TestMemoryMembershipRepository_CreateAndList(t *testing.T) {
	repo := workspace.NewMemoryMembershipRepository()
	ctx := context.Background()

	m, err := repo.Create(ctx, workspace.CreateMembershipParams{
		WorkspaceID: "ws1", UserID: "user1", Role: "admin", CreatedByUserID: "user1",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, m.ID)

	members, _ := repo.ListByWorkspace(ctx, "ws1")
	assert.Len(t, members, 1)
}

func TestMemoryMembershipRepository_Create_Duplicate(t *testing.T) {
	repo := workspace.NewMemoryMembershipRepository()
	ctx := context.Background()

	repo.Create(ctx, workspace.CreateMembershipParams{WorkspaceID: "ws1", UserID: "user1", Role: "admin"})
	_, err := repo.Create(ctx, workspace.CreateMembershipParams{WorkspaceID: "ws1", UserID: "user1", Role: "developer"})
	assert.ErrorIs(t, err, workspace.ErrAlreadyMember)
}

func TestMemoryMembershipRepository_FindByUserAndWorkspace(t *testing.T) {
	repo := workspace.NewMemoryMembershipRepository()
	ctx := context.Background()

	repo.Create(ctx, workspace.CreateMembershipParams{WorkspaceID: "ws1", UserID: "user1", Role: "admin"})
	m, err := repo.FindByUserAndWorkspace(ctx, "user1", "ws1")
	require.NoError(t, err)
	assert.Equal(t, "admin", m.Role)

	_, err = repo.FindByUserAndWorkspace(ctx, "user2", "ws1")
	assert.ErrorIs(t, err, workspace.ErrMembershipNotFound)
}

func TestMemoryMembershipRepository_UpdateRole(t *testing.T) {
	repo := workspace.NewMemoryMembershipRepository()
	ctx := context.Background()

	m, _ := repo.Create(ctx, workspace.CreateMembershipParams{WorkspaceID: "ws1", UserID: "user1", Role: "admin"})
	role := "developer"
	updated, err := repo.UpdateRole(ctx, m.ID, workspace.UpdateMembershipParams{Role: &role})
	require.NoError(t, err)
	assert.Equal(t, "developer", updated.Role)
}

func TestMemoryMembershipRepository_Delete(t *testing.T) {
	repo := workspace.NewMemoryMembershipRepository()
	ctx := context.Background()

	m, _ := repo.Create(ctx, workspace.CreateMembershipParams{WorkspaceID: "ws1", UserID: "user1", Role: "admin"})
	err := repo.Delete(ctx, m.ID)
	require.NoError(t, err)

	_, err = repo.FindByUserAndWorkspace(ctx, "user1", "ws1")
	assert.ErrorIs(t, err, workspace.ErrMembershipNotFound)
}

func TestMemoryMembershipRepository_ListByWorkspace_Scoped(t *testing.T) {
	repo := workspace.NewMemoryMembershipRepository()
	ctx := context.Background()

	repo.Create(ctx, workspace.CreateMembershipParams{WorkspaceID: "ws1", UserID: "user1", Role: "admin"})
	repo.Create(ctx, workspace.CreateMembershipParams{WorkspaceID: "ws1", UserID: "user2", Role: "developer"})
	repo.Create(ctx, workspace.CreateMembershipParams{WorkspaceID: "ws2", UserID: "user3", Role: "admin"})

	members, _ := repo.ListByWorkspace(ctx, "ws1")
	assert.Len(t, members, 2)
}

func TestMemoryGrantRepository_SetAndList(t *testing.T) {
	repo := workspace.NewMemoryGrantRepository()
	ctx := context.Background()

	grants, err := repo.Set(ctx, "user1", "ws1", []workspace.GrantInput{
		{ResourceType: "collection", ResourceID: "c1", AccessLevel: "read"},
		{ResourceType: "environment", ResourceID: "e1", AccessLevel: "admin"},
	})
	require.NoError(t, err)
	assert.Len(t, grants, 2)

	listed, _ := repo.ListByUserAndWorkspace(ctx, "user1", "ws1")
	assert.Len(t, listed, 2)
}

func TestMemoryGrantRepository_Set_ReplacesOld(t *testing.T) {
	repo := workspace.NewMemoryGrantRepository()
	ctx := context.Background()

	repo.Set(ctx, "user1", "ws1", []workspace.GrantInput{
		{ResourceType: "collection", ResourceID: "c1", AccessLevel: "read"},
	})
	repo.Set(ctx, "user1", "ws1", []workspace.GrantInput{
		{ResourceType: "collection", ResourceID: "c1", AccessLevel: "admin"},
	})

	listed, _ := repo.ListByUserAndWorkspace(ctx, "user1", "ws1")
	require.Len(t, listed, 1)
	assert.Equal(t, "admin", listed[0].AccessLevel)
}

func TestMemoryGrantRepository_ListByUserAndWorkspace_Isolated(t *testing.T) {
	repo := workspace.NewMemoryGrantRepository()
	ctx := context.Background()

	repo.Set(ctx, "user1", "ws1", []workspace.GrantInput{
		{ResourceType: "collection", ResourceID: "c1", AccessLevel: "read"},
	})
	repo.Set(ctx, "user2", "ws1", []workspace.GrantInput{
		{ResourceType: "collection", ResourceID: "c1", AccessLevel: "admin"},
	})
	repo.Set(ctx, "user1", "ws2", []workspace.GrantInput{
		{ResourceType: "collection", ResourceID: "c1", AccessLevel: "write"},
	})

	listed1, _ := repo.ListByUserAndWorkspace(ctx, "user1", "ws1")
	assert.Len(t, listed1, 1)

	listed2, _ := repo.ListByUserAndWorkspace(ctx, "user2", "ws1")
	assert.Len(t, listed2, 1)

	listed3, _ := repo.ListByUserAndWorkspace(ctx, "user1", "ws2")
	assert.Len(t, listed3, 1)
}

func TestMemoryGrantRepository_List_Empty(t *testing.T) {
	repo := workspace.NewMemoryGrantRepository()
	grants, err := repo.ListByUserAndWorkspace(context.Background(), "user1", "ws1")
	require.NoError(t, err)
	assert.Empty(t, grants)
}

func TestMemoryGrantRepository_Set_Empty(t *testing.T) {
	repo := workspace.NewMemoryGrantRepository()
	grants, err := repo.Set(context.Background(), "user1", "ws1", nil)
	require.NoError(t, err)
	assert.Empty(t, grants)
}

package gormstorage_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vue-api/backend/internal/workspace"
	gormstorage "vue-api/backend/internal/storage/gorm"
)

func TestGORMMembershipRepository_CreateAndList(t *testing.T) {
	db := workspaceTestDB(t)
	repo := gormstorage.NewMembershipRepository(db)
	ctx := context.Background()

	m, err := repo.Create(ctx, workspace.CreateMembershipParams{
		WorkspaceID: "ws1", UserID: "user1", Role: "admin", CreatedByUserID: "user1",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, m.ID)

	list, _ := repo.ListByWorkspace(ctx, "ws1")
	assert.Len(t, list, 1)
	assert.Equal(t, "user1", list[0].UserID)
}

func TestGORMMembershipRepository_Create_Duplicate(t *testing.T) {
	db := workspaceTestDB(t)
	repo := gormstorage.NewMembershipRepository(db)
	ctx := context.Background()

	repo.Create(ctx, workspace.CreateMembershipParams{WorkspaceID: "ws1", UserID: "user1", Role: "admin"})
	_, err := repo.Create(ctx, workspace.CreateMembershipParams{WorkspaceID: "ws1", UserID: "user1", Role: "developer"})
	assert.ErrorIs(t, err, workspace.ErrAlreadyMember)
}

func TestGORMMembershipRepository_FindByUserAndWorkspace(t *testing.T) {
	db := workspaceTestDB(t)
	repo := gormstorage.NewMembershipRepository(db)
	ctx := context.Background()

	repo.Create(ctx, workspace.CreateMembershipParams{WorkspaceID: "ws1", UserID: "user1", Role: "admin"})

	m, err := repo.FindByUserAndWorkspace(ctx, "user1", "ws1")
	require.NoError(t, err)
	assert.Equal(t, "admin", m.Role)

	_, err = repo.FindByUserAndWorkspace(ctx, "user2", "ws1")
	assert.ErrorIs(t, err, workspace.ErrMembershipNotFound)
}

func TestGORMMembershipRepository_UpdateRole(t *testing.T) {
	db := workspaceTestDB(t)
	repo := gormstorage.NewMembershipRepository(db)
	ctx := context.Background()

	m, _ := repo.Create(ctx, workspace.CreateMembershipParams{WorkspaceID: "ws1", UserID: "user1", Role: "admin"})
	role := "developer"
	updated, err := repo.UpdateRole(ctx, m.ID, workspace.UpdateMembershipParams{Role: &role})
	require.NoError(t, err)
	assert.Equal(t, "developer", updated.Role)
}

func TestGORMMembershipRepository_Delete(t *testing.T) {
	db := workspaceTestDB(t)
	repo := gormstorage.NewMembershipRepository(db)
	ctx := context.Background()

	m, _ := repo.Create(ctx, workspace.CreateMembershipParams{WorkspaceID: "ws1", UserID: "user1", Role: "admin"})
	err := repo.Delete(ctx, m.ID)
	require.NoError(t, err)

	_, err = repo.FindByUserAndWorkspace(ctx, "user1", "ws1")
	assert.ErrorIs(t, err, workspace.ErrMembershipNotFound)
}

func TestGORMMembershipRepository_ListByWorkspace_Scoped(t *testing.T) {
	db := workspaceTestDB(t)
	repo := gormstorage.NewMembershipRepository(db)
	ctx := context.Background()

	repo.Create(ctx, workspace.CreateMembershipParams{WorkspaceID: "ws1", UserID: "user1", Role: "admin"})
	repo.Create(ctx, workspace.CreateMembershipParams{WorkspaceID: "ws1", UserID: "user2", Role: "developer"})
	repo.Create(ctx, workspace.CreateMembershipParams{WorkspaceID: "ws2", UserID: "user3", Role: "admin"})

	list, _ := repo.ListByWorkspace(ctx, "ws1")
	assert.Len(t, list, 2)
}

func TestGORMMembershipRepository_CountByWorkspace(t *testing.T) {
	db := workspaceTestDB(t)
	repo := gormstorage.NewMembershipRepository(db)
	ctx := context.Background()

	repo.Create(ctx, workspace.CreateMembershipParams{WorkspaceID: "ws1", UserID: "user1", Role: "admin"})
	repo.Create(ctx, workspace.CreateMembershipParams{WorkspaceID: "ws1", UserID: "user2", Role: "developer"})

	count, _ := repo.CountByWorkspace(ctx, "ws1")
	assert.Equal(t, 2, count)
}

func TestGORMMembershipRepository_Delete_NotFound(t *testing.T) {
	db := workspaceTestDB(t)
	err := gormstorage.NewMembershipRepository(db).Delete(context.Background(), "nope")
	assert.ErrorIs(t, err, workspace.ErrMembershipNotFound)
}

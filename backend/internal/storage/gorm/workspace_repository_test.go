package gormstorage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"vue-api/backend/internal/workspace"
	gormstorage "vue-api/backend/internal/storage/gorm"
)

func workspaceTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=private", t.Name())), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, gormstorage.Migrate(db))
	return db
}

func TestGORMWorkspaceRepository_CreateAndFindByID(t *testing.T) {
	db := workspaceTestDB(t)
	repo := gormstorage.NewWorkspaceRepository(db)
	ctx := context.Background()

	ws, err := repo.Create(ctx, workspace.CreateWorkspaceParams{
		Name: "My WS", CreatedByUserID: "user1",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, ws.ID)
	assert.Equal(t, "My WS", ws.Name)

	found, err := repo.FindByID(ctx, ws.ID)
	require.NoError(t, err)
	assert.Equal(t, ws.Name, found.Name)
}

func TestGORMWorkspaceRepository_FindByID_NotFound(t *testing.T) {
	db := workspaceTestDB(t)
	_, err := gormstorage.NewWorkspaceRepository(db).FindByID(context.Background(), "nonexistent")
	assert.ErrorIs(t, err, workspace.ErrWorkspaceNotFound)
}

func TestGORMWorkspaceRepository_Update(t *testing.T) {
	db := workspaceTestDB(t)
	repo := gormstorage.NewWorkspaceRepository(db)
	ctx := context.Background()

	ws, _ := repo.Create(ctx, workspace.CreateWorkspaceParams{Name: "Old", CreatedByUserID: "user1"})
	name := "New"
	updated, err := repo.Update(ctx, ws.ID, workspace.UpdateWorkspaceParams{Name: &name})
	require.NoError(t, err)
	assert.Equal(t, "New", updated.Name)
}

func TestGORMWorkspaceRepository_Update_NotFound(t *testing.T) {
	db := workspaceTestDB(t)
	name := "X"
	_, err := gormstorage.NewWorkspaceRepository(db).Update(context.Background(), "nope", workspace.UpdateWorkspaceParams{Name: &name})
	assert.ErrorIs(t, err, workspace.ErrWorkspaceNotFound)
}

func TestGORMWorkspaceRepository_ListByUser(t *testing.T) {
	db := workspaceTestDB(t)
	repo := gormstorage.NewWorkspaceRepository(db)
	memberships := gormstorage.NewMembershipRepository(db)
	ctx := context.Background()

	ws1, _ := repo.Create(ctx, workspace.CreateWorkspaceParams{Name: "A", CreatedByUserID: "user1"})
	ws2, _ := repo.Create(ctx, workspace.CreateWorkspaceParams{Name: "B", CreatedByUserID: "user2"})
	ws3, _ := repo.Create(ctx, workspace.CreateWorkspaceParams{Name: "C", CreatedByUserID: "user1"})

	memberships.Create(ctx, workspace.CreateMembershipParams{WorkspaceID: ws1.ID, UserID: "user1", Role: "admin"})
	memberships.Create(ctx, workspace.CreateMembershipParams{WorkspaceID: ws2.ID, UserID: "user1", Role: "admin"})
	memberships.Create(ctx, workspace.CreateMembershipParams{WorkspaceID: ws3.ID, UserID: "user1", Role: "admin"})

	list, err := repo.ListByUser(ctx, "user1")
	require.NoError(t, err)
	assert.Len(t, list, 3)
}

func TestGORMWorkspaceRepository_ListByUser_Empty(t *testing.T) {
	db := workspaceTestDB(t)
	list, err := gormstorage.NewWorkspaceRepository(db).ListByUser(context.Background(), "nobody")
	require.NoError(t, err)
	assert.Empty(t, list)
}

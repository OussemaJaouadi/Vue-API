package gormstorage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"vue-api/backend/internal/collection"
	"vue-api/backend/internal/environment"
	gormstorage "vue-api/backend/internal/storage/gorm"
	"vue-api/backend/internal/workspace"
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

func TestGORMWorkspaceRepository_DeleteCleansWorkspaceData(t *testing.T) {
	db := workspaceTestDB(t)
	ctx := context.Background()
	repo := gormstorage.NewWorkspaceRepository(db)
	memberships := gormstorage.NewMembershipRepository(db)
	grants := gormstorage.NewGrantRepository(db)
	collections := gormstorage.NewCollectionRepository(db)
	environments := gormstorage.NewEnvironmentRepository(db)

	ws, err := repo.Create(ctx, workspace.CreateWorkspaceParams{Name: "A", CreatedByUserID: "user1"})
	require.NoError(t, err)
	otherWS, err := repo.Create(ctx, workspace.CreateWorkspaceParams{Name: "B", CreatedByUserID: "user1"})
	require.NoError(t, err)

	_, err = memberships.Create(ctx, workspace.CreateMembershipParams{WorkspaceID: ws.ID, UserID: "user1", Role: "admin", CreatedByUserID: "user1"})
	require.NoError(t, err)
	_, err = memberships.Create(ctx, workspace.CreateMembershipParams{WorkspaceID: otherWS.ID, UserID: "user1", Role: "admin", CreatedByUserID: "user1"})
	require.NoError(t, err)

	folder, err := collections.CreateFolder(ctx, collection.CreateFolderParams{WorkspaceID: ws.ID, Name: "Auth"})
	require.NoError(t, err)
	_, err = collections.CreateRequest(ctx, collection.CreateRequestParams{WorkspaceID: ws.ID, CollectionID: &folder.ID, Method: "GET", Name: "Me", Path: "/me"})
	require.NoError(t, err)
	_, err = collections.CreateRequest(ctx, collection.CreateRequestParams{WorkspaceID: otherWS.ID, Method: "GET", Name: "Health", Path: "/healthz"})
	require.NoError(t, err)

	env, err := environments.CreateEnvironment(ctx, environment.CreateEnvironmentParams{WorkspaceID: ws.ID, Name: "Local", CreatedByUserID: "user1"})
	require.NoError(t, err)
	_, err = environments.CreateVariable(ctx, environment.CreateVariableParams{EnvironmentID: env.ID, Key: "TOKEN", Value: "secret", Secret: true})
	require.NoError(t, err)

	_, err = grants.Set(ctx, "user1", ws.ID, []workspace.GrantInput{{ResourceType: "collection", ResourceID: folder.ID, AccessLevel: "admin"}})
	require.NoError(t, err)

	require.NoError(t, repo.Delete(ctx, ws.ID))

	_, err = repo.FindByID(ctx, ws.ID)
	assert.ErrorIs(t, err, workspace.ErrWorkspaceNotFound)

	assertCount(t, db, &gormstorage.WorkspaceMembershipModel{}, "workspace_id = ?", ws.ID, 0)
	assertCount(t, db, &gormstorage.ResourceGrantModel{}, "workspace_id = ?", ws.ID, 0)
	assertCount(t, db, &gormstorage.FolderModel{}, "workspace_id = ?", ws.ID, 0)
	assertCount(t, db, &gormstorage.RequestModel{}, "workspace_id = ?", ws.ID, 0)
	assertCount(t, db, &gormstorage.EnvironmentModel{}, "workspace_id = ?", ws.ID, 0)
	assertCount(t, db, &gormstorage.VariableModel{}, "environment_id = ?", env.ID, 0)

	assertCount(t, db, &gormstorage.WorkspaceModel{}, "id = ?", otherWS.ID, 1)
	assertCount(t, db, &gormstorage.RequestModel{}, "workspace_id = ?", otherWS.ID, 1)
}

func TestGORMWorkspaceRepository_Delete_NotFound(t *testing.T) {
	db := workspaceTestDB(t)
	err := gormstorage.NewWorkspaceRepository(db).Delete(context.Background(), "nope")
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

func assertCount(t *testing.T, db *gorm.DB, model any, query string, value any, expected int64) {
	t.Helper()
	var count int64
	require.NoError(t, db.Model(model).Where(query, value).Count(&count).Error)
	assert.Equal(t, expected, count)
}

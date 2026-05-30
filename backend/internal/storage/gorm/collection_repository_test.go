package gormstorage_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"vue-api/backend/internal/collection"
	gormstorage "vue-api/backend/internal/storage/gorm"
)

func newCollectionRepo(t *testing.T) (*gorm.DB, *gormstorage.CollectionRepository) {
	t.Helper()
	db := openTestDB(t)
	require.NoError(t, gormstorage.Migrate(db))
	repo := gormstorage.NewCollectionRepository(db)
	return db, repo
}

func TestCreateFolder_ReturnsFolderWithFields(t *testing.T) {
	_, repo := newCollectionRepo(t)

	folder, err := repo.CreateFolder(context.Background(), collection.CreateFolderParams{
		WorkspaceID: "ws1",
		Name:        "My API",
		Icon:        "PhGlobe",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, folder.ID)
	assert.Equal(t, "ws1", folder.WorkspaceID)
	assert.Equal(t, "My API", folder.Name)
	assert.Equal(t, "PhGlobe", folder.Icon)
	assert.False(t, folder.CreatedAt.IsZero())
	assert.False(t, folder.UpdatedAt.IsZero())
}

func TestCreateFolder_EmptyIconDefaultsToPhGlobe(t *testing.T) {
	_, repo := newCollectionRepo(t)

	folder, err := repo.CreateFolder(context.Background(), collection.CreateFolderParams{
		WorkspaceID: "ws1",
		Name:        "Defaults",
	})

	require.NoError(t, err)
	assert.Equal(t, "PhGlobe", folder.Icon)
}

func TestCreateFolder_DuplicateNameReturnsError(t *testing.T) {
	_, repo := newCollectionRepo(t)

	_, err := repo.CreateFolder(context.Background(), collection.CreateFolderParams{
		WorkspaceID: "ws1",
		Name:        "My API",
	})
	require.NoError(t, err)

	_, err = repo.CreateFolder(context.Background(), collection.CreateFolderParams{
		WorkspaceID: "ws1",
		Name:        "My API",
	})
	require.ErrorIs(t, err, collection.ErrFolderNameTaken)
}

func TestCreateFolder_SameNameDifferentWorkspaceIsOK(t *testing.T) {
	_, repo := newCollectionRepo(t)

	_, err := repo.CreateFolder(context.Background(), collection.CreateFolderParams{
		WorkspaceID: "ws1",
		Name:        "My API",
	})
	require.NoError(t, err)

	_, err = repo.CreateFolder(context.Background(), collection.CreateFolderParams{
		WorkspaceID: "ws2",
		Name:        "My API",
	})
	require.NoError(t, err)
}

func TestListFolders_ReturnsAllFoldersForWorkspace(t *testing.T) {
	_, repo := newCollectionRepo(t)

	_, _ = repo.CreateFolder(context.Background(), collection.CreateFolderParams{
		WorkspaceID: "ws1", Name: "Auth", Icon: "PhKey",
	})
	_, _ = repo.CreateFolder(context.Background(), collection.CreateFolderParams{
		WorkspaceID: "ws1", Name: "Users", Icon: "PhUser",
	})

	folders, err := repo.ListFolders(context.Background(), "ws1")

	require.NoError(t, err)
	assert.Len(t, folders, 2)
	assert.Equal(t, "Auth", folders[0].Name)
	assert.Equal(t, "Users", folders[1].Name)
}

func TestListFolders_FiltersByWorkspace(t *testing.T) {
	_, repo := newCollectionRepo(t)

	_, _ = repo.CreateFolder(context.Background(), collection.CreateFolderParams{
		WorkspaceID: "ws1", Name: "Auth",
	})
	_, _ = repo.CreateFolder(context.Background(), collection.CreateFolderParams{
		WorkspaceID: "ws2", Name: "Other",
	})

	folders, err := repo.ListFolders(context.Background(), "ws1")
	require.NoError(t, err)
	assert.Len(t, folders, 1)
	assert.Equal(t, "Auth", folders[0].Name)
}

func TestListFolders_EmptyWorkspaceReturnsEmpty(t *testing.T) {
	_, repo := newCollectionRepo(t)

	folders, err := repo.ListFolders(context.Background(), "empty")
	require.NoError(t, err)
	assert.Empty(t, folders)
}

func TestUpdateFolder_Name(t *testing.T) {
	_, repo := newCollectionRepo(t)

	folder, _ := repo.CreateFolder(context.Background(), collection.CreateFolderParams{
		WorkspaceID: "ws1", Name: "Old Name",
	})

	newName := "New Name"
	updated, err := repo.UpdateFolder(context.Background(), folder.ID, collection.UpdateFolderParams{
		Name: &newName,
	})

	require.NoError(t, err)
	assert.Equal(t, "New Name", updated.Name)
	assert.Equal(t, folder.Icon, updated.Icon)
}

func TestUpdateFolder_Icon(t *testing.T) {
	_, repo := newCollectionRepo(t)

	folder, _ := repo.CreateFolder(context.Background(), collection.CreateFolderParams{
		WorkspaceID: "ws1", Name: "Test",
	})

	newIcon := "PhLock"
	updated, err := repo.UpdateFolder(context.Background(), folder.ID, collection.UpdateFolderParams{
		Icon: &newIcon,
	})

	require.NoError(t, err)
	assert.Equal(t, "PhLock", updated.Icon)
	assert.Equal(t, "Test", updated.Name)
}

func TestUpdateFolder_NotFound(t *testing.T) {
	_, repo := newCollectionRepo(t)

	name := "Nope"
	_, err := repo.UpdateFolder(context.Background(), "nonexistent", collection.UpdateFolderParams{
		Name: &name,
	})

	require.ErrorIs(t, err, collection.ErrFolderNotFound)
}

func TestDeleteFolder_RemovesFolder(t *testing.T) {
	_, repo := newCollectionRepo(t)

	folder, _ := repo.CreateFolder(context.Background(), collection.CreateFolderParams{
		WorkspaceID: "ws1", Name: "Temp",
	})

	err := repo.DeleteFolder(context.Background(), folder.ID)
	require.NoError(t, err)

	folders, _ := repo.ListFolders(context.Background(), "ws1")
	assert.Empty(t, folders)
}

func TestDeleteFolder_CascadesToRequests(t *testing.T) {
	_, repo := newCollectionRepo(t)

	folder, _ := repo.CreateFolder(context.Background(), collection.CreateFolderParams{
		WorkspaceID: "ws1", Name: "Temp",
	})
	req, _ := repo.CreateRequest(context.Background(), collection.CreateRequestParams{
		CollectionID: &folder.ID,
		WorkspaceID:  "ws1",
		Method:       "GET",
		Name:         "test",
		Path:         "/test",
	})

	err := repo.DeleteFolder(context.Background(), folder.ID)
	require.NoError(t, err)

	err = repo.DeleteRequest(context.Background(), req.ID)
	require.ErrorIs(t, err, collection.ErrRequestNotFound)
}

func TestDeleteFolder_NotFound(t *testing.T) {
	_, repo := newCollectionRepo(t)

	err := repo.DeleteFolder(context.Background(), "nonexistent")
	require.ErrorIs(t, err, collection.ErrFolderNotFound)
}

func TestCreateRequest_RootLevel(t *testing.T) {
	_, repo := newCollectionRepo(t)

	req, err := repo.CreateRequest(context.Background(), collection.CreateRequestParams{
		WorkspaceID: "ws1",
		Method:      "GET",
		Name:        "Health Check",
		Path:        "/health",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, req.ID)
	assert.Nil(t, req.CollectionID)
	assert.Equal(t, "ws1", req.WorkspaceID)
	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "Health Check", req.Name)
	assert.Equal(t, "/health", req.Path)
	assert.False(t, req.CreatedAt.IsZero())
}

func TestCreateRequest_InFolder(t *testing.T) {
	_, repo := newCollectionRepo(t)

	folder, _ := repo.CreateFolder(context.Background(), collection.CreateFolderParams{
		WorkspaceID: "ws1", Name: "Auth",
	})

	req, err := repo.CreateRequest(context.Background(), collection.CreateRequestParams{
		CollectionID: &folder.ID,
		WorkspaceID:  "ws1",
		Method:       "POST",
		Name:         "Login",
		Path:         "/auth/login",
	})

	require.NoError(t, err)
	require.NotNil(t, req.CollectionID)
	assert.Equal(t, folder.ID, *req.CollectionID)
}

func TestListRequests_RootLevel(t *testing.T) {
	_, repo := newCollectionRepo(t)

	_, _ = repo.CreateRequest(context.Background(), collection.CreateRequestParams{
		WorkspaceID: "ws1", Method: "GET", Name: "Root1", Path: "/r1",
	})
	_, _ = repo.CreateRequest(context.Background(), collection.CreateRequestParams{
		WorkspaceID: "ws1", Method: "POST", Name: "Root2", Path: "/r2",
	})

	requests, err := repo.ListRootRequests(context.Background(), "ws1")

	require.NoError(t, err)
	assert.Len(t, requests, 2)
}

func TestListRequests_InFolder(t *testing.T) {
	_, repo := newCollectionRepo(t)

	folder, _ := repo.CreateFolder(context.Background(), collection.CreateFolderParams{
		WorkspaceID: "ws1", Name: "Auth",
	})

	_, _ = repo.CreateRequest(context.Background(), collection.CreateRequestParams{
		CollectionID: &folder.ID, WorkspaceID: "ws1", Method: "POST", Name: "Login", Path: "/login",
	})
	_, _ = repo.CreateRequest(context.Background(), collection.CreateRequestParams{
		CollectionID: &folder.ID, WorkspaceID: "ws1", Method: "GET", Name: "Logout", Path: "/logout",
	})

	requests, err := repo.ListRequests(context.Background(), "ws1", &folder.ID)

	require.NoError(t, err)
	assert.Len(t, requests, 2)
}

func TestListRequests_FiltersByWorkspaceAndFolder(t *testing.T) {
	_, repo := newCollectionRepo(t)

	folder, _ := repo.CreateFolder(context.Background(), collection.CreateFolderParams{
		WorkspaceID: "ws1", Name: "Auth",
	})
	_, _ = repo.CreateRequest(context.Background(), collection.CreateRequestParams{
		WorkspaceID: "ws1", Method: "GET", Name: "Root", Path: "/",
	})
	_, _ = repo.CreateRequest(context.Background(), collection.CreateRequestParams{
		CollectionID: &folder.ID, WorkspaceID: "ws1", Method: "POST", Name: "InFolder", Path: "/in",
	})
	_, _ = repo.CreateRequest(context.Background(), collection.CreateRequestParams{
		CollectionID: &folder.ID, WorkspaceID: "ws2", Method: "GET", Name: "OtherWS", Path: "/other",
	})

	rootReqs, _ := repo.ListRootRequests(context.Background(), "ws1")
	assert.Len(t, rootReqs, 1)
	assert.Equal(t, "Root", rootReqs[0].Name)

	folderReqs, _ := repo.ListRequests(context.Background(), "ws1", &folder.ID)
	assert.Len(t, folderReqs, 1)
	assert.Equal(t, "InFolder", folderReqs[0].Name)
}

func TestUpdateRequest_MethodNamePath(t *testing.T) {
	_, repo := newCollectionRepo(t)

	req, _ := repo.CreateRequest(context.Background(), collection.CreateRequestParams{
		WorkspaceID: "ws1", Method: "GET", Name: "Old", Path: "/old",
	})

	newMethod := "POST"
	newName := "Updated"
	newPath := "/new"
	updated, err := repo.UpdateRequest(context.Background(), req.ID, collection.UpdateRequestParams{
		Method: &newMethod,
		Name:   &newName,
		Path:   &newPath,
	})

	require.NoError(t, err)
	assert.Equal(t, "POST", updated.Method)
	assert.Equal(t, "Updated", updated.Name)
	assert.Equal(t, "/new", updated.Path)
}

func TestUpdateRequest_PartialUpdate(t *testing.T) {
	_, repo := newCollectionRepo(t)

	req, _ := repo.CreateRequest(context.Background(), collection.CreateRequestParams{
		WorkspaceID: "ws1", Method: "GET", Name: "Test", Path: "/test",
	})

	newMethod := "PUT"
	updated, err := repo.UpdateRequest(context.Background(), req.ID, collection.UpdateRequestParams{
		Method: &newMethod,
	})

	require.NoError(t, err)
	assert.Equal(t, "PUT", updated.Method)
	assert.Equal(t, "Test", updated.Name)
	assert.Equal(t, "/test", updated.Path)
}

func TestUpdateRequest_NotFound(t *testing.T) {
	_, repo := newCollectionRepo(t)

	name := "Nope"
	_, err := repo.UpdateRequest(context.Background(), "nonexistent", collection.UpdateRequestParams{
		Name: &name,
	})

	require.ErrorIs(t, err, collection.ErrRequestNotFound)
}

func TestDeleteRequest_RemovesRequest(t *testing.T) {
	_, repo := newCollectionRepo(t)

	req, _ := repo.CreateRequest(context.Background(), collection.CreateRequestParams{
		WorkspaceID: "ws1", Method: "GET", Name: "Temp", Path: "/temp",
	})

	err := repo.DeleteRequest(context.Background(), req.ID)
	require.NoError(t, err)

	requests, _ := repo.ListRootRequests(context.Background(), "ws1")
	assert.Empty(t, requests)
}

func TestDeleteRequest_NotFound(t *testing.T) {
	_, repo := newCollectionRepo(t)

	err := repo.DeleteRequest(context.Background(), "nonexistent")
	require.ErrorIs(t, err, collection.ErrRequestNotFound)
}

func TestCreateRequest_DefaultSortOrder(t *testing.T) {
	_, repo := newCollectionRepo(t)

	req, err := repo.CreateRequest(context.Background(), collection.CreateRequestParams{
		WorkspaceID: "ws1", Method: "GET", Name: "Test", Path: "/test",
	})

	require.NoError(t, err)
	assert.Equal(t, 0, req.SortOrder)
}

func TestReorderRequests_MovesRequestBetweenRootAndFolder(t *testing.T) {
	_, repo := newCollectionRepo(t)

	folder, _ := repo.CreateFolder(context.Background(), collection.CreateFolderParams{
		WorkspaceID: "ws1", Name: "Auth",
	})
	first, _ := repo.CreateRequest(context.Background(), collection.CreateRequestParams{
		WorkspaceID: "ws1", Method: "GET", Name: "First", Path: "/first",
	})
	second, _ := repo.CreateRequest(context.Background(), collection.CreateRequestParams{
		WorkspaceID: "ws1", Method: "GET", Name: "Second", Path: "/second",
	})

	err := repo.ReorderRequests(context.Background(), "ws1", []collection.RequestOrderGroup{
		{CollectionID: nil, RequestIDs: []string{second.ID}},
		{CollectionID: &folder.ID, RequestIDs: []string{first.ID}},
	})

	require.NoError(t, err)

	rootRequests, err := repo.ListRootRequests(context.Background(), "ws1")
	require.NoError(t, err)
	require.Len(t, rootRequests, 1)
	assert.Equal(t, second.ID, rootRequests[0].ID)
	assert.Equal(t, 0, rootRequests[0].SortOrder)

	folderRequests, err := repo.ListRequests(context.Background(), "ws1", &folder.ID)
	require.NoError(t, err)
	require.Len(t, folderRequests, 1)
	assert.Equal(t, first.ID, folderRequests[0].ID)
	assert.Equal(t, 0, folderRequests[0].SortOrder)
}

func TestReorderRequests_RejectsCollectionFromAnotherWorkspace(t *testing.T) {
	_, repo := newCollectionRepo(t)

	folder, _ := repo.CreateFolder(context.Background(), collection.CreateFolderParams{
		WorkspaceID: "ws2", Name: "Other",
	})
	req, _ := repo.CreateRequest(context.Background(), collection.CreateRequestParams{
		WorkspaceID: "ws1", Method: "GET", Name: "First", Path: "/first",
	})

	err := repo.ReorderRequests(context.Background(), "ws1", []collection.RequestOrderGroup{
		{CollectionID: &folder.ID, RequestIDs: []string{req.ID}},
	})

	require.ErrorIs(t, err, collection.ErrFolderNotFound)
}

func TestCreateFolder_DefaultSortOrder(t *testing.T) {
	_, repo := newCollectionRepo(t)

	folder, err := repo.CreateFolder(context.Background(), collection.CreateFolderParams{
		WorkspaceID: "ws1", Name: "Test",
	})

	require.NoError(t, err)
	assert.Equal(t, 0, folder.SortOrder)
}

func TestListRequests_EmptyFolder(t *testing.T) {
	_, repo := newCollectionRepo(t)

	folder, _ := repo.CreateFolder(context.Background(), collection.CreateFolderParams{
		WorkspaceID: "ws1", Name: "Empty",
	})

	requests, err := repo.ListRequests(context.Background(), "ws1", &folder.ID)
	require.NoError(t, err)
	assert.Empty(t, requests)
}

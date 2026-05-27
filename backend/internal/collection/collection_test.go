package collection_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"vue-api/backend/internal/collection"
)

func TestErrorSentinels_AreDistinct(t *testing.T) {
	assert.ErrorIs(t, collection.ErrFolderNotFound, collection.ErrFolderNotFound)
	assert.ErrorIs(t, collection.ErrRequestNotFound, collection.ErrRequestNotFound)
	assert.ErrorIs(t, collection.ErrFolderNameTaken, collection.ErrFolderNameTaken)

	assert.False(t, errors.Is(collection.ErrFolderNotFound, collection.ErrRequestNotFound))
	assert.False(t, errors.Is(collection.ErrFolderNotFound, collection.ErrFolderNameTaken))
	assert.False(t, errors.Is(collection.ErrRequestNotFound, collection.ErrFolderNameTaken))
}

func TestCreateFolderParams_HoldsValues(t *testing.T) {
	params := collection.CreateFolderParams{
		WorkspaceID: "ws1",
		Name:        "My Folder",
		Icon:        "PhKey",
	}

	assert.Equal(t, "ws1", params.WorkspaceID)
	assert.Equal(t, "My Folder", params.Name)
	assert.Equal(t, "PhKey", params.Icon)
}

func TestUpdateFolderParams_PartialUpdates(t *testing.T) {
	name := "Renamed"
	icon := "PhLock"
	params := collection.UpdateFolderParams{
		Name: &name,
		Icon: &icon,
	}

	assert.NotNil(t, params.Name)
	assert.NotNil(t, params.Icon)
	assert.Equal(t, "Renamed", *params.Name)
	assert.Equal(t, "PhLock", *params.Icon)
}

func TestUpdateFolderParams_NilFields(t *testing.T) {
	params := collection.UpdateFolderParams{}

	assert.Nil(t, params.Name)
	assert.Nil(t, params.Icon)
}

func TestCreateRequestParams_HoldsValues(t *testing.T) {
	collID := "folder1"
	params := collection.CreateRequestParams{
		CollectionID: &collID,
		WorkspaceID:  "ws1",
		Method:       "GET",
		Name:         "Health",
		Path:         "/health",
	}

	assert.NotNil(t, params.CollectionID)
	assert.Equal(t, "folder1", *params.CollectionID)
	assert.Equal(t, "ws1", params.WorkspaceID)
	assert.Equal(t, "GET", params.Method)
	assert.Equal(t, "Health", params.Name)
	assert.Equal(t, "/health", params.Path)
}

func TestCreateRequestParams_NilCollectionID(t *testing.T) {
	params := collection.CreateRequestParams{
		WorkspaceID: "ws1",
		Method:      "POST",
		Name:        "Login",
		Path:        "/login",
	}

	assert.Nil(t, params.CollectionID)
}

func TestUpdateRequestParams_PartialUpdates(t *testing.T) {
	method := "PUT"
	name := "Updated"
	path := "/new"
	params := collection.UpdateRequestParams{
		Method: &method,
		Name:   &name,
		Path:   &path,
	}

	assert.Equal(t, "PUT", *params.Method)
	assert.Equal(t, "Updated", *params.Name)
	assert.Equal(t, "/new", *params.Path)
}

func TestUpdateRequestParams_NilFields(t *testing.T) {
	params := collection.UpdateRequestParams{}

	assert.Nil(t, params.Method)
	assert.Nil(t, params.Name)
	assert.Nil(t, params.Path)
}

func TestFolderWithRequests_HoldsValues(t *testing.T) {
	fwr := collection.FolderWithRequests{
		Folder: collection.Folder{
			ID:   "f1",
			Name: "Auth",
		},
		Requests: []collection.Request{
			{ID: "r1", Name: "Login"},
			{ID: "r2", Name: "Logout"},
		},
	}

	assert.Equal(t, "f1", fwr.Folder.ID)
	assert.Len(t, fwr.Requests, 2)
	assert.Equal(t, "Login", fwr.Requests[0].Name)
	assert.Equal(t, "Logout", fwr.Requests[1].Name)
}

func TestFolder_DefaultZeroValues(t *testing.T) {
	var f collection.Folder

	assert.Empty(t, f.ID)
	assert.Empty(t, f.WorkspaceID)
	assert.Empty(t, f.Name)
	assert.Empty(t, f.Icon)
	assert.Equal(t, 0, f.SortOrder)
	assert.True(t, f.CreatedAt.IsZero())
	assert.True(t, f.UpdatedAt.IsZero())
}

func TestRequest_DefaultZeroValues(t *testing.T) {
	var r collection.Request

	assert.Empty(t, r.ID)
	assert.Nil(t, r.CollectionID)
	assert.Empty(t, r.WorkspaceID)
	assert.Empty(t, r.Method)
	assert.Empty(t, r.Name)
	assert.Empty(t, r.Path)
	assert.Equal(t, 0, r.SortOrder)
	assert.True(t, r.CreatedAt.IsZero())
	assert.True(t, r.UpdatedAt.IsZero())
}

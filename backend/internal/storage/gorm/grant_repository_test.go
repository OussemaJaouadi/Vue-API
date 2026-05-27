package gormstorage_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vue-api/backend/internal/workspace"
	gormstorage "vue-api/backend/internal/storage/gorm"
)

func TestGORMGrantRepository_SetAndList(t *testing.T) {
	db := workspaceTestDB(t)
	repo := gormstorage.NewGrantRepository(db)
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

func TestGORMGrantRepository_Set_ReplacesOld(t *testing.T) {
	db := workspaceTestDB(t)
	repo := gormstorage.NewGrantRepository(db)
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

func TestGORMGrantRepository_Set_UserIsolation(t *testing.T) {
	db := workspaceTestDB(t)
	repo := gormstorage.NewGrantRepository(db)
	ctx := context.Background()

	repo.Set(ctx, "user1", "ws1", []workspace.GrantInput{{ResourceType: "collection", ResourceID: "c1", AccessLevel: "read"}})
	repo.Set(ctx, "user2", "ws1", []workspace.GrantInput{{ResourceType: "collection", ResourceID: "c1", AccessLevel: "admin"}})

	list1, _ := repo.ListByUserAndWorkspace(ctx, "user1", "ws1")
	list2, _ := repo.ListByUserAndWorkspace(ctx, "user2", "ws1")
	assert.Len(t, list1, 1)
	assert.Len(t, list2, 1)
	assert.Equal(t, "read", list1[0].AccessLevel)
	assert.Equal(t, "admin", list2[0].AccessLevel)
}

func TestGORMGrantRepository_List_Empty(t *testing.T) {
	db := workspaceTestDB(t)
	grants, err := gormstorage.NewGrantRepository(db).ListByUserAndWorkspace(context.Background(), "user1", "ws1")
	require.NoError(t, err)
	assert.Empty(t, grants)
}

func TestGORMGrantRepository_Set_Empty(t *testing.T) {
	db := workspaceTestDB(t)
	repo := gormstorage.NewGrantRepository(db)
	grants, err := repo.Set(context.Background(), "user1", "ws1", nil)
	require.NoError(t, err)
	assert.Empty(t, grants)
}

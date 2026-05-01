package gormstorage_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"vue-api/backend/internal/auth"
	"vue-api/backend/internal/auth/authtest"
	gormstorage "vue-api/backend/internal/storage/gorm"
)

func TestUserRepositoryContract(t *testing.T) {
	authtest.RunUserRepositoryContract(t, func(t *testing.T) auth.UserRepository {
		return newTestUserRepository(t)
	})
}

func newTestUserRepository(t *testing.T) *gormstorage.UserRepository {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, gormstorage.Migrate(db))

	repo := gormstorage.NewUserRepository(db)
	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})

	return repo
}

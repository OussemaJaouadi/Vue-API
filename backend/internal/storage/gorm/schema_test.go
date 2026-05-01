package gormstorage_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	gormstorage "vue-api/backend/internal/storage/gorm"
)

func TestPlanReportsSafeChangesForEmptyDatabase(t *testing.T) {
	db := openTestDB(t)

	plan, err := gormstorage.Plan(db)
	require.NoError(t, err)

	require.True(t, plan.HasChanges())
	require.False(t, plan.HasManualChanges())
	require.Contains(t, plan.String(), "create table users")
	require.Contains(t, plan.String(), "create table workspaces")
	require.Contains(t, plan.String(), "create table workspace_memberships")
	require.Contains(t, plan.String(), "create unique index idx_users_email")
}

func TestMigrateAppliesSafePlan(t *testing.T) {
	db := openTestDB(t)

	require.NoError(t, gormstorage.Migrate(db))

	plan, err := gormstorage.Plan(db)
	require.NoError(t, err)
	require.False(t, plan.HasChanges())
	require.NoError(t, gormstorage.VerifySchema(db))
}

func TestGeneratePlanAndApplyMigrationFiles(t *testing.T) {
	db := openTestDB(t)
	dir := t.TempDir()

	plan, err := gormstorage.Plan(db)
	require.NoError(t, err)

	path, err := gormstorage.GenerateMigrationFile(dir, "initial_schema", plan, time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	require.Equal(t, filepath.Join(dir, "20260501120000_initial_schema.sql"), path)

	content, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Contains(t, string(content), "CREATE TABLE users")

	pending, err := gormstorage.PendingMigrationPlan(db, dir)
	require.NoError(t, err)
	require.Contains(t, pending.String(), "20260501120000_initial_schema.sql")

	require.NoError(t, gormstorage.ApplyMigrationFiles(db, dir))
	require.NoError(t, gormstorage.VerifySchema(db))

	pending, err = gormstorage.PendingMigrationPlan(db, dir)
	require.NoError(t, err)
	require.False(t, pending.HasChanges())
}

func TestGeneratePlanUsesExistingMigrationFilesAsBaseline(t *testing.T) {
	db := openTestDB(t)
	dir := t.TempDir()

	plan, err := gormstorage.GenerateMigrationPlan(db, dir)
	require.NoError(t, err)

	_, err = gormstorage.GenerateMigrationFile(dir, "initial_schema", plan, time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC))
	require.NoError(t, err)

	plan, err = gormstorage.GenerateMigrationPlan(db, dir)
	require.NoError(t, err)
	require.False(t, plan.HasChanges())
}

func TestPlanReportsManualColumnChanges(t *testing.T) {
	db := openTestDB(t)
	require.NoError(t, db.Exec(`
		CREATE TABLE users (
			id text NOT NULL PRIMARY KEY,
			email integer,
			username text NOT NULL,
			password_hash text NOT NULL,
			global_role text NOT NULL,
			token_version integer NOT NULL DEFAULT 1,
			active integer NOT NULL DEFAULT 1,
			created_at datetime,
			updated_at text NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`).Error)

	plan, err := gormstorage.Plan(db)
	require.NoError(t, err)

	require.True(t, plan.HasManualChanges())
	require.Contains(t, plan.String(), "users.email: expected text NOT NULL")
	require.Contains(t, plan.String(), "users.created_at: expected datetime NOT NULL")
}

func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	return db
}

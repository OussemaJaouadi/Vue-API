package auth_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"vue-api/backend/internal/auth"
)

func TestScopedRolePermissions(t *testing.T) {
	tests := []struct {
		name       string
		globalRole string
		scopedRole string
		permission string
		allowed    bool
	}{
		{
			name:       "manager can manage members without scoped role",
			globalRole: auth.GlobalRoleManager,
			scopedRole: "",
			permission: auth.PermissionManageMembers,
			allowed:    true,
		},
		{
			name:       "admin can delete projects",
			globalRole: auth.GlobalRoleUser,
			scopedRole: auth.ScopedRoleAdmin,
			permission: auth.PermissionDeleteProject,
			allowed:    true,
		},
		{
			name:       "developer can manage collections",
			globalRole: auth.GlobalRoleUser,
			scopedRole: auth.ScopedRoleDeveloper,
			permission: auth.PermissionManageCollections,
			allowed:    true,
		},
		{
			name:       "developer cannot invite members",
			globalRole: auth.GlobalRoleUser,
			scopedRole: auth.ScopedRoleDeveloper,
			permission: auth.PermissionManageMembers,
			allowed:    false,
		},
		{
			name:       "tester can send requests",
			globalRole: auth.GlobalRoleUser,
			scopedRole: auth.ScopedRoleTester,
			permission: auth.PermissionSendRequests,
			allowed:    true,
		},
		{
			name:       "tester cannot manage collections",
			globalRole: auth.GlobalRoleUser,
			scopedRole: auth.ScopedRoleTester,
			permission: auth.PermissionManageCollections,
			allowed:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.allowed, auth.Can(tt.globalRole, tt.scopedRole, tt.permission))
		})
	}
}

func TestValidScopedRole(t *testing.T) {
	require.True(t, auth.ValidScopedRole(auth.ScopedRoleAdmin))
	require.True(t, auth.ValidScopedRole(auth.ScopedRoleDeveloper))
	require.True(t, auth.ValidScopedRole(auth.ScopedRoleTester))
	require.False(t, auth.ValidScopedRole("owner"))
}

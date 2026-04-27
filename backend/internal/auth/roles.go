package auth

const (
	ScopedRoleAdmin     = "admin"
	ScopedRoleDeveloper = "developer"
	ScopedRoleTester    = "tester"
)

const (
	PermissionManageMembers     = "manage_members"
	PermissionCreateProject     = "create_project"
	PermissionDeleteProject     = "delete_project"
	PermissionManageProject     = "manage_project"
	PermissionManageCollections = "manage_collections"
	PermissionManageEnvironment = "manage_environment"
	PermissionManageAuthProfile = "manage_auth_profile"
	PermissionSendRequests      = "send_requests"
	PermissionViewCollections   = "view_collections"
)

func ValidScopedRole(role string) bool {
	switch role {
	case ScopedRoleAdmin, ScopedRoleDeveloper, ScopedRoleTester:
		return true
	default:
		return false
	}
}

func Can(globalRole string, scopedRole string, permission string) bool {
	if globalRole == GlobalRoleManager {
		return true
	}

	permissions, ok := scopedRolePermissions[scopedRole]
	if !ok {
		return false
	}

	return permissions[permission]
}

var scopedRolePermissions = map[string]map[string]bool{
	ScopedRoleAdmin: {
		PermissionManageMembers:     true,
		PermissionCreateProject:     true,
		PermissionDeleteProject:     true,
		PermissionManageProject:     true,
		PermissionManageCollections: true,
		PermissionManageEnvironment: true,
		PermissionManageAuthProfile: true,
		PermissionSendRequests:      true,
		PermissionViewCollections:   true,
	},
	ScopedRoleDeveloper: {
		PermissionManageCollections: true,
		PermissionManageEnvironment: true,
		PermissionManageAuthProfile: true,
		PermissionSendRequests:      true,
		PermissionViewCollections:   true,
	},
	ScopedRoleTester: {
		PermissionSendRequests:    true,
		PermissionViewCollections: true,
	},
}

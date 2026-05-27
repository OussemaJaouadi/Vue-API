package workspace_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"vue-api/backend/internal/workspace"
)

func TestCreateWorkspaceParams_ZeroValues(t *testing.T) {
	p := workspace.CreateWorkspaceParams{}
	assert.Empty(t, p.Name)
	assert.Empty(t, p.CreatedByUserID)
}

func TestCreateMembershipParams_ZeroValues(t *testing.T) {
	p := workspace.CreateMembershipParams{}
	assert.Empty(t, p.WorkspaceID)
	assert.Empty(t, p.UserID)
	assert.Empty(t, p.Role)
}

func TestUpdateWorkspaceParams_ZeroValues(t *testing.T) {
	p := workspace.UpdateWorkspaceParams{}
	assert.Nil(t, p.Name)
}

func TestUpdateMembershipParams_ZeroValues(t *testing.T) {
	p := workspace.UpdateMembershipParams{}
	assert.Nil(t, p.Role)
}

func TestGrantInput_ZeroValues(t *testing.T) {
	g := workspace.GrantInput{}
	assert.Empty(t, g.ResourceType)
	assert.Empty(t, g.ResourceID)
	assert.Empty(t, g.AccessLevel)
}

func TestErrorSentinels(t *testing.T) {
	assert.NotNil(t, workspace.ErrWorkspaceNotFound)
	assert.NotNil(t, workspace.ErrWorkspaceNameTaken)
	assert.NotNil(t, workspace.ErrMembershipNotFound)
	assert.NotNil(t, workspace.ErrAlreadyMember)
	assert.NotNil(t, workspace.ErrGrantNotFound)
}

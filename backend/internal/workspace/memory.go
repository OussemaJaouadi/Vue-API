package workspace

import (
	"context"
	"sort"
	"sync"
	"time"

	"vue-api/backend/internal/id"
)

type MemoryWorkspaceRepository struct {
	mu         sync.RWMutex
	workspaces map[string]Workspace
}

func NewMemoryWorkspaceRepository() *MemoryWorkspaceRepository {
	return &MemoryWorkspaceRepository{
		workspaces: make(map[string]Workspace),
	}
}

func (r *MemoryWorkspaceRepository) ListByUser(ctx context.Context, userID string) ([]Workspace, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []Workspace
	for _, ws := range r.workspaces {
		if ws.CreatedByUserID == userID {
			result = append(result, ws)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})
	return result, nil
}

func (r *MemoryWorkspaceRepository) Create(ctx context.Context, params CreateWorkspaceParams) (Workspace, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, ws := range r.workspaces {
		if ws.Name == params.Name && ws.CreatedByUserID == params.CreatedByUserID {
			return Workspace{}, ErrWorkspaceNameTaken
		}
	}
	now := time.Now().UTC()
	ws := Workspace{
		ID:              id.NewUUIDV7(),
		Name:            params.Name,
		CreatedByUserID: params.CreatedByUserID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	r.workspaces[ws.ID] = ws
	return ws, nil
}

func (r *MemoryWorkspaceRepository) FindByID(ctx context.Context, id string) (Workspace, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ws, exists := r.workspaces[id]
	if !exists {
		return Workspace{}, ErrWorkspaceNotFound
	}
	return ws, nil
}

func (r *MemoryWorkspaceRepository) Update(ctx context.Context, id string, params UpdateWorkspaceParams) (Workspace, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ws, exists := r.workspaces[id]
	if !exists {
		return Workspace{}, ErrWorkspaceNotFound
	}
	if params.Name != nil {
		for _, other := range r.workspaces {
			if other.ID != id && other.Name == *params.Name && other.CreatedByUserID == ws.CreatedByUserID {
				return Workspace{}, ErrWorkspaceNameTaken
			}
		}
		ws.Name = *params.Name
	}
	ws.UpdatedAt = time.Now().UTC()
	r.workspaces[id] = ws
	return ws, nil
}

func (r *MemoryWorkspaceRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.workspaces[id]; !exists {
		return ErrWorkspaceNotFound
	}
	delete(r.workspaces, id)
	return nil
}

type MemoryMembershipRepository struct {
	mu          sync.RWMutex
	memberships map[string]Membership
}

func NewMemoryMembershipRepository() *MemoryMembershipRepository {
	return &MemoryMembershipRepository{
		memberships: make(map[string]Membership),
	}
}

func (r *MemoryMembershipRepository) ListByWorkspace(ctx context.Context, workspaceID string) ([]Membership, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []Membership
	for _, m := range r.memberships {
		if m.WorkspaceID == workspaceID {
			result = append(result, m)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})
	return result, nil
}

func (r *MemoryMembershipRepository) FindByUserAndWorkspace(ctx context.Context, userID string, workspaceID string) (Membership, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, m := range r.memberships {
		if m.UserID == userID && m.WorkspaceID == workspaceID {
			return m, nil
		}
	}
	return Membership{}, ErrMembershipNotFound
}

func (r *MemoryMembershipRepository) Create(ctx context.Context, params CreateMembershipParams) (Membership, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, m := range r.memberships {
		if m.UserID == params.UserID && m.WorkspaceID == params.WorkspaceID {
			return Membership{}, ErrAlreadyMember
		}
	}
	now := time.Now().UTC()
	m := Membership{
		ID:              id.NewUUIDV7(),
		WorkspaceID:     params.WorkspaceID,
		UserID:          params.UserID,
		Role:            params.Role,
		CreatedByUserID: params.CreatedByUserID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	r.memberships[m.ID] = m
	return m, nil
}

func (r *MemoryMembershipRepository) UpdateRole(ctx context.Context, id string, params UpdateMembershipParams) (Membership, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	m, exists := r.memberships[id]
	if !exists {
		return Membership{}, ErrMembershipNotFound
	}
	if params.Role != nil {
		m.Role = *params.Role
	}
	m.UpdatedAt = time.Now().UTC()
	r.memberships[id] = m
	return m, nil
}

func (r *MemoryMembershipRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.memberships[id]; !exists {
		return ErrMembershipNotFound
	}
	delete(r.memberships, id)
	return nil
}

func (r *MemoryMembershipRepository) CountByWorkspace(ctx context.Context, workspaceID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var count int
	for _, m := range r.memberships {
		if m.WorkspaceID == workspaceID {
			count++
		}
	}
	return count, nil
}

type MemoryGrantRepository struct {
	mu     sync.RWMutex
	grants map[string]ResourceGrant
}

func NewMemoryGrantRepository() *MemoryGrantRepository {
	return &MemoryGrantRepository{
		grants: make(map[string]ResourceGrant),
	}
}

func (r *MemoryGrantRepository) ListByUserAndWorkspace(ctx context.Context, userID string, workspaceID string) ([]ResourceGrant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []ResourceGrant
	for _, g := range r.grants {
		if g.UserID == userID && g.WorkspaceID == workspaceID {
			result = append(result, g)
		}
	}
	return result, nil
}

func (r *MemoryGrantRepository) Set(ctx context.Context, userID string, workspaceID string, grants []GrantInput) ([]ResourceGrant, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, g := range r.grants {
		if g.UserID == userID && g.WorkspaceID == workspaceID {
			delete(r.grants, id)
		}
	}

	now := time.Now().UTC()
	var result []ResourceGrant
	for _, input := range grants {
		g := ResourceGrant{
			ID:           id.NewUUIDV7(),
			WorkspaceID:  workspaceID,
			UserID:       userID,
			ResourceType: input.ResourceType,
			ResourceID:   input.ResourceID,
			AccessLevel:  input.AccessLevel,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		r.grants[g.ID] = g
		result = append(result, g)
	}
	return result, nil
}

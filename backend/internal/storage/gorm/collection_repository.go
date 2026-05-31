package gormstorage

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"vue-api/backend/internal/collection"
	"vue-api/backend/internal/id"
)

type CollectionRepository struct {
	db *gorm.DB
}

func NewCollectionRepository(db *gorm.DB) *CollectionRepository {
	return &CollectionRepository{db: db}
}

func (repo *CollectionRepository) ListFolders(ctx context.Context, workspaceID string) ([]collection.Folder, error) {
	var models []FolderModel
	if err := repo.db.WithContext(ctx).
		Where("workspace_id = ?", workspaceID).
		Order("sort_order ASC, name ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	folders := make([]collection.Folder, 0, len(models))
	for _, m := range models {
		folders = append(folders, toDomainFolder(m))
	}

	return folders, nil
}

func (repo *CollectionRepository) ListRequests(ctx context.Context, workspaceID string, collectionID *string) ([]collection.Request, error) {
	var models []RequestModel
	query := repo.db.WithContext(ctx).Where("workspace_id = ?", workspaceID)
	if collectionID == nil {
		query = query.Where("collection_id IS NULL")
	} else {
		query = query.Where("collection_id = ?", *collectionID)
	}

	if err := query.Order("sort_order ASC, name ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	requests := make([]collection.Request, 0, len(models))
	for _, m := range models {
		requests = append(requests, toDomainRequest(m))
	}

	return requests, nil
}

func (repo *CollectionRepository) ListRootRequests(ctx context.Context, workspaceID string) ([]collection.Request, error) {
	return repo.ListRequests(ctx, workspaceID, nil)
}

func (repo *CollectionRepository) ListEnvironmentPolicies(ctx context.Context, workspaceID string, collectionID string) ([]collection.EnvironmentPolicy, error) {
	var rows []struct {
		CollectionID          string
		EnvironmentID         string
		EnvironmentName       string
		EnvironmentVisibility string
		DefaultEnvironment    bool
	}

	if err := repo.db.WithContext(ctx).
		Table("collection_environment_policies AS policy").
		Select("policy.collection_id, policy.environment_id, environments.name AS environment_name, environments.visibility AS environment_visibility, policy.default_environment").
		Joins("JOIN environments ON environments.id = policy.environment_id").
		Where("policy.workspace_id = ? AND policy.collection_id = ?", workspaceID, collectionID).
		Order("environments.sort_order ASC, environments.name ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	policies := make([]collection.EnvironmentPolicy, 0, len(rows))
	for _, row := range rows {
		policies = append(policies, collection.EnvironmentPolicy{
			CollectionID:          row.CollectionID,
			EnvironmentID:         row.EnvironmentID,
			EnvironmentName:       row.EnvironmentName,
			EnvironmentVisibility: row.EnvironmentVisibility,
			Default:               row.DefaultEnvironment,
		})
	}

	return policies, nil
}

func (repo *CollectionRepository) CreateFolder(ctx context.Context, params collection.CreateFolderParams) (collection.Folder, error) {
	var count int64
	if err := repo.db.WithContext(ctx).Model(&FolderModel{}).
		Where("workspace_id = ? AND name = ?", params.WorkspaceID, params.Name).
		Count(&count).Error; err != nil {
		return collection.Folder{}, err
	}
	if count > 0 {
		return collection.Folder{}, collection.ErrFolderNameTaken
	}

	icon := params.Icon
	if icon == "" {
		icon = "PhGlobe"
	}

	model := FolderModel{
		ID:          id.NewUUIDV7(),
		WorkspaceID: params.WorkspaceID,
		Name:        params.Name,
		Icon:        icon,
		SortOrder:   0,
	}

	if err := repo.db.WithContext(ctx).Create(&model).Error; err != nil {
		return collection.Folder{}, err
	}

	return toDomainFolder(model), nil
}

func (repo *CollectionRepository) SetEnvironmentPolicies(ctx context.Context, params collection.SetEnvironmentPolicyParams) ([]collection.EnvironmentPolicy, error) {
	err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var collectionCount int64
		if err := tx.Model(&FolderModel{}).
			Where("id = ? AND workspace_id = ?", params.CollectionID, params.WorkspaceID).
			Count(&collectionCount).Error; err != nil {
			return err
		}
		if collectionCount == 0 {
			return collection.ErrFolderNotFound
		}

		allowed := uniqueStrings(params.AllowedEnvironmentIDs)
		if params.DefaultEnvironmentID != "" && !containsString(allowed, params.DefaultEnvironmentID) {
			allowed = append([]string{params.DefaultEnvironmentID}, allowed...)
		}

		if len(allowed) > 0 {
			var envCount int64
			if err := tx.Model(&EnvironmentModel{}).
				Where("workspace_id = ? AND id IN ?", params.WorkspaceID, allowed).
				Count(&envCount).Error; err != nil {
				return err
			}
			if envCount != int64(len(allowed)) {
				return collection.ErrEnvironmentNotFound
			}
		}

		if err := tx.Where("workspace_id = ? AND collection_id = ?", params.WorkspaceID, params.CollectionID).
			Delete(&CollectionEnvironmentPolicyModel{}).Error; err != nil {
			return err
		}

		for _, environmentID := range allowed {
			model := CollectionEnvironmentPolicyModel{
				ID:                 id.NewUUIDV7(),
				WorkspaceID:        params.WorkspaceID,
				CollectionID:       params.CollectionID,
				EnvironmentID:      environmentID,
				DefaultEnvironment: environmentID == params.DefaultEnvironmentID,
			}
			if err := tx.Create(&model).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return repo.ListEnvironmentPolicies(ctx, params.WorkspaceID, params.CollectionID)
}

func (repo *CollectionRepository) CreateRequest(ctx context.Context, params collection.CreateRequestParams) (collection.Request, error) {
	model := RequestModel{
		ID:              id.NewUUIDV7(),
		CollectionID:    params.CollectionID,
		WorkspaceID:     params.WorkspaceID,
		Method:          params.Method,
		Name:            params.Name,
		Path:            params.Path,
		QueryParamsJSON: defaultJSON(params.QueryParamsJSON, "[]"),
		HeadersJSON:     defaultJSON(params.HeadersJSON, "[]"),
		Body:            params.Body,
		BodyLanguage:    defaultJSON(params.BodyLanguage, "json"),
		AuthConfigJSON:  defaultJSON(params.AuthConfigJSON, "{}"),
		SortOrder:       0,
	}

	if err := repo.db.WithContext(ctx).Create(&model).Error; err != nil {
		return collection.Request{}, err
	}

	return toDomainRequest(model), nil
}

func (repo *CollectionRepository) UpdateFolder(ctx context.Context, folderID string, params collection.UpdateFolderParams) (collection.Folder, error) {
	var model FolderModel
	if err := repo.db.WithContext(ctx).Where("id = ?", folderID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return collection.Folder{}, collection.ErrFolderNotFound
		}
		return collection.Folder{}, err
	}

	if params.Name != nil {
		model.Name = *params.Name
	}
	if params.Icon != nil {
		model.Icon = *params.Icon
	}

	if err := repo.db.WithContext(ctx).Save(&model).Error; err != nil {
		return collection.Folder{}, err
	}

	return toDomainFolder(model), nil
}

func (repo *CollectionRepository) UpdateRequest(ctx context.Context, requestID string, params collection.UpdateRequestParams) (collection.Request, error) {
	var model RequestModel
	if err := repo.db.WithContext(ctx).Where("id = ?", requestID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return collection.Request{}, collection.ErrRequestNotFound
		}
		return collection.Request{}, err
	}

	if params.Method != nil {
		model.Method = *params.Method
	}
	if params.Name != nil {
		model.Name = *params.Name
	}
	if params.Path != nil {
		model.Path = *params.Path
	}
	if params.QueryParamsJSON != nil {
		model.QueryParamsJSON = defaultJSON(*params.QueryParamsJSON, "[]")
	}
	if params.HeadersJSON != nil {
		model.HeadersJSON = defaultJSON(*params.HeadersJSON, "[]")
	}
	if params.Body != nil {
		model.Body = *params.Body
	}
	if params.BodyLanguage != nil {
		model.BodyLanguage = defaultJSON(*params.BodyLanguage, "json")
	}
	if params.AuthConfigJSON != nil {
		model.AuthConfigJSON = defaultJSON(*params.AuthConfigJSON, "{}")
	}

	if err := repo.db.WithContext(ctx).Save(&model).Error; err != nil {
		return collection.Request{}, err
	}

	return toDomainRequest(model), nil
}

func (repo *CollectionRepository) ReorderRequests(ctx context.Context, workspaceID string, groups []collection.RequestOrderGroup) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, group := range groups {
			if group.CollectionID != nil {
				var count int64
				if err := tx.Model(&FolderModel{}).
					Where("id = ? AND workspace_id = ?", *group.CollectionID, workspaceID).
					Count(&count).Error; err != nil {
					return err
				}
				if count == 0 {
					return collection.ErrFolderNotFound
				}
			}

			for index, requestID := range group.RequestIDs {
				updates := map[string]any{
					"collection_id": group.CollectionID,
					"sort_order":    index,
				}

				result := tx.Model(&RequestModel{}).
					Where("id = ? AND workspace_id = ?", requestID, workspaceID).
					Updates(updates)
				if result.Error != nil {
					return result.Error
				}
				if result.RowsAffected == 0 {
					return collection.ErrRequestNotFound
				}
			}
		}

		return nil
	})
}

func (repo *CollectionRepository) DeleteFolder(ctx context.Context, folderID string) error {
	result := repo.db.WithContext(ctx).Where("id = ?", folderID).Delete(&FolderModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return collection.ErrFolderNotFound
	}

	if err := repo.db.WithContext(ctx).Where("collection_id = ?", folderID).Delete(&CollectionEnvironmentPolicyModel{}).Error; err != nil {
		return err
	}

	if err := repo.db.WithContext(ctx).Where("collection_id = ?", folderID).Delete(&RequestModel{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *CollectionRepository) DeleteRequest(ctx context.Context, requestID string) error {
	result := repo.db.WithContext(ctx).Where("id = ?", requestID).Delete(&RequestModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return collection.ErrRequestNotFound
	}

	return nil
}

func toDomainFolder(model FolderModel) collection.Folder {
	return collection.Folder{
		ID:          model.ID,
		WorkspaceID: model.WorkspaceID,
		Name:        model.Name,
		Icon:        model.Icon,
		SortOrder:   model.SortOrder,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func toDomainRequest(model RequestModel) collection.Request {
	return collection.Request{
		ID:              model.ID,
		CollectionID:    model.CollectionID,
		WorkspaceID:     model.WorkspaceID,
		Method:          model.Method,
		Name:            model.Name,
		Path:            model.Path,
		QueryParamsJSON: defaultJSON(model.QueryParamsJSON, "[]"),
		HeadersJSON:     defaultJSON(model.HeadersJSON, "[]"),
		Body:            model.Body,
		BodyLanguage:    defaultJSON(model.BodyLanguage, "json"),
		AuthConfigJSON:  defaultJSON(model.AuthConfigJSON, "{}"),
		SortOrder:       model.SortOrder,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
	}
}

func defaultJSON(value string, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]bool, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	return result
}

func containsString(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}

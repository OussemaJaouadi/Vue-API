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
	if err := repo.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Order("sort_order ASC, name ASC").Find(&models).Error; err != nil {
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

func (repo *CollectionRepository) CreateFolder(ctx context.Context, params collection.CreateFolderParams) (collection.Folder, error) {
	var count int64
	if err := repo.db.WithContext(ctx).Model(&FolderModel{}).Where("workspace_id = ? AND name = ?", params.WorkspaceID, params.Name).Count(&count).Error; err != nil {
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

func (repo *CollectionRepository) CreateRequest(ctx context.Context, params collection.CreateRequestParams) (collection.Request, error) {
	model := RequestModel{
		ID:           id.NewUUIDV7(),
		CollectionID: params.CollectionID,
		WorkspaceID:  params.WorkspaceID,
		Method:       params.Method,
		Name:         params.Name,
		Path:         params.Path,
		SortOrder:    0,
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

	if err := repo.db.WithContext(ctx).Save(&model).Error; err != nil {
		return collection.Request{}, err
	}

	return toDomainRequest(model), nil
}

func (repo *CollectionRepository) DeleteFolder(ctx context.Context, folderID string) error {
	result := repo.db.WithContext(ctx).Where("id = ?", folderID).Delete(&FolderModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return collection.ErrFolderNotFound
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
		ID:           model.ID,
		CollectionID: model.CollectionID,
		WorkspaceID:  model.WorkspaceID,
		Method:       model.Method,
		Name:         model.Name,
		Path:         model.Path,
		SortOrder:    model.SortOrder,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}
}

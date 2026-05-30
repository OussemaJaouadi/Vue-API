package collection

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const WorkbenchExportSchema = "vue-api-workbench.collection.v1"

var (
	ErrUnsupportedImportFormat = errors.New("unsupported import format")
	ErrInvalidImportPayload    = errors.New("invalid import payload")
)

type ImportResult struct {
	Format             string
	CollectionsCreated int
	RequestsCreated    int
	Warnings           []string
}

type workbenchExport struct {
	Schema       string                `json:"schema"`
	Collections  []workbenchCollection `json:"collections"`
	RootRequests []workbenchRequest    `json:"rootRequests"`
}

type workbenchCollection struct {
	Name     string             `json:"name"`
	Icon     string             `json:"icon"`
	Requests []workbenchRequest `json:"requests"`
}

type workbenchRequest struct {
	Method       string          `json:"method"`
	Name         string          `json:"name"`
	Path         string          `json:"path"`
	QueryParams  json.RawMessage `json:"queryParams"`
	Headers      json.RawMessage `json:"headers"`
	Body         string          `json:"body"`
	BodyLanguage string          `json:"bodyLanguage"`
	AuthConfig   json.RawMessage `json:"authConfig"`
}

func ImportWorkbenchExport(ctx context.Context, repo Repository, workspaceID string, payload json.RawMessage) (ImportResult, error) {
	if repo == nil {
		return ImportResult{}, errors.New("collection repository is required")
	}
	if strings.TrimSpace(workspaceID) == "" {
		return ImportResult{}, ErrInvalidImportPayload
	}

	var envelope struct {
		Schema string `json:"schema"`
	}
	if err := json.Unmarshal(payload, &envelope); err != nil {
		return ImportResult{}, ErrInvalidImportPayload
	}
	if envelope.Schema != WorkbenchExportSchema {
		return ImportResult{}, ErrUnsupportedImportFormat
	}

	var input workbenchExport
	if err := json.Unmarshal(payload, &input); err != nil {
		return ImportResult{}, ErrInvalidImportPayload
	}

	result := ImportResult{
		Format: "Workbench export",
	}

	existingNames, err := folderNameSet(ctx, repo, workspaceID)
	if err != nil {
		return ImportResult{}, err
	}

	orderGroups := make([]RequestOrderGroup, 0, len(input.Collections)+1)
	for _, importedCollection := range input.Collections {
		name := uniqueFolderName(cleanCollectionName(importedCollection.Name), existingNames)
		folder, err := repo.CreateFolder(ctx, CreateFolderParams{
			WorkspaceID: workspaceID,
			Name:        name,
			Icon:        importedCollection.Icon,
		})
		if err != nil {
			return ImportResult{}, err
		}
		existingNames[folder.Name] = true
		result.CollectionsCreated++

		requestIDs, created, err := importRequests(ctx, repo, workspaceID, &folder.ID, importedCollection.Requests)
		if err != nil {
			return ImportResult{}, err
		}
		result.RequestsCreated += created
		orderGroups = append(orderGroups, RequestOrderGroup{
			CollectionID: &folder.ID,
			RequestIDs:   requestIDs,
		})
	}

	rootRequestIDs, created, err := importRequests(ctx, repo, workspaceID, nil, input.RootRequests)
	if err != nil {
		return ImportResult{}, err
	}
	result.RequestsCreated += created
	if len(rootRequestIDs) > 0 {
		orderGroups = append(orderGroups, RequestOrderGroup{
			CollectionID: nil,
			RequestIDs:   rootRequestIDs,
		})
	}

	if len(orderGroups) > 0 {
		if err := repo.ReorderRequests(ctx, workspaceID, orderGroups); err != nil {
			return ImportResult{}, err
		}
	}

	return result, nil
}

func folderNameSet(ctx context.Context, repo Repository, workspaceID string) (map[string]bool, error) {
	folders, err := repo.ListFolders(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	names := make(map[string]bool, len(folders))
	for _, folder := range folders {
		names[folder.Name] = true
	}

	return names, nil
}

func importRequests(ctx context.Context, repo Repository, workspaceID string, collectionID *string, requests []workbenchRequest) ([]string, int, error) {
	requestIDs := make([]string, 0, len(requests))
	for _, importedRequest := range requests {
		created, err := repo.CreateRequest(ctx, CreateRequestParams{
			CollectionID:    collectionID,
			WorkspaceID:     workspaceID,
			Method:          cleanMethod(importedRequest.Method),
			Name:            cleanRequestName(importedRequest.Name, importedRequest.Path),
			Path:            importedRequest.Path,
			QueryParamsJSON: rawOrDefault(importedRequest.QueryParams, "[]"),
			HeadersJSON:     rawOrDefault(importedRequest.Headers, "[]"),
			Body:            importedRequest.Body,
			BodyLanguage:    stringOrDefault(importedRequest.BodyLanguage, "json"),
			AuthConfigJSON:  rawOrDefault(importedRequest.AuthConfig, "{}"),
		})
		if err != nil {
			return nil, 0, err
		}
		requestIDs = append(requestIDs, created.ID)
	}

	return requestIDs, len(requestIDs), nil
}

func cleanCollectionName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "Imported Collection"
	}

	return name
}

func cleanRequestName(name string, path string) string {
	name = strings.TrimSpace(name)
	if name != "" {
		return name
	}
	path = strings.TrimSpace(path)
	if path != "" {
		return path
	}

	return "Imported Request"
}

func cleanMethod(method string) string {
	method = strings.ToUpper(strings.TrimSpace(method))
	if method == "" {
		return "GET"
	}

	return method
}

func uniqueFolderName(base string, existing map[string]bool) string {
	if !existing[base] {
		return base
	}

	for suffix := 2; ; suffix++ {
		candidate := fmt.Sprintf("%s (import %d)", base, suffix)
		if !existing[candidate] {
			return candidate
		}
	}
}

func rawOrDefault(raw json.RawMessage, fallback string) string {
	if len(raw) == 0 || string(raw) == "null" {
		return fallback
	}

	return string(raw)
}

func stringOrDefault(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}

	return value
}

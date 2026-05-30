package collection

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"
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

type ImportPreview struct {
	FileName string   `json:"fileName"`
	Format   string   `json:"format"`
	Status   string   `json:"status"`
	Summary  string   `json:"summary"`
	Details  []string `json:"details"`
}

type workbenchExport struct {
	Schema       string                `json:"schema"`
	ExportedAt   string                `json:"exportedAt,omitempty"`
	Collections  []workbenchCollection `json:"collections"`
	RootRequests []workbenchRequest    `json:"rootRequests"`
}

type workbenchCollection struct {
	ID        string             `json:"id,omitempty"`
	Name      string             `json:"name"`
	Icon      string             `json:"icon"`
	SortOrder int                `json:"sortOrder"`
	Requests  []workbenchRequest `json:"requests"`
}

type workbenchRequest struct {
	ID           string          `json:"id,omitempty"`
	Method       string          `json:"method"`
	Name         string          `json:"name"`
	Path         string          `json:"path"`
	QueryParams  json.RawMessage `json:"queryParams"`
	Headers      json.RawMessage `json:"headers"`
	Body         string          `json:"body"`
	BodyLanguage string          `json:"bodyLanguage"`
	AuthConfig   json.RawMessage `json:"authConfig"`
	SortOrder    int             `json:"sortOrder"`
}

type openAPIImport struct {
	OpenAPI string                    `json:"openapi"`
	Info    openAPIInfo               `json:"info"`
	Paths   map[string]map[string]any `json:"paths"`
}

type openAPIInfo struct {
	Title string `json:"title"`
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

func ImportContent(ctx context.Context, repo Repository, workspaceID string, payload json.RawMessage) (ImportResult, error) {
	var envelope map[string]any
	if err := json.Unmarshal(payload, &envelope); err != nil {
		return ImportResult{}, ErrInvalidImportPayload
	}

	if envelope["schema"] == WorkbenchExportSchema {
		return ImportWorkbenchExport(ctx, repo, workspaceID, payload)
	}
	if version, ok := envelope["openapi"].(string); ok && version != "" {
		return importOpenAPI(ctx, repo, workspaceID, payload)
	}

	return ImportResult{}, ErrUnsupportedImportFormat
}

func PreviewImportContent(fileName string, content string) ImportPreview {
	extension := strings.ToLower(filepath.Ext(fileName))
	if extension == ".yaml" || extension == ".yml" {
		return ImportPreview{
			FileName: fileName,
			Format:   "YAML spec",
			Status:   "unsupported",
			Summary:  "YAML upload selected",
			Details: []string{
				"YAML parsing belongs in the backend parser service.",
				"The UI accepts the file type so the flow is visible now.",
			},
		}
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(content), &payload); err != nil {
		return ImportPreview{
			FileName: fileName,
			Format:   "Invalid JSON",
			Status:   "error",
			Summary:  "Could not parse file",
			Details:  []string{"Check that the file is valid JSON before importing."},
		}
	}

	if payload["schema"] == WorkbenchExportSchema {
		collections := arrayValue(payload["collections"])
		rootRequests := arrayValue(payload["rootRequests"])
		requests := len(rootRequests)
		for _, collectionValue := range collections {
			collectionObject, ok := collectionValue.(map[string]any)
			if !ok {
				continue
			}
			requests += len(arrayValue(collectionObject["requests"]))
		}

		return ImportPreview{
			FileName: fileName,
			Format:   "Workbench export",
			Status:   "ready",
			Summary:  fmt.Sprintf("%d collections / %d requests", len(collections), requests),
			Details: []string{
				"Ready for backend persistence.",
				"Collections and request order will be stored in the selected workspace.",
			},
		}
	}

	if version, ok := payload["openapi"].(string); ok && version != "" {
		paths := objectValue(payload["paths"])
		return ImportPreview{
			FileName: fileName,
			Format:   fmt.Sprintf("OpenAPI %s", version),
			Status:   "ready",
			Summary:  fmt.Sprintf("%d paths detected", len(paths)),
			Details: []string{
				"Backend detected the spec shape.",
				"Operations will be imported into one collection with request rows.",
			},
		}
	}

	if payload["swagger"] == "2.0" {
		paths := objectValue(payload["paths"])
		return ImportPreview{
			FileName: fileName,
			Format:   "Swagger 2.0",
			Status:   "unsupported",
			Summary:  fmt.Sprintf("%d paths detected", len(paths)),
			Details: []string{
				"Backend detected the legacy spec shape.",
				"Parser adapter will normalize it before import.",
			},
		}
	}

	if _, ok := payload["item"]; ok {
		items := arrayValue(payload["item"])
		return ImportPreview{
			FileName: fileName,
			Format:   "Postman collection",
			Status:   "unsupported",
			Summary:  fmt.Sprintf("%d top-level items detected", len(items)),
			Details: []string{
				"Backend detected a Postman collection.",
				"Postman normalization is documented as a later parser adapter.",
			},
		}
	}

	return ImportPreview{
		FileName: fileName,
		Format:   "Unknown JSON",
		Status:   "error",
		Summary:  "No supported collection shape detected",
		Details:  []string{"Expected Workbench export, OpenAPI, Swagger, or Postman collection JSON."},
	}
}

func arrayValue(value any) []any {
	values, ok := value.([]any)
	if !ok {
		return nil
	}

	return values
}

func objectValue(value any) map[string]any {
	values, ok := value.(map[string]any)
	if !ok {
		return nil
	}

	return values
}

func ExportWorkbenchExport(ctx context.Context, repo Repository, workspaceID string, exportedAt time.Time) (workbenchExport, error) {
	if repo == nil {
		return workbenchExport{}, errors.New("collection repository is required")
	}
	if strings.TrimSpace(workspaceID) == "" {
		return workbenchExport{}, ErrInvalidImportPayload
	}

	folders, err := repo.ListFolders(ctx, workspaceID)
	if err != nil {
		return workbenchExport{}, err
	}

	output := workbenchExport{
		Schema:       WorkbenchExportSchema,
		ExportedAt:   exportedAt.UTC().Format(time.RFC3339),
		Collections:  make([]workbenchCollection, 0, len(folders)),
		RootRequests: []workbenchRequest{},
	}

	for _, folder := range folders {
		requests, err := repo.ListRequests(ctx, workspaceID, &folder.ID)
		if err != nil {
			return workbenchExport{}, err
		}

		output.Collections = append(output.Collections, workbenchCollection{
			ID:        folder.ID,
			Name:      folder.Name,
			Icon:      folder.Icon,
			SortOrder: folder.SortOrder,
			Requests:  exportRequests(requests),
		})
	}

	rootRequests, err := repo.ListRootRequests(ctx, workspaceID)
	if err != nil {
		return workbenchExport{}, err
	}
	output.RootRequests = exportRequests(rootRequests)

	return output, nil
}

func importOpenAPI(ctx context.Context, repo Repository, workspaceID string, payload json.RawMessage) (ImportResult, error) {
	if repo == nil {
		return ImportResult{}, errors.New("collection repository is required")
	}
	if strings.TrimSpace(workspaceID) == "" {
		return ImportResult{}, ErrInvalidImportPayload
	}

	var input openAPIImport
	if err := json.Unmarshal(payload, &input); err != nil {
		return ImportResult{}, ErrInvalidImportPayload
	}
	if input.OpenAPI == "" {
		return ImportResult{}, ErrUnsupportedImportFormat
	}

	existingNames, err := folderNameSet(ctx, repo, workspaceID)
	if err != nil {
		return ImportResult{}, err
	}

	collectionName := uniqueFolderName(cleanCollectionName(input.Info.Title), existingNames)
	folder, err := repo.CreateFolder(ctx, CreateFolderParams{
		WorkspaceID: workspaceID,
		Name:        collectionName,
		Icon:        "PhGlobe",
	})
	if err != nil {
		return ImportResult{}, err
	}

	requests := openAPIRequests(input)
	requestIDs, created, err := importRequests(ctx, repo, workspaceID, &folder.ID, requests)
	if err != nil {
		return ImportResult{}, err
	}
	if len(requestIDs) > 0 {
		if err := repo.ReorderRequests(ctx, workspaceID, []RequestOrderGroup{{
			CollectionID: &folder.ID,
			RequestIDs:   requestIDs,
		}}); err != nil {
			return ImportResult{}, err
		}
	}

	return ImportResult{
		Format:             fmt.Sprintf("OpenAPI %s", input.OpenAPI),
		CollectionsCreated: 1,
		RequestsCreated:    created,
		Warnings:           []string{},
	}, nil
}

func openAPIRequests(input openAPIImport) []workbenchRequest {
	paths := make([]string, 0, len(input.Paths))
	for path := range input.Paths {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	requests := make([]workbenchRequest, 0)
	for _, path := range paths {
		pathItem := input.Paths[path]
		for _, method := range []string{"get", "post", "put", "patch", "delete", "options", "head"} {
			operation := objectValue(pathItem[method])
			if operation == nil {
				continue
			}

			requests = append(requests, workbenchRequest{
				Method:       strings.ToUpper(method),
				Name:         openAPIOperationName(method, path, operation),
				Path:         path,
				QueryParams:  openAPIQueryParams(operation),
				Headers:      json.RawMessage("[]"),
				Body:         openAPIRequestBody(operation),
				BodyLanguage: "json",
				AuthConfig:   json.RawMessage("{}"),
			})
		}
	}

	return requests
}

func openAPIOperationName(method string, path string, operation map[string]any) string {
	for _, key := range []string{"summary", "operationId"} {
		if value, ok := operation[key].(string); ok && strings.TrimSpace(value) != "" {
			return value
		}
	}

	return fmt.Sprintf("%s %s", strings.ToUpper(method), path)
}

func openAPIQueryParams(operation map[string]any) json.RawMessage {
	parameters := arrayValue(operation["parameters"])
	rows := make([]map[string]any, 0)
	for _, parameterValue := range parameters {
		parameter := objectValue(parameterValue)
		if parameter["in"] != "query" {
			continue
		}
		name, _ := parameter["name"].(string)
		if strings.TrimSpace(name) == "" {
			continue
		}
		rows = append(rows, map[string]any{
			"id":      fmt.Sprintf("openapi-query-%s", name),
			"key":     name,
			"value":   "",
			"enabled": false,
		})
	}

	return mustJSON(rows, "[]")
}

func openAPIRequestBody(operation map[string]any) string {
	requestBody := objectValue(operation["requestBody"])
	content := objectValue(requestBody["content"])
	if _, ok := content["application/json"]; ok {
		return "{}"
	}

	return ""
}

func mustJSON(value any, fallback string) json.RawMessage {
	raw, err := json.Marshal(value)
	if err != nil {
		return json.RawMessage(fallback)
	}

	return raw
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

func exportRequests(requests []Request) []workbenchRequest {
	output := make([]workbenchRequest, 0, len(requests))
	for _, request := range requests {
		output = append(output, workbenchRequest{
			ID:           request.ID,
			Method:       request.Method,
			Name:         request.Name,
			Path:         request.Path,
			QueryParams:  json.RawMessage(rawOrDefault(json.RawMessage(request.QueryParamsJSON), "[]")),
			Headers:      json.RawMessage(rawOrDefault(json.RawMessage(request.HeadersJSON), "[]")),
			Body:         request.Body,
			BodyLanguage: stringOrDefault(request.BodyLanguage, "json"),
			AuthConfig:   json.RawMessage(rawOrDefault(json.RawMessage(request.AuthConfigJSON), "{}")),
			SortOrder:    request.SortOrder,
		})
	}

	return output
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

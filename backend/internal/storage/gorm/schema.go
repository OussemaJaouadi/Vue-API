package gormstorage

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"gorm.io/gorm"
)

type ChangeSafety string

const (
	SafeChange   ChangeSafety = "safe"
	ManualChange ChangeSafety = "manual"
)

type Change struct {
	Safety      ChangeSafety
	Description string
	SQL         string
}

type SchemaPlan struct {
	Changes []Change
}

var ErrManualMigrationRequired = errors.New("manual migration required")
var ErrSchemaOutOfDate = errors.New("database schema is out of date")

func (plan SchemaPlan) HasChanges() bool {
	return len(plan.Changes) > 0
}

func (plan SchemaPlan) HasManualChanges() bool {
	for _, change := range plan.Changes {
		if change.Safety == ManualChange {
			return true
		}
	}

	return false
}

func (plan SchemaPlan) SafeChanges() []Change {
	changes := make([]Change, 0, len(plan.Changes))
	for _, change := range plan.Changes {
		if change.Safety == SafeChange {
			changes = append(changes, change)
		}
	}

	return changes
}

func (plan SchemaPlan) String() string {
	if len(plan.Changes) == 0 {
		return "Database schema is up to date."
	}

	var builder strings.Builder
	builder.WriteString("Pending database changes:\n")
	for _, change := range plan.Changes {
		builder.WriteString("- ")
		if change.Safety == ManualChange {
			builder.WriteString("[manual] ")
		}
		builder.WriteString(change.Description)
		builder.WriteString("\n")
	}

	return strings.TrimRight(builder.String(), "\n")
}

func Plan(db *gorm.DB) (SchemaPlan, error) {
	actual, err := inspectSchema(db)
	if err != nil {
		return SchemaPlan{}, err
	}

	plan := SchemaPlan{}
	for _, table := range expectedSchema() {
		actualTable, exists := actual.Tables[table.Name]
		if !exists {
			plan.Changes = append(plan.Changes, Change{
				Safety:      SafeChange,
				Description: "create table " + table.Name,
				SQL:         table.CreateSQL,
			})
			for _, index := range table.Indexes {
				plan.Changes = append(plan.Changes, Change{
					Safety:      SafeChange,
					Description: createIndexDescription(index),
					SQL:         index.CreateSQL,
				})
			}
			continue
		}

		for _, column := range table.Columns {
			actualColumn, exists := actualTable.Columns[column.Name]
			if !exists {
				plan.Changes = append(plan.Changes, Change{
					Safety:      SafeChange,
					Description: fmt.Sprintf("%s.%s: add column %s", table.Name, column.Name, column.Definition),
					SQL:         fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s", quoteIdent(table.Name), column.Definition),
				})
				continue
			}

			if !column.matches(actualColumn) {
				plan.Changes = append(plan.Changes, Change{
					Safety:      ManualChange,
					Description: fmt.Sprintf("%s.%s: expected %s, found %s", table.Name, column.Name, column.ExpectedDescription(), actualColumn.Description()),
				})
			}
		}

		for _, index := range table.Indexes {
			actualIndex, exists := actualTable.Indexes[index.Name]
			if !exists {
				plan.Changes = append(plan.Changes, Change{
					Safety:      SafeChange,
					Description: createIndexDescription(index),
					SQL:         index.CreateSQL,
				})
				continue
			}

			if !index.matches(actualIndex) {
				plan.Changes = append(plan.Changes, Change{
					Safety: ManualChange,
					Description: fmt.Sprintf("%s: expected %s index on %s(%s), found %s index on %s(%s)",
						index.Name,
						uniqueLabel(index.Unique),
						table.Name,
						strings.Join(index.Columns, ", "),
						uniqueLabel(actualIndex.Unique),
						table.Name,
						strings.Join(actualIndex.Columns, ", "),
					),
				})
			}
		}
	}

	return plan, nil
}

func Migrate(db *gorm.DB) error {
	plan, err := Plan(db)
	if err != nil {
		return err
	}
	if plan.HasManualChanges() {
		return fmt.Errorf("%w:\n%s", ErrManualMigrationRequired, plan.String())
	}

	for _, change := range plan.SafeChanges() {
		if change.SQL == "" {
			continue
		}
		if err := db.Exec(change.SQL).Error; err != nil {
			return err
		}
	}

	return nil
}

func VerifySchema(db *gorm.DB) error {
	plan, err := Plan(db)
	if err != nil {
		return err
	}
	if plan.HasChanges() {
		return fmt.Errorf("%w:\n%s", ErrSchemaOutOfDate, plan.String())
	}

	return nil
}

type expectedTable struct {
	Name      string
	CreateSQL string
	Columns   []expectedColumn
	Indexes   []expectedIndex
}

type expectedColumn struct {
	Name       string
	Type       string
	NotNull    bool
	PrimaryKey bool
	Default    string
	Definition string
}

type expectedIndex struct {
	Name      string
	Unique    bool
	Columns   []string
	CreateSQL string
}

func expectedSchema() []expectedTable {
	return []expectedTable{
		{
			Name: "users",
			CreateSQL: strings.Join([]string{
				"CREATE TABLE users (",
				"id text NOT NULL PRIMARY KEY,",
				"email text NOT NULL,",
				"username text NOT NULL,",
				"password_hash text NOT NULL,",
				"global_role text NOT NULL,",
				"token_version integer NOT NULL DEFAULT 1,",
				"active numeric NOT NULL DEFAULT true,",
				"created_at datetime NOT NULL,",
				"updated_at datetime NOT NULL",
				")",
			}, " "),
			Columns: []expectedColumn{
				{Name: "id", Type: "text", NotNull: true, PrimaryKey: true, Definition: "id text NOT NULL PRIMARY KEY"},
				{Name: "email", Type: "text", NotNull: true, Definition: "email text NOT NULL"},
				{Name: "username", Type: "text", NotNull: true, Definition: "username text NOT NULL"},
				{Name: "password_hash", Type: "text", NotNull: true, Definition: "password_hash text NOT NULL"},
				{Name: "global_role", Type: "text", NotNull: true, Definition: "global_role text NOT NULL"},
				{Name: "token_version", Type: "integer", NotNull: true, Default: "1", Definition: "token_version integer NOT NULL DEFAULT 1"},
				{Name: "active", Type: "numeric", NotNull: true, Default: "true", Definition: "active numeric NOT NULL DEFAULT true"},
				{Name: "created_at", Type: "datetime", NotNull: true, Definition: "created_at datetime NOT NULL"},
				{Name: "updated_at", Type: "datetime", NotNull: true, Definition: "updated_at datetime NOT NULL"},
			},
			Indexes: []expectedIndex{
				{Name: "idx_users_email", Unique: true, Columns: []string{"email"}, CreateSQL: "CREATE UNIQUE INDEX idx_users_email ON users(email)"},
				{Name: "idx_users_username", Unique: true, Columns: []string{"username"}, CreateSQL: "CREATE UNIQUE INDEX idx_users_username ON users(username)"},
			},
		},
		{
			Name: "workspaces",
			CreateSQL: strings.Join([]string{
				"CREATE TABLE workspaces (",
				"id text NOT NULL PRIMARY KEY,",
				"name text NOT NULL,",
				"created_by_user_id text NOT NULL,",
				"created_at datetime NOT NULL,",
				"updated_at datetime NOT NULL",
				")",
			}, " "),
			Columns: []expectedColumn{
				{Name: "id", Type: "text", NotNull: true, PrimaryKey: true, Definition: "id text NOT NULL PRIMARY KEY"},
				{Name: "name", Type: "text", NotNull: true, Definition: "name text NOT NULL"},
				{Name: "created_by_user_id", Type: "text", NotNull: true, Definition: "created_by_user_id text NOT NULL"},
				{Name: "created_at", Type: "datetime", NotNull: true, Definition: "created_at datetime NOT NULL"},
				{Name: "updated_at", Type: "datetime", NotNull: true, Definition: "updated_at datetime NOT NULL"},
			},
		},
		{
			Name: "workspace_memberships",
			CreateSQL: strings.Join([]string{
				"CREATE TABLE workspace_memberships (",
				"id text NOT NULL PRIMARY KEY,",
				"workspace_id text NOT NULL,",
				"user_id text NOT NULL,",
				"role text NOT NULL,",
				"created_by_user_id text NOT NULL,",
				"created_at datetime NOT NULL,",
				"updated_at datetime NOT NULL",
				")",
			}, " "),
			Columns: []expectedColumn{
				{Name: "id", Type: "text", NotNull: true, PrimaryKey: true, Definition: "id text NOT NULL PRIMARY KEY"},
				{Name: "workspace_id", Type: "text", NotNull: true, Definition: "workspace_id text NOT NULL"},
				{Name: "user_id", Type: "text", NotNull: true, Definition: "user_id text NOT NULL"},
				{Name: "role", Type: "text", NotNull: true, Definition: "role text NOT NULL"},
				{Name: "created_by_user_id", Type: "text", NotNull: true, Definition: "created_by_user_id text NOT NULL"},
				{Name: "created_at", Type: "datetime", NotNull: true, Definition: "created_at datetime NOT NULL"},
				{Name: "updated_at", Type: "datetime", NotNull: true, Definition: "updated_at datetime NOT NULL"},
			},
			Indexes: []expectedIndex{
				{
					Name:      "idx_workspace_memberships_workspace_user",
					Unique:    true,
					Columns:   []string{"workspace_id", "user_id"},
					CreateSQL: "CREATE UNIQUE INDEX idx_workspace_memberships_workspace_user ON workspace_memberships(workspace_id, user_id)",
				},
			},
		},
	}
}

type inspectedSchema struct {
	Tables map[string]inspectedTable
}

type inspectedTable struct {
	Columns map[string]inspectedColumn
	Indexes map[string]inspectedIndex
}

type inspectedColumn struct {
	Name       string
	Type       string
	NotNull    bool
	PrimaryKey bool
	Default    string
}

type inspectedIndex struct {
	Name    string
	Unique  bool
	Columns []string
}

func inspectSchema(db *gorm.DB) (inspectedSchema, error) {
	var tableRows []struct {
		Name string
	}
	if err := db.Raw("SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%'").Scan(&tableRows).Error; err != nil {
		return inspectedSchema{}, err
	}

	schema := inspectedSchema{Tables: make(map[string]inspectedTable)}
	for _, tableRow := range tableRows {
		table := inspectedTable{
			Columns: make(map[string]inspectedColumn),
			Indexes: make(map[string]inspectedIndex),
		}

		var columnRows []struct {
			CID       int     `gorm:"column:cid"`
			Name      string  `gorm:"column:name"`
			Type      string  `gorm:"column:type"`
			NotNull   int     `gorm:"column:notnull"`
			Default   *string `gorm:"column:dflt_value"`
			PrimaryID int     `gorm:"column:pk"`
		}
		if err := db.Raw("PRAGMA table_info(" + quoteIdent(tableRow.Name) + ")").Scan(&columnRows).Error; err != nil {
			return inspectedSchema{}, err
		}
		for _, columnRow := range columnRows {
			defaultValue := ""
			if columnRow.Default != nil {
				defaultValue = normalizeDefault(*columnRow.Default)
			}
			table.Columns[columnRow.Name] = inspectedColumn{
				Name:       columnRow.Name,
				Type:       normalizeType(columnRow.Type),
				NotNull:    columnRow.NotNull == 1 || columnRow.PrimaryID > 0,
				PrimaryKey: columnRow.PrimaryID > 0,
				Default:    defaultValue,
			}
		}

		var indexRows []struct {
			Name   string `gorm:"column:name"`
			Unique int    `gorm:"column:unique"`
		}
		if err := db.Raw("PRAGMA index_list(" + quoteIdent(tableRow.Name) + ")").Scan(&indexRows).Error; err != nil {
			return inspectedSchema{}, err
		}
		for _, indexRow := range indexRows {
			var indexInfoRows []struct {
				SeqNo int    `gorm:"column:seqno"`
				Name  string `gorm:"column:name"`
			}
			if err := db.Raw("PRAGMA index_info(" + quoteIdent(indexRow.Name) + ")").Scan(&indexInfoRows).Error; err != nil {
				return inspectedSchema{}, err
			}
			sort.Slice(indexInfoRows, func(i int, j int) bool {
				return indexInfoRows[i].SeqNo < indexInfoRows[j].SeqNo
			})

			columns := make([]string, 0, len(indexInfoRows))
			for _, indexInfoRow := range indexInfoRows {
				columns = append(columns, indexInfoRow.Name)
			}

			table.Indexes[indexRow.Name] = inspectedIndex{
				Name:    indexRow.Name,
				Unique:  indexRow.Unique == 1,
				Columns: columns,
			}
		}

		schema.Tables[tableRow.Name] = table
	}

	return schema, nil
}

func (expected expectedColumn) matches(actual inspectedColumn) bool {
	return normalizeType(expected.Type) == actual.Type &&
		expected.NotNull == actual.NotNull &&
		expected.PrimaryKey == actual.PrimaryKey &&
		normalizeDefault(expected.Default) == actual.Default
}

func (expected expectedColumn) ExpectedDescription() string {
	parts := []string{expected.Type}
	if expected.NotNull {
		parts = append(parts, "NOT NULL")
	}
	if expected.PrimaryKey {
		parts = append(parts, "PRIMARY KEY")
	}
	if expected.Default != "" {
		parts = append(parts, "DEFAULT "+expected.Default)
	}

	return strings.Join(parts, " ")
}

func (actual inspectedColumn) Description() string {
	parts := []string{actual.Type}
	if actual.NotNull {
		parts = append(parts, "NOT NULL")
	}
	if actual.PrimaryKey {
		parts = append(parts, "PRIMARY KEY")
	}
	if actual.Default != "" {
		parts = append(parts, "DEFAULT "+actual.Default)
	}

	return strings.Join(parts, " ")
}

func (expected expectedIndex) matches(actual inspectedIndex) bool {
	return expected.Unique == actual.Unique && strings.Join(expected.Columns, ",") == strings.Join(actual.Columns, ",")
}

func createIndexDescription(index expectedIndex) string {
	return fmt.Sprintf("create %sindex %s", uniquePrefix(index.Unique), index.Name)
}

func quoteIdent(identifier string) string {
	return "`" + strings.ReplaceAll(identifier, "`", "``") + "`"
}

func normalizeType(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func normalizeDefault(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.TrimPrefix(value, "(")
	value = strings.TrimSuffix(value, ")")
	return strings.Trim(value, "'\"")
}

func uniquePrefix(unique bool) string {
	if unique {
		return "unique "
	}

	return ""
}

func uniqueLabel(unique bool) string {
	if unique {
		return "unique"
	}

	return "non-unique"
}

package gormstorage

import (
	"fmt"
	"strconv"
	"sync/atomic"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"vue-api/backend/internal/config"
)

var inMemoryDBCounter uint64

func Open(database config.DatabaseConfig) (*gorm.DB, error) {
	if database.Driver != "sqlite" {
		return nil, fmt.Errorf("unsupported gorm database driver %q", database.Driver)
	}

	return gorm.Open(sqlite.Open(database.URL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
}

func openInMemoryDB() (*gorm.DB, error) {
	id := atomic.AddUint64(&inMemoryDBCounter, 1)
	return gorm.Open(sqlite.Open("file:gormschema"+strconv.FormatUint(id, 10)+"?mode=memory&cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
}

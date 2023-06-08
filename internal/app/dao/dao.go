package dao

import (
	"boilerplate/internal/pkg/config"
	"strings"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if dbType := config.C.Gorm.DBType; strings.ToLower(dbType) == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	}

	return db.AutoMigrate()
}

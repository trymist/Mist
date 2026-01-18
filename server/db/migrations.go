package db

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/corecollectives/mist/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

func migrateExistingTable(db *gorm.DB, model interface{}) error {
	migrator := db.Migrator()

	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return fmt.Errorf("failed to parse model: %w", err)
	}

	tableName := stmt.Schema.Table

	for _, field := range stmt.Schema.Fields {
		if field.DBName == "" {
			continue
		}

		if !migrator.HasColumn(model, field.DBName) {
			columnType := getSQLiteColumnType(field)
			if columnType == "" {
				continue
			}

			sql := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `%s` %s", tableName, field.DBName, columnType)

			if field.HasDefaultValue && field.DefaultValue != "" {
				sql += fmt.Sprintf(" DEFAULT %s", field.DefaultValue)
			}

			if err := db.Exec(sql).Error; err != nil {
				fmt.Printf("migration.go: warning adding column %s.%s: %v\n", tableName, field.DBName, err)
			}
		}
	}

	return nil
}

func getSQLiteColumnType(field *schema.Field) string {
	switch field.DataType {
	case schema.Bool:
		return "BOOLEAN"
	case schema.Int, schema.Uint:
		return "INTEGER"
	case schema.Float:
		return "REAL"
	case schema.String:
		return "TEXT"
	case schema.Time:
		return "DATETIME"
	case schema.Bytes:
		return "BLOB"
	default:
		kind := field.FieldType.Kind()
		if kind == reflect.Ptr {
			kind = field.FieldType.Elem().Kind()
		}
		switch kind {
		case reflect.Bool:
			return "BOOLEAN"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return "INTEGER"
		case reflect.Float32, reflect.Float64:
			return "REAL"
		case reflect.String:
			return "TEXT"
		case reflect.Struct:
			typeName := field.FieldType.String()
			if strings.Contains(typeName, "Time") || strings.Contains(typeName, "DeletedAt") {
				return "DATETIME"
			}
			return "TEXT"
		default:
			return "TEXT"
		}
	}
}

func MigrateDB(dbInstance *gorm.DB) error {
	return migrateDbInternal(dbInstance)
}

func migrateDbInternal(dbInstance *gorm.DB) error {
	migrator := dbInstance.Migrator()

	allModels := []interface{}{
		&models.User{},
		&models.ApiToken{},
		&models.App{},
		&models.AuditLog{},
		&models.Backup{},
		&models.Deployment{},
		&models.EnvVariable{},
		&models.GithubApp{},
		&models.Project{},
		&models.ProjectMember{},
		&models.GitProvider{},
		&models.GithubInstallation{},
		&models.AppRepositories{},
		&models.Domain{},
		&models.Volume{},
		&models.Cron{},
		&models.Registry{},
		&models.SystemSettingEntry{},
		&models.Logs{},
		&models.ServiceTemplate{},
		&models.Session{},
		&models.Notification{},
		&models.UpdateLog{},
	}

	for _, model := range allModels {
		if migrator.HasTable(model) {
			if err := migrateExistingTable(dbInstance, model); err != nil {
				fmt.Printf("migration.go: warning migrating existing table: %v\n", err)
			}
		} else {
			if err := dbInstance.AutoMigrate(model); err != nil {
				fmt.Printf("migration.go: error creating table: %v\n", err)
				return err
			}
		}
	}

	var wildCardDomain = models.SystemSettingEntry{
		Key:   "wildcard_domain",
		Value: " ",
	}
	var MistAppName = models.SystemSettingEntry{
		Key:   "mist_app_name",
		Value: "mist",
	}
	var Version = models.SystemSettingEntry{
		Key:   "version",
		Value: "1.0.3",
	}

	dbInstance.Clauses(clause.Insert{Modifier: "OR IGNORE"}).Create(&wildCardDomain)
	dbInstance.Clauses(clause.Insert{Modifier: "OR IGNORE"}).Create(&MistAppName)
	dbInstance.Clauses(clause.Insert{Modifier: "OR REPLACE"}).Create(&Version)

	return nil
}

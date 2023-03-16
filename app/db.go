package app

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"github.com/spf13/viper"
)

var (
	db *sql.DB
)

func connectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.user"),
		viper.GetString("db.pass"),
		viper.GetString("db.name"))
}

// InitDatabase For Vanilla SQL
func InitDatabase() (*sql.DB, error) {

	enabled := viper.GetBool("db.enabled")
	if !enabled {
		return nil, nil
	}

	driver := viper.GetString("db.driver")
	localDB, err := sql.Open(driver, connectionString())
	if err != nil {
		return nil, err
	}
	log.Info("Database connection was created")

	db = localDB

	return localDB, nil
}

func RunMigrations(rootDir ...string) error {

	enabled := viper.GetBool("db.enabled")
	if !enabled {
		return nil
	}

	basePath := "."
	if len(rootDir) != 0 {
		basePath = rootDir[0]
	}

	err := goose.Up(db, basePath+"/migrations")
	if err != nil {
		return err
	}
	return nil
}

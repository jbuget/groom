package db

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var Database *sql.DB

func InitDatabase(datasourceName string) {
	database, err := sql.Open("pgx", datasourceName)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	Database = database
}

func RunMigrations(migrationsPath string, databaseName string) error {
	driver, err := postgres.WithInstance(Database, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		databaseName,
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Migrations applied successfully")
	return nil
}

func Connect(databaseURL string) (*sql.DB, error) {
	return sql.Open("pgx", databaseURL)
}

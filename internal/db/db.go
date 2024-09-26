package db

import (
    "database/sql"
    "log"
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/jackc/pgx/v4/stdlib"
)

// Fonction pour exécuter les migrations
func RunMigrations(db *sql.DB, migrationsPath string) error {
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        return err
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://"+migrationsPath, // chemin vers les fichiers de migration
        "postgres",               // nom de la base de données
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
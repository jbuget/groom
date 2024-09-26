package db

import (
    "database/sql"
    _ "github.com/jackc/pgx/v4/stdlib"
)

func Connect(databaseURL string) (*sql.DB, error) {
    return sql.Open("pgx", databaseURL)
}
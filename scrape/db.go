package Result_NITH

import (
	"Result-NITH/db"
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"os"
)

//go:embed db/sql/schema.sql
var DDL string

func GetDbQueries() (*sql.DB, *db.Queries) {

	database, err := sql.Open("sqlite3", "./result.db")
	if err != nil {
		log.Fatal(err)
	}

	queries := db.New(database)
	return database, queries
}

func GetDbQueriesForNewDb(fileName string) (*sql.DB, *db.Queries) {
	if _, err := os.Stat(fileName); os.IsExist(err) {
		fmt.Printf("File %s already exist\n", fileName)
	}

	database, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	// create tables
	_, err = database.ExecContext(ctx, DDL)
	if err != nil {
		log.Fatal(err)
	}
	queries := db.New(database)
	return database, queries
}

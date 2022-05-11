package pingtraq

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/osm/migrator"
	"github.com/osm/migrator/repository"
)

var dbClient *sql.DB

func initDB(db string) error {
	var err error
	if dbClient, err = sql.Open("sqlite3", db); err != nil {
		return fmt.Errorf("can't initialize database connection: %v", err)
	}

	return migrator.ToLatest(dbClient, getDatabaseRepository())
}

func getDatabaseRepository() repository.Source {
	return repository.FromMemory(map[int]string{
		1: "CREATE TABLE migration (version text NOT NULL PRIMARY KEY);",
		2: "CREATE TABLE ping (id text NOT NULL PRIMARY KEY, name text NOT NULL UNIQUE, created_at timestamp with time zone NOT NULL);",
		3: "CREATE TABLE ping_record (id text PRIMARY KEY, ping_id text NOT NULL, client text NOT NULL, address text NOT NULL, user_agent text NOT NULL, battery_level text, created_at timestamp with time zone NOT NULL, FOREIGN KEY (ping_id) REFERENCES ping (id));",
		4: "ALTER TABLE ping ADD COLUMN after_hook text;",
	})
}

func prepare(query string) (*sql.Stmt, error) {
	return dbClient.Prepare(query)
}

func query(query string, args ...interface{}) (*sql.Rows, error) {
	return dbClient.Query(query, args...)
}

func queryRow(query string, args ...interface{}) *sql.Row {
	return dbClient.QueryRow(query, args...)
}

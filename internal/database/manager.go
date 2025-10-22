package database

import (
	"cmd/obyavigo/main.go/internal/config"
	"cmd/obyavigo/main.go/internal/database/mongodb"
	"cmd/obyavigo/main.go/internal/database/postgres"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type DB struct {
	Psql  *postgres.Postgres
	Mongo *mongodb.Mongo
}

func Init(cfg *config.Config) (*DB, error) {
	var q map[string]string
	queries, err := os.ReadFile(cfg.Database.Queries)
	if err != nil {
		return nil, fmt.Errorf("error open queries file: %w", err)
	}

	if err = json.Unmarshal(queries, &q); err != nil {
		return nil, fmt.Errorf("error parsing queries file: %w", err)
	}

	psql, err := postgres.Connect(cfg, q)
	if err != nil {
		return nil, err
	}

	mongo, err := mongodb.Connect(cfg, q)
	if err != nil {
		return nil, err
	}

	return &DB{
		Psql:  psql,
		Mongo: mongo,
	}, nil
}

func (db *DB) CheckConnections() error {
	if err := db.Mongo.Ping(); err != nil {
		return err
	}

	if err := db.Psql.Ping(); err != nil {
		return err
	}

	return nil
}

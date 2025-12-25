package main

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	// Open a connection to ClickHouse
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{"127.0.0.1:19000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "otus",
			Password: "otus",
		},
		Settings:    clickhouse.Settings{"max_execution_time": 60},
		DialTimeout: time.Second * 30,
	})
	defer conn.Close()

	// Create a table
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS users (
			id UInt32,
			name String
		) ENGINE = MergeTree()
		ORDER BY id
	`

	_, err := conn.ExecContext(ctx, createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	// Start a transaction
	tx, err := conn.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback() // Rollback if we encounter an error

	// Prepare the insert statement within the transaction
	insertSQL := `
		INSERT INTO users (id, name) VALUES (?, ?)
	`

	// Create a prepared statement for the insert
	stmt, err := tx.PrepareContext(ctx, insertSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Define data to be inserted
	data := [][]interface{}{
		{1, "John"},
		{2, "Jane"},
		{3, "Doe"},
	}

	// Batch insert data
	for _, row := range data {
		_, err = stmt.ExecContext(ctx, row...)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	// Query data from the table
	querySQL := `
		SELECT * FROM users
	`

	rows, err := conn.QueryContext(ctx, querySQL)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the result rows
	for rows.Next() {
		var id uint32
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

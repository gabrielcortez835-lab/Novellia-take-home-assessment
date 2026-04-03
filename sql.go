package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	FilePath string
	DB       *sql.DB
}

func NewDatabase(filePath string) (*Database, error) {
	dbObj := &Database{FilePath: filePath}

	// Step 1: Check if DB file exists
	_, err := os.Stat(filePath)
	dbExists := err == nil

	// Step 2: Open a connection (it will create file if not exists)
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	dbObj.DB = db

	// Step 3: If not exists, initialize schema
	if !dbExists {
		log.Println("Database file does not exist. Creating and initializing...")
		if err := dbObj.initializeSchema(); err != nil {
			return nil, fmt.Errorf("failed to initialize schema: %w", err)
		}
	} else {
		log.Println("Database file exists. Running upgrade/migration...")
		if err := dbObj.upgradeSchema(); err != nil {
			return nil, fmt.Errorf("failed to upgrade schema: %w", err)
		}
	}

	return dbObj, nil
}

// Connect opens (and creates if not exists) a SQLite DB file
func Connect(filePath string) (*sql.DB, error) {
	// Check if file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Println("Database file does not exist — will create:", filePath)
	}

	// Open connection — automatically creates file if not exists
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite DB: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping SQLite DB: %w", err)
	}

	return db, nil
}

// Simulates "CREATE or UPDATE stored proc" using Go code to run SQL
func (d *Database) initializeSchema() error {
	// Example: create a table
	schema := fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
        id TEXT NOT NULL PRIMARY KEY,
        resourceType TEXT NOT NULL,
        subject TEXT NOT NULL,
        metadata TEXT
    );
`, FHIRDataTableName)
	_, err := d.DB.Exec(schema)
	return err
}

// In real DBs that support stored procs, you’d call the proc directly.
// Here we'd apply migrations (for SQLite, just altering tables if needed)
func (d *Database) upgradeSchema() error {
	// Example: ensure a column exists
	addColumn := `
    ALTER TABLE users ADD COLUMN email TEXT;
    `
	_, err := d.DB.Exec(addColumn)
	if err != nil {
		// Ignore "duplicate column" errors in simplistic migration logic.
		log.Println("Upgrade note:", err)
	}
	return nil
}

// InsertResource inserts a new record into resources table
func InsertResource(db *sql.DB, id, resourceType, subject, metadata string) (int64, error) {
	stmt, err := db.Prepare(fmt.Sprintf(`
    INSERT INTO %s (id, resourceType, subject, metadata)
        VALUES (?, ?, ?, ?);
    );
	`, FHIRDataTableName))

	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(id, resourceType, subject, metadata)
	if err != nil {
		return 0, fmt.Errorf("failed to execute insert: %w", err)
	}

	// Return the inserted ID
	insertedId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return insertedId, nil
}

// InsertResource inserts a new record into resources table
func GetResourceById(db *sql.DB, id string) (string, error) {
	query := fmt.Sprintf(`
    SELECT metadata
    FROM %s
    WHERE id = ?;
    `, FHIRDataTableName)
	row := db.QueryRow(query, id)

	var r string
	err := row.Scan(&r)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // Not found
		}
		return "", fmt.Errorf("failed to scan row: %w", err)
	}

	return r, nil
}

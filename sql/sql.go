package sql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gabrielcortez835-lab/Novellia-take-home-assessment/constants"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tidwall/gjson"
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

// SqlConnect opens (and creates if not exists) a SQLite DB file
func SqlConnect(filePath string) (*sql.DB, error) {
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
`, constants.FHIRDataTableName)
	_, err := d.DB.Exec(schema)
	return err
}

func (d *Database) upgradeSchema() error {
	return nil
}

// SqlInsertResource inserts a new record into resources table
func SqlInsertResource(id, resourceType, subject, metadata string) error {

	db, err := SqlConnect(constants.SqlDBName)
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(fmt.Sprintf(`
    INSERT INTO %s (id, resourceType, subject, metadata)
        VALUES (?, ?, ?, ?);
    );
	`, constants.FHIRDataTableName))

	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, resourceType, subject, metadata)
	if err != nil {
		return fmt.Errorf("failed to execute insert: %w", err)
	}
	return nil
}

// InsertResource inserts a new record into resources table
func SqlGetRecordsById(id string) ([]gjson.Result, error) {
	db, err := SqlConnect(constants.SqlDBName)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	query := fmt.Sprintf(`
    SELECT metadata
    FROM %s
    WHERE id = ?;
    `, constants.FHIRDataTableName)

	return sqlGetQueryInternal(db, query, id)
}

func SqlGetRecords(resourceType string, subject string) ([]gjson.Result, error) {
	db, err := SqlConnect(constants.SqlDBName)

	if err != nil {
		return nil, err
	}
	defer db.Close()

	if resourceType != "" && subject != "" {
		return sqlGetRecordsByResourceTypeAndSubjectInternal(db, subject, subject)
	} else if resourceType != "" {
		return sqlGetRecordsByResourceTypeInternal(db, subject)
	} else if subject != "" {
		return sqlGetRecordsBySubjectInternal(db, subject)
	} else {
		return sqlGetRecordsInternal(db)
	}

}

func sqlGetRecordsInternal(db *sql.DB) ([]gjson.Result, error) {
	query := fmt.Sprintf(`
        SELECT metadata
        FROM %s;
    `, constants.FHIRDataTableName)

	return sqlGetQueryInternal(db, query)
}

func sqlGetRecordsByResourceTypeInternal(db *sql.DB, resourceType string) ([]gjson.Result, error) {
	query := fmt.Sprintf(`
        SELECT metadata
        FROM %s
        WHERE resourceType = ?;
    `, constants.FHIRDataTableName)

	return sqlGetQueryInternal(db, query, resourceType)
}

func sqlGetRecordsByResourceTypeAndSubjectInternal(db *sql.DB, resourceType string, subject string) ([]gjson.Result, error) {
	query := fmt.Sprintf(`
        SELECT metadata
        FROM %s
        WHERE resourceType = ? and subject = ?;
    `, constants.FHIRDataTableName)

	return sqlGetQueryInternal(db, query, resourceType, subject)
}

func sqlGetRecordsBySubjectInternal(db *sql.DB, subject string) ([]gjson.Result, error) {
	query := fmt.Sprintf(`
        SELECT metadata
        FROM %s
        WHERE subject = ?;
    `, constants.FHIRDataTableName)

	return sqlGetQueryInternal(db, query, subject)
}

func sqlGetQueryInternal(db *sql.DB, query string, args ...interface{}) ([]gjson.Result, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query records: %w", err)
	}
	defer rows.Close()

	var results []gjson.Result

	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		results = append(results, gjson.Parse(raw))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

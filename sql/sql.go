package sql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gabrielcortez835-lab/Novellia-take-home-assessment/constants"
	"github.com/gabrielcortez835-lab/Novellia-take-home-assessment/objects"
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
		`, constants.FHIRDataTableName)
	_, err := d.DB.Exec(schema)

	if err != nil {
		return err
	}

	schema = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		entryId TEXT,
		validationError TEXT
		);
		`, constants.ValidationErrorTableName)
	_, err = d.DB.Exec(schema)

	return err

}

func (d *Database) upgradeSchema() error {
	return nil
}

// InsertResource inserts a new record into resources table
func InsertResource(id, resourceType, subject, metadata string) error {

	db, err := Connect(constants.SqlDBName)
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
func GetRecordsById(id string) ([]gjson.Result, error) {
	db, err := Connect(constants.SqlDBName)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	query := fmt.Sprintf(`
    SELECT metadata
    FROM %s
    WHERE id = ?;
    `, constants.FHIRDataTableName)

	return GetQueryInternal(db, query, id)
}

func GetRecords(resourceType string, subject string) ([]gjson.Result, error) {
	db, err := Connect(constants.SqlDBName)

	if err != nil {
		return nil, err
	}
	defer db.Close()

	if resourceType != "" && subject != "" {
		return GetRecordsByResourceTypeAndSubjectInternal(db, resourceType, subject)
	} else if resourceType != "" {
		return GetRecordsByResourceTypeInternal(db, resourceType)
	} else if subject != "" {
		return GetRecordsBySubjectInternal(db, subject)
	} else {
		return GetRecordsInternal(db)
	}

}

func GetRecordsInternal(db *sql.DB) ([]gjson.Result, error) {
	query := fmt.Sprintf(`
        SELECT metadata
        FROM %s;
    `, constants.FHIRDataTableName)

	return GetQueryInternal(db, query)
}

func GetRecordsByResourceTypeInternal(db *sql.DB, resourceType string) ([]gjson.Result, error) {
	query := fmt.Sprintf(`
        SELECT metadata
        FROM %s
        WHERE resourceType = ?;
    `, constants.FHIRDataTableName)

	return GetQueryInternal(db, query, resourceType)
}

func GetRecordsByResourceTypeAndSubjectInternal(db *sql.DB, resourceType string, subject string) ([]gjson.Result, error) {
	query := fmt.Sprintf(`
        SELECT metadata
        FROM %s
        WHERE resourceType = ? and subject = ?;
    `, constants.FHIRDataTableName)

	return GetQueryInternal(db, query, resourceType, subject)
}

func GetRecordsBySubjectInternal(db *sql.DB, subject string) ([]gjson.Result, error) {
	query := fmt.Sprintf(`
        SELECT metadata
        FROM %s
        WHERE subject = ?;
    `, constants.FHIRDataTableName)

	return GetQueryInternal(db, query, subject)
}

func GetQueryInternal(db *sql.DB, query string, args ...interface{}) ([]gjson.Result, error) {
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

func GetRecordCountForResourceType(resourceType string) (int, error) {
	db, err := Connect(constants.SqlDBName)

	if err != nil {
		return 0, err
	}
	defer db.Close()

	return getRecordCountForResourceTypeInternal(db, resourceType)
}

func getRecordCountForResourceTypeInternal(db *sql.DB, resourceType string) (int, error) {
	var count int

	query := fmt.Sprintf(`
        SELECT COUNT(*) 
		FROM %s 
		WHERE resourceType = ?;
    `, constants.FHIRDataTableName)

	err := db.QueryRow(query, resourceType).Scan(&count)

	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetRecordCountForSubject() (int, error) {
	db, err := Connect(constants.SqlDBName)

	if err != nil {
		return 0, err
	}
	defer db.Close()

	return getRecordCountForSubjectInternal(db)
}

func getRecordCountForSubjectInternal(db *sql.DB) (int, error) {
	var count int

	query := fmt.Sprintf(`
        SELECT COUNT(DISTINCT subject) AS unique_count 
		FROM %s;
    `, constants.FHIRDataTableName)

	err := db.QueryRow(query).Scan(&count)

	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetRecordCountPerPatient() (map[string]int, error) {
	db, err := Connect(constants.SqlDBName)

	if err != nil {
		return nil, err
	}
	defer db.Close()

	return getRecordCountPerPatientInternal(db)
}

func getRecordCountPerPatientInternal(db *sql.DB) (map[string]int, error) {

	query := fmt.Sprintf(`
        SELECT subject, COUNT(*) AS record_count
		FROM %s
		GROUP BY subject;
    `, constants.FHIRDataTableName)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make(map[string]int)

	// Ensure any errors from iteration are caught
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for rows.Next() {
		var subject string
		var count int
		if err := rows.Scan(&subject, &count); err != nil {
			return nil, err
		}
		results[subject] = count
	}

	return results, nil
}

func InsertValidationError(entryID string, validationError string) error {
	db, err := Connect(constants.SqlDBName)
	if err != nil {
		return err
	}
	defer db.Close()

	return insertValidationErrorInternal(db, entryID, validationError)
}

func insertValidationErrorInternal(db *sql.DB, entryID string, validationError string) error {
	query := fmt.Sprintf(`
        INSERT INTO %s (entryId, validationError) 
		VALUES (?, ?)
    `, constants.ValidationErrorTableName)
	_, err := db.Exec(query, entryID, validationError)
	return err
}

func GetAllValidationErrors() ([]objects.ValidationError, error) {
	db, err := Connect(constants.SqlDBName)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	return getAllValidationErrorsInternal(db)
}

func getAllValidationErrorsInternal(db *sql.DB) ([]objects.ValidationError, error) {
	query := fmt.Sprintf(`
        SELECT entryId, validationError
        FROM %s
    `, constants.ValidationErrorTableName)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []objects.ValidationError
	for rows.Next() {
		var entryID string
		var rawValErr string

		if err := rows.Scan(&entryID, &rawValErr); err != nil {
			return nil, err
		}

		var valErrList []string
		if err := json.Unmarshal([]byte(rawValErr), &valErrList); err != nil {
			return nil, fmt.Errorf("failed to unmarshal ValidationError JSON for entry %s: %w", entryID, err)
		}

		results = append(results, objects.ValidationError{
			EntryID:         entryID,
			ValidationError: valErrList,
		})
	}
	return results, rows.Err()
}

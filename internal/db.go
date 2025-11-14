package runner

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type ScriptUsage struct {
	ID         int
	Directory  string
	ScriptName string
	Source     string // "make", "npm", "pnpm", "yarn"
	LastUsed   time.Time
	UseCount   int
	IsPinned   bool
}

type Database struct {
	db *sql.DB
}

func InitDatabase() (*Database, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "alex-runner")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	dbPath := filepath.Join(configDir, "alex-runner.sqlite.db")
	return InitDatabaseWithPath(dbPath)
}

func InitDatabaseWithPath(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := createSchema(db); err != nil {
		db.Close()
		return nil, err
	}

	return &Database{db: db}, nil
}

func createSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS script_usage (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		directory TEXT NOT NULL,
		script_name TEXT NOT NULL,
		source TEXT DEFAULT '',
		last_used TIMESTAMP NOT NULL,
		use_count INTEGER DEFAULT 1,
		is_pinned INTEGER DEFAULT 0,
		UNIQUE(directory, script_name, source)
	);

	CREATE INDEX IF NOT EXISTS idx_directory ON script_usage(directory);
	CREATE INDEX IF NOT EXISTS idx_frecency ON script_usage(directory, last_used DESC, use_count DESC);

	CREATE TABLE IF NOT EXISTS package_manager_cache (
		directory TEXT PRIMARY KEY,
		package_manager TEXT NOT NULL,
		detected_at TIMESTAMP NOT NULL
	);
	`

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	// Add columns to existing tables (for migration)
	_, _ = db.Exec(`ALTER TABLE script_usage ADD COLUMN is_pinned INTEGER DEFAULT 0`)
	_, _ = db.Exec(`ALTER TABLE script_usage ADD COLUMN source TEXT DEFAULT ''`)

	return nil
}

func (d *Database) RecordUsage(directory string, scriptName string, source string) error {
	query := `
	INSERT INTO script_usage (directory, script_name, source, last_used, use_count)
	VALUES (?, ?, ?, ?, 1)
	ON CONFLICT(directory, script_name, source)
	DO UPDATE SET
		last_used = ?,
		use_count = use_count + 1
	`

	now := time.Now()
	_, err := d.db.Exec(query, directory, scriptName, source, now, now)
	if err != nil {
		return fmt.Errorf("failed to record usage: %w", err)
	}

	return nil
}

func (d *Database) GetUsageStats(directory string) ([]ScriptUsage, error) {
	query := `
	SELECT id, directory, script_name, COALESCE(source, ''), last_used, use_count, COALESCE(is_pinned, 0)
	FROM script_usage
	WHERE directory = ?
	ORDER BY is_pinned DESC, last_used DESC, use_count DESC
	`

	rows, err := d.db.Query(query, directory)
	if err != nil {
		return nil, fmt.Errorf("failed to query usage stats: %w", err)
	}
	defer rows.Close()

	var usages []ScriptUsage
	for rows.Next() {
		var usage ScriptUsage
		var isPinnedInt int
		err := rows.Scan(&usage.ID, &usage.Directory, &usage.ScriptName, &usage.Source, &usage.LastUsed, &usage.UseCount, &isPinnedInt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		usage.IsPinned = isPinnedInt != 0
		usages = append(usages, usage)
	}

	return usages, nil
}

func (d *Database) ResetDirectory(directory string) error {
	query := `DELETE FROM script_usage WHERE directory = ?`
	_, err := d.db.Exec(query, directory)
	if err != nil {
		return fmt.Errorf("failed to reset directory: %w", err)
	}
	return nil
}

func (d *Database) ResetAll() error {
	query := `DELETE FROM script_usage`
	_, err := d.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to reset all: %w", err)
	}
	return nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

// GetCachedPackageManager retrieves the cached package manager for a directory
func (d *Database) GetCachedPackageManager(directory string) (string, error) {
	query := `SELECT package_manager FROM package_manager_cache WHERE directory = ?`
	var packageManager string
	err := d.db.QueryRow(query, directory).Scan(&packageManager)
	if err == sql.ErrNoRows {
		return "", nil // No cache entry found
	}
	if err != nil {
		return "", fmt.Errorf("failed to query package manager cache: %w", err)
	}
	return packageManager, nil
}

// SetCachedPackageManager stores the detected package manager for a directory
func (d *Database) SetCachedPackageManager(directory string, packageManager string) error {
	query := `
	INSERT INTO package_manager_cache (directory, package_manager, detected_at)
	VALUES (?, ?, ?)
	ON CONFLICT(directory)
	DO UPDATE SET
		package_manager = ?,
		detected_at = ?
	`
	now := time.Now()
	_, err := d.db.Exec(query, directory, packageManager, now, packageManager, now)
	if err != nil {
		return fmt.Errorf("failed to cache package manager: %w", err)
	}
	return nil
}

// PinScript pins a script for the given directory
func (d *Database) PinScript(directory string, scriptName string, source string) error {
	query := `
	INSERT INTO script_usage (directory, script_name, source, last_used, use_count, is_pinned)
	VALUES (?, ?, ?, ?, 0, 1)
	ON CONFLICT(directory, script_name, source)
	DO UPDATE SET is_pinned = 1
	`
	_, err := d.db.Exec(query, directory, scriptName, source, time.Now())
	if err != nil {
		return fmt.Errorf("failed to pin script: %w", err)
	}
	return nil
}

// UnpinScript unpins a script for the given directory
func (d *Database) UnpinScript(directory string, scriptName string, source string) error {
	query := `UPDATE script_usage SET is_pinned = 0 WHERE directory = ? AND script_name = ? AND source = ?`
	_, err := d.db.Exec(query, directory, scriptName, source)
	if err != nil {
		return fmt.Errorf("failed to unpin script: %w", err)
	}
	return nil
}

// TogglePin toggles the pin status of a script
func (d *Database) TogglePin(directory string, scriptName string, source string) (bool, error) {
	// Check current pin status
	var isPinned int
	query := `SELECT COALESCE(is_pinned, 0) FROM script_usage WHERE directory = ? AND script_name = ? AND source = ?`
	err := d.db.QueryRow(query, directory, scriptName, source).Scan(&isPinned)

	if err == sql.ErrNoRows {
		// Script doesn't exist in usage table yet, pin it
		if err := d.PinScript(directory, scriptName, source); err != nil {
			return false, err
		}
		return true, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to check pin status: %w", err)
	}

	// Toggle the pin status
	if isPinned != 0 {
		if err := d.UnpinScript(directory, scriptName, source); err != nil {
			return false, err
		}
		return false, nil
	} else {
		if err := d.PinScript(directory, scriptName, source); err != nil {
			return false, err
		}
		return true, nil
	}
}

// IsPinned checks if a script is pinned
func (d *Database) IsPinned(directory string, scriptName string, source string) (bool, error) {
	var isPinned int
	query := `SELECT COALESCE(is_pinned, 0) FROM script_usage WHERE directory = ? AND script_name = ? AND source = ?`
	err := d.db.QueryRow(query, directory, scriptName, source).Scan(&isPinned)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to check pin status: %w", err)
	}

	return isPinned != 0, nil
}

// FindScriptsByName finds all scripts with the given name across all sources
func (d *Database) FindScriptsByName(directory string, scriptName string) ([]ScriptUsage, error) {
	query := `
	SELECT id, directory, script_name, COALESCE(source, ''), last_used, use_count, COALESCE(is_pinned, 0)
	FROM script_usage
	WHERE directory = ? AND script_name = ?
	ORDER BY source
	`

	rows, err := d.db.Query(query, directory, scriptName)
	if err != nil {
		return nil, fmt.Errorf("failed to query scripts: %w", err)
	}
	defer rows.Close()

	var usages []ScriptUsage
	for rows.Next() {
		var usage ScriptUsage
		var isPinnedInt int
		err := rows.Scan(&usage.ID, &usage.Directory, &usage.ScriptName, &usage.Source, &usage.LastUsed, &usage.UseCount, &isPinnedInt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		usage.IsPinned = isPinnedInt != 0
		usages = append(usages, usage)
	}

	return usages, nil
}

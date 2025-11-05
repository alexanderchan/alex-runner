package runner

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupTestDB(t *testing.T) (*Database, string) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := InitDatabaseWithPath(dbPath)
	if err != nil {
		t.Fatalf("failed to initialize test database: %v", err)
	}

	return db, dbPath
}

func TestDatabaseInit(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	if db.db == nil {
		t.Fatal("database connection is nil")
	}
}

func TestRecordUsage(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	directory := "/test/project"
	scriptName := "dev"

	// Record first usage
	err := db.RecordUsage(directory, scriptName)
	if err != nil {
		t.Fatalf("failed to record usage: %v", err)
	}

	// Check usage stats
	stats, err := db.GetUsageStats(directory)
	if err != nil {
		t.Fatalf("failed to get usage stats: %v", err)
	}

	if len(stats) != 1 {
		t.Fatalf("expected 1 usage stat, got %d", len(stats))
	}

	if stats[0].ScriptName != scriptName {
		t.Errorf("expected script name %s, got %s", scriptName, stats[0].ScriptName)
	}

	if stats[0].UseCount != 1 {
		t.Errorf("expected use count 1, got %d", stats[0].UseCount)
	}
}

func TestRecordUsageIncrementsCount(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	directory := "/test/project"
	scriptName := "build"

	// Record usage multiple times
	for i := 0; i < 5; i++ {
		err := db.RecordUsage(directory, scriptName)
		if err != nil {
			t.Fatalf("failed to record usage: %v", err)
		}
		time.Sleep(10 * time.Millisecond) // Small delay to ensure timestamp changes
	}

	stats, err := db.GetUsageStats(directory)
	if err != nil {
		t.Fatalf("failed to get usage stats: %v", err)
	}

	if len(stats) != 1 {
		t.Fatalf("expected 1 usage stat, got %d", len(stats))
	}

	if stats[0].UseCount != 5 {
		t.Errorf("expected use count 5, got %d", stats[0].UseCount)
	}
}

func TestMultipleScriptsInDirectory(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	directory := "/test/project"
	scripts := []string{"dev", "build", "test", "lint"}

	// Record usage for multiple scripts
	for _, script := range scripts {
		err := db.RecordUsage(directory, script)
		if err != nil {
			t.Fatalf("failed to record usage for %s: %v", script, err)
		}
	}

	stats, err := db.GetUsageStats(directory)
	if err != nil {
		t.Fatalf("failed to get usage stats: %v", err)
	}

	if len(stats) != len(scripts) {
		t.Errorf("expected %d usage stats, got %d", len(scripts), len(stats))
	}
}

func TestMultipleDirectories(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	dir1 := "/test/project1"
	dir2 := "/test/project2"
	scriptName := "dev"

	// Record usage in both directories
	err := db.RecordUsage(dir1, scriptName)
	if err != nil {
		t.Fatalf("failed to record usage in dir1: %v", err)
	}

	err = db.RecordUsage(dir2, scriptName)
	if err != nil {
		t.Fatalf("failed to record usage in dir2: %v", err)
	}

	// Check stats for dir1
	stats1, err := db.GetUsageStats(dir1)
	if err != nil {
		t.Fatalf("failed to get usage stats for dir1: %v", err)
	}

	if len(stats1) != 1 {
		t.Errorf("expected 1 stat for dir1, got %d", len(stats1))
	}

	// Check stats for dir2
	stats2, err := db.GetUsageStats(dir2)
	if err != nil {
		t.Fatalf("failed to get usage stats for dir2: %v", err)
	}

	if len(stats2) != 1 {
		t.Errorf("expected 1 stat for dir2, got %d", len(stats2))
	}
}

func TestResetDirectory(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	directory := "/test/project"
	scripts := []string{"dev", "build", "test"}

	// Record usage for multiple scripts
	for _, script := range scripts {
		err := db.RecordUsage(directory, script)
		if err != nil {
			t.Fatalf("failed to record usage: %v", err)
		}
	}

	// Reset directory
	err := db.ResetDirectory(directory)
	if err != nil {
		t.Fatalf("failed to reset directory: %v", err)
	}

	// Check that stats are empty
	stats, err := db.GetUsageStats(directory)
	if err != nil {
		t.Fatalf("failed to get usage stats: %v", err)
	}

	if len(stats) != 0 {
		t.Errorf("expected 0 stats after reset, got %d", len(stats))
	}
}

func TestResetAll(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	directories := []string{"/test/project1", "/test/project2", "/test/project3"}

	// Record usage in multiple directories
	for _, dir := range directories {
		err := db.RecordUsage(dir, "dev")
		if err != nil {
			t.Fatalf("failed to record usage: %v", err)
		}
	}

	// Reset all
	err := db.ResetAll()
	if err != nil {
		t.Fatalf("failed to reset all: %v", err)
	}

	// Check that all directories are empty
	for _, dir := range directories {
		stats, err := db.GetUsageStats(dir)
		if err != nil {
			t.Fatalf("failed to get usage stats: %v", err)
		}

		if len(stats) != 0 {
			t.Errorf("expected 0 stats for %s after reset, got %d", dir, len(stats))
		}
	}
}

func TestLastUsedTimestamp(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	directory := "/test/project"
	scriptName := "dev"

	beforeTime := time.Now()
	time.Sleep(10 * time.Millisecond)

	err := db.RecordUsage(directory, scriptName)
	if err != nil {
		t.Fatalf("failed to record usage: %v", err)
	}

	time.Sleep(10 * time.Millisecond)
	afterTime := time.Now()

	stats, err := db.GetUsageStats(directory)
	if err != nil {
		t.Fatalf("failed to get usage stats: %v", err)
	}

	if len(stats) != 1 {
		t.Fatalf("expected 1 stat, got %d", len(stats))
	}

	lastUsed := stats[0].LastUsed
	if lastUsed.Before(beforeTime) || lastUsed.After(afterTime) {
		t.Errorf("last_used timestamp %v is not between %v and %v", lastUsed, beforeTime, afterTime)
	}
}

func TestGetUsageStatsEmptyDirectory(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	directory := "/test/empty-project"

	stats, err := db.GetUsageStats(directory)
	if err != nil {
		t.Fatalf("failed to get usage stats: %v", err)
	}

	if len(stats) != 0 {
		t.Errorf("expected 0 stats for empty directory, got %d", len(stats))
	}
}

func TestDatabasePersistence(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "persist.db")

	// Create database and add data
	db1, err := InitDatabaseWithPath(dbPath)
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}

	directory := "/test/project"
	scriptName := "dev"

	err = db1.RecordUsage(directory, scriptName)
	if err != nil {
		t.Fatalf("failed to record usage: %v", err)
	}

	db1.Close()

	// Reopen database and check data persists
	db2, err := InitDatabaseWithPath(dbPath)
	if err != nil {
		t.Fatalf("failed to reopen database: %v", err)
	}
	defer db2.Close()

	stats, err := db2.GetUsageStats(directory)
	if err != nil {
		t.Fatalf("failed to get usage stats: %v", err)
	}

	if len(stats) != 1 {
		t.Fatalf("expected 1 stat after reopening, got %d", len(stats))
	}

	if stats[0].ScriptName != scriptName {
		t.Errorf("expected script name %s, got %s", scriptName, stats[0].ScriptName)
	}
}

func TestInitDatabaseCreatesFile(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "new.db")

	// Ensure file doesn't exist
	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		t.Fatal("database file should not exist before initialization")
	}

	db, err := InitDatabaseWithPath(dbPath)
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	// Check file was created
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatal("database file should exist after initialization")
	}
}

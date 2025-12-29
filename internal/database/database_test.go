// Package database provides database connection and testing utilities.
package database

import (
	"os"
	"path/filepath"
	"testing"
)

// TestBuildPath tests the buildPath function
func TestBuildPath(t *testing.T) {
	tests := []struct {
		name       string
		params     []string
		wantErr    bool
		errType    error
	}{
		{
			name:    "Valid two params",
			params:  []string{"testdb", "shard1"},
			wantErr: false,
		},
		{
			name:    "Valid three params",
			params:  []string{"testdb", "shard1", "data.db"},
			wantErr: false,
		},
		{
			name:    "Single param should fail",
			params:  []string{"testdb"},
			wantErr: true,
			errType: ErrInvalidPath,
		},
		{
			name:    "Empty params should fail",
			params:  []string{},
			wantErr: true,
			errType: ErrInvalidPath,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := buildPath(tt.params...)

			if tt.wantErr {
				if err == nil {
					t.Errorf("buildPath() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("buildPath() unexpected error: %v", err)
				return
			}

			if path == "" {
				t.Errorf("buildPath() returned empty path")
			}

			// Cleanup
			if len(tt.params) >= 2 {
				os.RemoveAll(tt.params[0])
			}
		})
	}
}

// TestCreateLevelDBDatabase tests the LevelDB database creation
func TestCreateLevelDBDatabase(t *testing.T) {
	testDir := filepath.Join(os.TempDir(), "dm-backend-test")
	defer os.RemoveAll(testDir)

	tests := []struct {
		name    string
		params  []string
		wantErr bool
	}{
		{
			name:    "Valid database creation",
			params:  []string{testDir, "test_shard", "test.db"},
			wantErr: false,
		},
		{
			name:    "Invalid params should fail",
			params:  []string{testDir},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := CreateLevelDBDatabase(tt.params...)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateLevelDBDatabase() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("CreateLevelDBDatabase() unexpected error: %v", err)
				return
			}

			if db == nil {
				t.Errorf("CreateLevelDBDatabase() returned nil database")
				return
			}

			// Test basic operations
			testKey := []byte("test_key")
			testValue := []byte("test_value")

			err = db.Put(testKey, testValue, nil)
			if err != nil {
				t.Errorf("Failed to put value: %v", err)
			}

			value, err := db.Get(testKey, nil)
			if err != nil {
				t.Errorf("Failed to get value: %v", err)
			}

			if string(value) != string(testValue) {
				t.Errorf("Expected value %s, got %s", testValue, value)
			}

			db.Close()
		})
	}
}

// TestCreateSQLiteDatabase tests the SQLite database creation
func TestCreateSQLiteDatabase(t *testing.T) {
	testDir := filepath.Join(os.TempDir(), "dm-backend-test-sqlite")
	defer os.RemoveAll(testDir)

	tests := []struct {
		name    string
		params  []string
		wantErr bool
	}{
		{
			name:    "Valid database creation",
			params:  []string{testDir, "test_shard", "test.sqlite"},
			wantErr: false,
		},
		{
			name:    "Invalid params should fail",
			params:  []string{testDir},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := CreateSQLiteDatabase(tt.params...)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateSQLiteDatabase() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("CreateSQLiteDatabase() unexpected error: %v", err)
				return
			}

			if db == nil {
				t.Errorf("CreateSQLiteDatabase() returned nil database")
				return
			}

			// Test basic operations
			_, err = db.Exec(`CREATE TABLE IF NOT EXISTS test (id INTEGER PRIMARY KEY, name TEXT)`)
			if err != nil {
				t.Errorf("Failed to create table: %v", err)
			}

			_, err = db.Exec(`INSERT INTO test (name) VALUES (?)`, "test_name")
			if err != nil {
				t.Errorf("Failed to insert: %v", err)
			}

			var name string
			err = db.QueryRow(`SELECT name FROM test WHERE id = 1`).Scan(&name)
			if err != nil {
				t.Errorf("Failed to query: %v", err)
			}

			if name != "test_name" {
				t.Errorf("Expected name 'test_name', got '%s'", name)
			}

			db.Close()
		})
	}
}

// TestCleanPath tests the CleanPath utility function
func TestCleanPath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/path/to/file/", "path/to/file"},
		{"path/to/file", "path/to/file"},
		{"/path/", "path"},
		{"//path//to//", "path/to"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := CleanPath(tt.input)
			if result != tt.expected {
				t.Errorf("CleanPath(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

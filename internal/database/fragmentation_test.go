// Package database provides fragmentation management testing.
package database

import (
	"os"
	"path/filepath"
	"testing"
)

// setupTestEnvironment creates a temporary test directory
func setupTestEnvironment(t *testing.T) (cleanup func()) {
	testDir := filepath.Join(os.TempDir(), "dm-backend-fragmentation-test")
	os.MkdirAll(testDir, 0755)

	// Override the default paths by managing our own test paths
	return func() {
		os.RemoveAll(testDir)
		os.RemoveAll("Database")
	}
}

// TestFragmentationAdd tests adding fragmentation entries
func TestFragmentationAdd(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	tests := []struct {
		name    string
		shard   int64
		keys    []string
		wantErr bool
	}{
		{
			name:    "Add single key",
			shard:   1,
			keys:    []string{"user_001"},
			wantErr: false,
		},
		{
			name:    "Add multiple keys",
			shard:   2,
			keys:    []string{"user_002", "user_003", "user_004"},
			wantErr: false,
		},
		{
			name:    "Empty keys should fail",
			shard:   3,
			keys:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := FragmentationAdd(tt.shard, tt.keys...)

			if tt.wantErr {
				if err == nil {
					t.Errorf("FragmentationAdd() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("FragmentationAdd() unexpected error: %v", err)
				return
			}

			// Verify keys were added
			for _, key := range tt.keys {
				shard, err := FragmentationGet(key)
				if err != nil {
					t.Errorf("Failed to get key %s: %v", key, err)
					continue
				}
				if shard != tt.shard {
					t.Errorf("Expected shard %d for key %s, got %d", tt.shard, key, shard)
				}
			}
		})
	}
}

// TestFragmentationGet tests retrieving fragmentation entries
func TestFragmentationGet(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Setup: Add a test entry
	testKey := "test_get_key"
	testShard := int64(42)
	err := FragmentationAdd(testShard, testKey)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	tests := []struct {
		name      string
		rowID     string
		wantShard int64
		wantErr   bool
	}{
		{
			name:      "Get existing key",
			rowID:     testKey,
			wantShard: testShard,
			wantErr:   false,
		},
		{
			name:      "Get non-existing key",
			rowID:     "non_existing_key",
			wantShard: -1,
			wantErr:   true,
		},
		{
			name:      "Empty key should fail",
			rowID:     "",
			wantShard: -1,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shard, err := FragmentationGet(tt.rowID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("FragmentationGet() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("FragmentationGet() unexpected error: %v", err)
				return
			}

			if shard != tt.wantShard {
				t.Errorf("FragmentationGet() = %d, want %d", shard, tt.wantShard)
			}
		})
	}
}

// TestFragmentationRemove tests removing fragmentation entries
func TestFragmentationRemove(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Setup: Add test entries
	testKey := "test_remove_key"
	err := FragmentationAdd(1, testKey)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{
			name:    "Remove existing key",
			key:     testKey,
			wantErr: false,
		},
		{
			name:    "Remove non-existing key should fail",
			key:     "non_existing_key",
			wantErr: true,
		},
		{
			name:    "Empty key should fail",
			key:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := FragmentationRemove(tt.key)

			if tt.wantErr {
				if err == nil {
					t.Errorf("FragmentationRemove() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("FragmentationRemove() unexpected error: %v", err)
				return
			}

			// Verify key was removed
			_, err = FragmentationGet(tt.key)
			if err == nil {
				t.Errorf("Key %s should have been removed", tt.key)
			}
		})
	}
}

// TestFragmentationExists tests checking if keys exist
func TestFragmentationExists(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Setup: Add a test entry
	testKey := "test_exists_key"
	err := FragmentationAdd(1, testKey)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	tests := []struct {
		name       string
		key        string
		wantExists bool
		wantErr    bool
	}{
		{
			name:       "Check existing key",
			key:        testKey,
			wantExists: true,
			wantErr:    false,
		},
		{
			name:       "Check non-existing key",
			key:        "non_existing_key",
			wantExists: false,
			wantErr:    false,
		},
		{
			name:       "Empty key should fail",
			key:        "",
			wantExists: false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := FragmentationExists(tt.key)

			if tt.wantErr {
				if err == nil {
					t.Errorf("FragmentationExists() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("FragmentationExists() unexpected error: %v", err)
				return
			}

			if exists != tt.wantExists {
				t.Errorf("FragmentationExists() = %v, want %v", exists, tt.wantExists)
			}
		})
	}
}

// TestFragmentationUpdate tests updating fragmentation entries
func TestFragmentationUpdate(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Setup: Add a test entry
	testKey := "test_update_key"
	originalShard := int64(1)
	err := FragmentationAdd(originalShard, testKey)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	tests := []struct {
		name     string
		key      string
		newShard int64
		wantErr  bool
	}{
		{
			name:     "Update existing key",
			key:      testKey,
			newShard: 99,
			wantErr:  false,
		},
		{
			name:     "Update non-existing key should fail",
			key:      "non_existing_key",
			newShard: 99,
			wantErr:  true,
		},
		{
			name:     "Empty key should fail",
			key:      "",
			newShard: 99,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := FragmentationUpdate(tt.key, tt.newShard)

			if tt.wantErr {
				if err == nil {
					t.Errorf("FragmentationUpdate() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("FragmentationUpdate() unexpected error: %v", err)
				return
			}

			// Verify the update
			shard, err := FragmentationGet(tt.key)
			if err != nil {
				t.Errorf("Failed to get updated key: %v", err)
				return
			}

			if shard != tt.newShard {
				t.Errorf("Expected shard %d, got %d", tt.newShard, shard)
			}
		})
	}
}

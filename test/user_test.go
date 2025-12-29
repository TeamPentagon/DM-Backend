// Package test provides integration tests for the DM-Backend application.
// These tests verify the correct operation of user management functionality.
package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/TeamPentagon/DM-Backend/internal/model"
)

// setupTestEnvironment creates a clean test environment
func setupTestEnvironment(t *testing.T) func() {
	// Store the current directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Create a temporary test directory
	testDir := filepath.Join(os.TempDir(), "dm-backend-user-test")
	os.MkdirAll(testDir, 0755)
	os.Chdir(testDir)

	return func() {
		os.Chdir(originalDir)
		os.RemoveAll(testDir)
	}
}

// createTestUser creates a user with test data
func createTestUser() *model.User {
	return &model.User{
		UserName:      "testUser",
		PersonId:      "testPersonId",
		Profile:       "testProfile",
		Password:      "testPassword",
		Email:         "test@example.com",
		ProfilePicUrl: "http://test.com/test.jpg",
		AccountTime:   1234567890,
		BirthDate:     1233567890,
		Gender:        "testGender",
		LastEdit:      1234567890,
		PhoneNumber:   "1234567890",
	}
}

// TestSaveUserData tests creating and saving a new user
func TestSaveUserData(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	user := createTestUser()

	err := user.SaveUserData()
	if err != nil {
		t.Errorf("SaveUserData() returned an error: %v", err)
	}

	// Verify data wasn't modified
	if user.UserName != "testUser" {
		t.Errorf("Expected UserName to be 'testUser', but got '%s'", user.UserName)
	}
	if user.PersonId != "testPersonId" {
		t.Errorf("Expected PersonId to be 'testPersonId', but got '%s'", user.PersonId)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected Email to be 'test@example.com', but got '%s'", user.Email)
	}
}

// TestSaveUserDataWithEmptyPersonId tests validation
func TestSaveUserDataWithEmptyPersonId(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	user := &model.User{
		UserName: "testUser",
		PersonId: "", // Empty - should fail
	}

	err := user.SaveUserData()
	if err == nil {
		t.Error("SaveUserData() should return error for empty PersonId")
	}
}

// TestGetUserData tests retrieving user data
func TestGetUserData(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// First save the user
	originalUser := createTestUser()
	err := originalUser.SaveUserData()
	if err != nil {
		t.Fatalf("Setup failed - SaveUserData() returned an error: %v", err)
	}

	// Now retrieve it
	retrievedUser := &model.User{}
	err = retrievedUser.GetUserData("testPersonId")

	if err != nil {
		t.Errorf("GetUserData() returned an error: %v", err)
	}

	// Verify retrieved data matches
	if retrievedUser.UserName != originalUser.UserName {
		t.Errorf("Expected UserName '%s', got '%s'", originalUser.UserName, retrievedUser.UserName)
	}
	if retrievedUser.PersonId != originalUser.PersonId {
		t.Errorf("Expected PersonId '%s', got '%s'", originalUser.PersonId, retrievedUser.PersonId)
	}
	if retrievedUser.Email != originalUser.Email {
		t.Errorf("Expected Email '%s', got '%s'", originalUser.Email, retrievedUser.Email)
	}
	if retrievedUser.PhoneNumber != originalUser.PhoneNumber {
		t.Errorf("Expected PhoneNumber '%s', got '%s'", originalUser.PhoneNumber, retrievedUser.PhoneNumber)
	}
}

// TestGetUserDataNotFound tests retrieving non-existent user
func TestGetUserDataNotFound(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	user := &model.User{}
	err := user.GetUserData("nonExistentUser")

	if err == nil {
		t.Error("GetUserData() should return error for non-existent user")
	}
}

// TestGetUserDataEmptyUsername tests validation
func TestGetUserDataEmptyUsername(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	user := &model.User{}
	err := user.GetUserData("")

	if err == nil {
		t.Error("GetUserData() should return error for empty username")
	}
}

// TestUpdateUserData tests updating user data
func TestUpdateUserData(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// First save the user
	user := createTestUser()
	err := user.SaveUserData()
	if err != nil {
		t.Fatalf("Setup failed - SaveUserData() returned an error: %v", err)
	}

	// Update the user
	user.Profile = "updatedProfile"
	user.Email = "updated@example.com"
	err = user.UpdateUserData("testPersonId")

	if err != nil {
		t.Errorf("UpdateUserData() returned an error: %v", err)
	}

	// Verify the update
	retrievedUser := &model.User{}
	err = retrievedUser.GetUserData("testPersonId")
	if err != nil {
		t.Fatalf("GetUserData() failed: %v", err)
	}

	if retrievedUser.Profile != "updatedProfile" {
		t.Errorf("Expected Profile 'updatedProfile', got '%s'", retrievedUser.Profile)
	}
	if retrievedUser.Email != "updated@example.com" {
		t.Errorf("Expected Email 'updated@example.com', got '%s'", retrievedUser.Email)
	}
}

// TestUpdateUserDataNotFound tests updating non-existent user
func TestUpdateUserDataNotFound(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	user := createTestUser()
	err := user.UpdateUserData("nonExistentUser")

	if err == nil {
		t.Error("UpdateUserData() should return error for non-existent user")
	}
}

// TestDeleteUserData tests deleting user data
func TestDeleteUserData(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// First save the user
	user := createTestUser()
	err := user.SaveUserData()
	if err != nil {
		t.Fatalf("Setup failed - SaveUserData() returned an error: %v", err)
	}

	// Delete the user
	err = user.DeleteUserData("testPersonId")
	if err != nil {
		t.Errorf("DeleteUserData() returned an error: %v", err)
	}

	// Verify the user is deleted
	deletedUser := &model.User{}
	err = deletedUser.GetUserData("testPersonId")
	if err == nil {
		t.Error("GetUserData() should return error for deleted user")
	}
}

// TestDeleteUserDataNotFound tests deleting non-existent user
func TestDeleteUserDataNotFound(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	user := &model.User{}
	err := user.DeleteUserData("nonExistentUser")

	if err == nil {
		t.Error("DeleteUserData() should return error for non-existent user")
	}
}

// TestUserExists tests the UserExists helper function
func TestUserExists(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// First save a user
	user := createTestUser()
	err := user.SaveUserData()
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// Test existing user
	exists, err := model.UserExists("testPersonId")
	if err != nil {
		t.Errorf("UserExists() returned error: %v", err)
	}
	if !exists {
		t.Error("UserExists() should return true for existing user")
	}

	// Test non-existing user
	exists, err = model.UserExists("nonExistentUser")
	if err != nil {
		t.Errorf("UserExists() returned unexpected error: %v", err)
	}
	if exists {
		t.Error("UserExists() should return false for non-existing user")
	}
}

// TestUserCRUDLifecycle tests the complete CRUD lifecycle
func TestUserCRUDLifecycle(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// CREATE
	user := createTestUser()
	err := user.SaveUserData()
	if err != nil {
		t.Fatalf("CREATE failed: %v", err)
	}

	// READ
	readUser := &model.User{}
	err = readUser.GetUserData("testPersonId")
	if err != nil {
		t.Fatalf("READ failed: %v", err)
	}
	if readUser.UserName != user.UserName {
		t.Errorf("READ verification failed: expected '%s', got '%s'", user.UserName, readUser.UserName)
	}

	// UPDATE
	readUser.Profile = "New Profile"
	readUser.Email = "new@example.com"
	err = readUser.UpdateUserData("testPersonId")
	if err != nil {
		t.Fatalf("UPDATE failed: %v", err)
	}

	// Verify UPDATE
	updatedUser := &model.User{}
	err = updatedUser.GetUserData("testPersonId")
	if err != nil {
		t.Fatalf("UPDATE verification failed: %v", err)
	}
	if updatedUser.Profile != "New Profile" {
		t.Errorf("UPDATE verification failed: expected 'New Profile', got '%s'", updatedUser.Profile)
	}

	// DELETE
	err = updatedUser.DeleteUserData("testPersonId")
	if err != nil {
		t.Fatalf("DELETE failed: %v", err)
	}

	// Verify DELETE
	exists, _ := model.UserExists("testPersonId")
	if exists {
		t.Error("DELETE verification failed: user still exists")
	}
}

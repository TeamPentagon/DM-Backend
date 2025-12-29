package test

import (
	"testing"

	"github.com/TeamPentagon/DM-Backend/internal/model"
)

func TestSaveUserData(t *testing.T) {
	// Initialize test data
	user := &model.User{
		UserName:      "testUser",
		PersonId:      "testPersonId",
		Profile:       "testProfile",
		Password:      "testPassword",
		Email:         "testEmail",
		ProfilePicUrl: "http://test.com/test.jpg",
		AccountTime:   1234567890, //unix time
		BirthDate:     1233567890,
		Gender:        "testGender",
		LastEdit:      1234567890,
		PhoneNumber:   "1234567890",
	}

	// Call the function being tested
	err := user.SaveUserData()

	// Check if the function returned an error
	if err != nil {
		t.Errorf("SaveUserData() returned an error: %v", err)
	}

	// Add additional assertions if needed
	if user.UserName != "testUser" {
		t.Errorf("Expected UserName to be 'testUser', but got '%s'", user.UserName)
	}
	if user.PersonId != "testPersonId" {
		t.Errorf("Expected PersonId to be 'testPersonId', but got '%s'", user.PersonId)
	}
	if user.Profile != "testProfile" {
		t.Errorf("Expected Profile to be 'testProfile', but got '%s'", user.Profile)
	}
	if user.Password != "testPassword" {
		t.Errorf("Expected Password to be 'testPassword', but got '%s'", user.Password)
	}
	if user.Email != "testEmail" {
		t.Errorf("Expected Email to be 'testEmail', but got '%s'", user.Email)
	}
	if user.ProfilePicUrl != "http://test.com/test.jpg" {
		t.Errorf("Expected ProfilePicUrl to be 'http://test.com/test.jpg', but got '%s'", user.ProfilePicUrl)
	}
	if user.AccountTime != 1234567890 {
		t.Errorf("Expected AccountTime to be 1234567890, but got '%d'", user.AccountTime)
	}
	if user.BirthDate != 1233567890 {
		t.Errorf("Expected BirthDate to be 1233567890, but got '%d'", user.BirthDate)
	}
	if user.Gender != "testGender" {
		t.Errorf("Expected Gender to be 'testGender', but got '%s'", user.Gender)
	}
	if user.LastEdit != 1234567890 {
		t.Errorf("Expected LastEdit to be 1234567890, but got '%d'", user.LastEdit)
	}
	if user.PhoneNumber != "1234567890" {
		t.Errorf("Expected PhoneNumber to be '1234567890', but got '%s'", user.PhoneNumber)
	}
}

func TestGetUserData(t *testing.T) {
	// Initialize test data
	user := &model.User{}
	// Call the function being tested
	err := user.GetUserData("testUser")

	// Check if the function returned an error
	if err != nil {
		t.Errorf("GetUserData() returned an error: %v", err)
	}

	// Add additional assertions if needed
	// Add additional assertions if needed
	if user.UserName != "testUser" {
		t.Errorf("Expected UserName to be 'testUser', but got '%s'", user.UserName)
	}
	if user.PersonId != "testPersonId" {
		t.Errorf("Expected PersonId to be 'testPersonId', but got '%s'", user.PersonId)
	}
	if user.Profile != "testProfile" {
		t.Errorf("Expected Profile to be 'testProfile', but got '%s'", user.Profile)
	}
	if user.Password != "testPassword" {
		t.Errorf("Expected Password to be 'testPassword', but got '%s'", user.Password)
	}
	if user.Email != "testEmail" {
		t.Errorf("Expected Email to be 'testEmail', but got '%s'", user.Email)
	}
	if user.ProfilePicUrl != "http://test.com/test.jpg" {
		t.Errorf("Expected ProfilePicUrl to be 'http://test.com/test.jpg', but got '%s'", user.ProfilePicUrl)
	}
	if user.AccountTime != 1234567890 {
		t.Errorf("Expected AccountTime to be 1234567890, but got '%d'", user.AccountTime)
	}
	if user.BirthDate != 1233567890 {
		t.Errorf("Expected BirthDate to be 1233567890, but got '%d'", user.BirthDate)
	}
	if user.Gender != "testGender" {
		t.Errorf("Expected Gender to be 'testGender', but got '%s'", user.Gender)
	}
	if user.LastEdit != 1234567890 {
		t.Errorf("Expected LastEdit to be 1234567890, but got '%d'", user.LastEdit)
	}
	if user.PhoneNumber != "1234567890" {
		t.Errorf("Expected PhoneNumber to be '1234567890', but got '%s'", user.PhoneNumber)
	}
}

func TestUpdatedUserData(t *testing.T) {
	// Initialize test data
	userName := "testUser"

	// Call the function being tested
	user := &model.User{}
	err := user.GetUserData(userName)
	if err != nil {
		t.Errorf("GetUserData() returned an error: %v", err)
	}
	user.Profile = "testProfile"
	err = user.UpdatedUserData(userName)

	// Check if the function returned an error
	if err != nil {
		t.Errorf("UpdatedUserData() returned an error: %v", err)
	}

	// Add additional assertions if needed
}

func TestDeleteUserData(t *testing.T) {

	// Initialize test data
	userName := "testUser"

	// Call the function being tested
	user := &model.User{}
	err := user.GetUserData(userName)
	if err != nil {
		t.Errorf("GetUserData() returned an error: %v", err)
	}
	err = user.DeleteUserData(userName)

	// Check if the function returned an error
	if err != nil {
		t.Errorf("DeleteUserData() returned an error: %v", err)
	}

	// Add additional assertions if needed
}

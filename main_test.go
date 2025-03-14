package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var originalHome string
var tempHome string

// TestMain sets up a temporary HOME so that our config file is isolated during tests.
func TestMain(m *testing.M) {
	originalHome = os.Getenv("HOME")
	var err error
	tempHome, err = ioutil.TempDir("", "git-user-test")
	if err != nil {
		panic(err)
	}
	os.Setenv("HOME", tempHome)
	
	code := m.Run()
	
	os.Setenv("HOME", originalHome)
	os.RemoveAll(tempHome)
	os.Exit(code)
}

func TestGetConfigFilePath(t *testing.T) {
	path, err := getConfigFilePath()
	if err != nil {
		t.Fatalf("getConfigFilePath error: %v", err)
	}
	expectedSuffix := filepath.Join(".config", "git-user", "profiles.json")
	if !strings.HasSuffix(path, expectedSuffix) {
		t.Errorf("Expected path to end with %s, got %s", expectedSuffix, path)
	}
}

func TestLoadProfilesEmpty(t *testing.T) {
	path, err := getConfigFilePath()
	if err != nil {
		t.Fatalf("getConfigFilePath error: %v", err)
	}
	// Ensure the file is removed so that loadProfiles returns an empty slice.
	os.Remove(path)
	profiles, err := loadProfiles()
	if err != nil {
		t.Fatalf("loadProfiles error: %v", err)
	}
	if len(profiles) != 0 {
		t.Errorf("Expected 0 profiles, got %d", len(profiles))
	}
}

func TestAddProfile(t *testing.T) {
	path, err := getConfigFilePath()
	if err != nil {
		t.Fatalf("getConfigFilePath error: %v", err)
	}
	os.Remove(path)
	err = addProfile("TestUser", "test@example.com")
	if err != nil {
		t.Fatalf("addProfile error: %v", err)
	}
	profiles, err := loadProfiles()
	if err != nil {
		t.Fatalf("loadProfiles error: %v", err)
	}
	if len(profiles) != 1 {
		t.Errorf("Expected 1 profile, got %d", len(profiles))
	}
	if profiles[0].Name != "TestUser" || profiles[0].Email != "test@example.com" {
		t.Errorf("Profile data mismatch: got %+v", profiles[0])
	}
}

func TestDeleteProfile(t *testing.T) {
	path, err := getConfigFilePath()
	if err != nil {
		t.Fatalf("getConfigFilePath error: %v", err)
	}
	os.Remove(path)
	// Add two profiles.
	err = addProfile("User1", "user1@example.com")
	if err != nil {
		t.Fatalf("addProfile error: %v", err)
	}
	err = addProfile("User2", "user2@example.com")
	if err != nil {
		t.Fatalf("addProfile error: %v", err)
	}
	// Delete one profile.
	err = deleteProfile("User1")
	if err != nil {
		t.Fatalf("deleteProfile error: %v", err)
	}
	profiles, err := loadProfiles()
	if err != nil {
		t.Fatalf("loadProfiles error: %v", err)
	}
	if len(profiles) != 1 {
		t.Errorf("Expected 1 profile after deletion, got %d", len(profiles))
	}
	if profiles[0].Name != "User2" {
		t.Errorf("Expected remaining profile to be User2, got %s", profiles[0].Name)
	}
}

func TestSaveAndLoadProfiles(t *testing.T) {
	profilesToSave := []GitProfile{
		{Name: "UserA", Email: "a@example.com"},
		{Name: "UserB", Email: "b@example.com"},
	}
	err := saveProfiles(profilesToSave)
	if err != nil {
		t.Fatalf("saveProfiles error: %v", err)
	}
	loaded, err := loadProfiles()
	if err != nil {
		t.Fatalf("loadProfiles error: %v", err)
	}
	if len(loaded) != 2 {
		t.Errorf("Expected 2 profiles, got %d", len(loaded))
	}
}

func TestListProfilesOutput(t *testing.T) {
	// Set up test profiles.
	profiles := []GitProfile{
		{Name: "UserX", Email: "x@example.com"},
		{Name: "UserY", Email: "y@example.com"},
	}
	err := saveProfiles(profiles)
	if err != nil {
		t.Fatalf("saveProfiles error: %v", err)
	}
	
	// Capture stdout.
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	err = listProfiles()
	if err != nil {
		t.Fatalf("listProfiles error: %v", err)
	}
	
	w.Close()
	var buf strings.Builder
	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Fatalf("Failed to read stdout: %v", err)
	}
	output := buf.String()
	os.Stdout = oldStdout
	
	if !strings.Contains(output, "UserX") || !strings.Contains(output, "UserY") {
		t.Errorf("Output does not contain expected profiles: %s", output)
	}
}

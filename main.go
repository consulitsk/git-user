package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// GitProfile represents a git user's profile.
type GitProfile struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Path for storing profiles.
const configDir = ".config/git-user"
const configFile = "profiles.json"

// getConfigFilePath returns the full path to the profiles file,
// creating the directory if it doesn't exist.
func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, configDir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
	}
	return filepath.Join(dir, configFile), nil
}

// loadProfiles loads profiles from the file. If the file does not exist, it returns an empty list.
func loadProfiles() ([]GitProfile, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []GitProfile{}, nil
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var profiles []GitProfile
	err = json.Unmarshal(data, &profiles)
	if err != nil {
		return nil, err
	}
	return profiles, nil
}

// saveProfiles saves profiles to the file.
func saveProfiles(profiles []GitProfile) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(profiles, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}

// addProfile adds a new profile. If a profile with the same name exists, it returns an error.
func addProfile(name, email string) error {
	profiles, err := loadProfiles()
	if err != nil {
		return err
	}
	for _, p := range profiles {
		if p.Name == name {
			return fmt.Errorf("profile with name '%s' already exists", name)
		}
	}
	profiles = append(profiles, GitProfile{Name: name, Email: email})
	return saveProfiles(profiles)
}

// deleteProfile removes a profile by name.
func deleteProfile(name string) error {
	profiles, err := loadProfiles()
	if err != nil {
		return err
	}
	index := -1
	for i, p := range profiles {
		if p.Name == name {
			index = i
			break
		}
	}
	if index == -1 {
		return fmt.Errorf("profile with name '%s' not found", name)
	}
	profiles = append(profiles[:index], profiles[index+1:]...)
	return saveProfiles(profiles)
}

// getCurrentGitConfig returns the current git configuration value for the given key.
func getCurrentGitConfig(key string) (string, error) {
	out, err := exec.Command("git", "config", "--global", key).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// listProfiles prints the list of stored profiles, marking the currently set git profile.
// If the current git profile is not among the stored profiles, a message is printed.
func listProfiles() error {
	profiles, err := loadProfiles()
	if err != nil {
		return err
	}
	currentName, err := getCurrentGitConfig("user.name")
	if err != nil {
		currentName = ""
	}
	currentEmail, err := getCurrentGitConfig("user.email")
	if err != nil {
		currentEmail = ""
	}

	foundCurrent := false
	fmt.Println("Stored profiles:")
	for i, p := range profiles {
		marker := " "
		if p.Name == currentName && p.Email == currentEmail {
			marker = "*"
			foundCurrent = true
		}
		fmt.Printf("[%d] %s - %s %s\n", i, p.Name, p.Email, marker)
	}
	if currentName != "" && currentEmail != "" && !foundCurrent {
		fmt.Printf("Current active git profile: %s - %s\n", currentName, currentEmail)
		fmt.Println("Note: The current git profile is not imported in the profiles list.")
		fmt.Println("You can import it using the 'import' subcommand.")
	}
	return nil
}

// setGitConfig sets the global git configuration for user.name and user.email.
func setGitConfig(name, email string) error {
	if err := exec.Command("git", "config", "--global", "user.name", name).Run(); err != nil {
		return err
	}
	return exec.Command("git", "config", "--global", "user.email", email).Run()
}

// switchProfile interactively switches the git profile by letting the user choose from a list.
func switchProfile() error {
	profiles, err := loadProfiles()
	if err != nil {
		return err
	}
	if len(profiles) == 0 {
		return fmt.Errorf("no profiles found, add a profile using 'git-user add'")
	}
	fmt.Println("Select a profile:")
	for i, p := range profiles {
		fmt.Printf("[%d] %s - %s\n", i, p.Name, p.Email)
	}
	fmt.Print("Enter profile number: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	input = strings.TrimSpace(input)
	index, err := strconv.Atoi(input)
	if err != nil || index < 0 || index >= len(profiles) {
		return fmt.Errorf("invalid selection")
	}
	selected := profiles[index]
	if err := setGitConfig(selected.Name, selected.Email); err != nil {
		return err
	}
	fmt.Printf("Profile switched to: %s - %s\n", selected.Name, selected.Email)
	return nil
}

// importProfile imports the current global git profile into the stored profiles.
func importProfile() error {
	name, err := getCurrentGitConfig("user.name")
	if err != nil || name == "" {
		return fmt.Errorf("no current git user name found")
	}
	email, err := getCurrentGitConfig("user.email")
	if err != nil || email == "" {
		return fmt.Errorf("no current git user email found")
	}
	profiles, err := loadProfiles()
	if err != nil {
		return err
	}
	for _, p := range profiles {
		if p.Name == name && p.Email == email {
			return fmt.Errorf("current git profile already exists in profiles")
		}
	}
	profiles = append(profiles, GitProfile{Name: name, Email: email})
	if err := saveProfiles(profiles); err != nil {
		return err
	}
	fmt.Printf("Current git profile imported: %s - %s\n", name, email)
	return nil
}

func main() {
	// On startup, display current git profile information and available subcommands.
	currentName, errName := getCurrentGitConfig("user.name")
	currentEmail, errEmail := getCurrentGitConfig("user.email")
	if errName != nil || currentName == "" || errEmail != nil || currentEmail == "" {
		fmt.Println("No active git profile is set in the global configuration.")
	} else {
		fmt.Printf("Current active git profile: %s - %s\n", currentName, currentEmail)
		profiles, err := loadProfiles()
		imported := false
		if err == nil {
			for _, p := range profiles {
				if p.Name == currentName && p.Email == currentEmail {
					imported = true
					break
				}
			}
		}
		if !imported {
			fmt.Println("Note: The current git profile is not imported in the profiles list.")
			fmt.Println("You can import it using the 'import' subcommand.")
		}
	}
	fmt.Println("Available subcommands: add, delete, list, switch, import")

	// If no subcommand is provided, display usage information and exit.
	if len(os.Args) < 2 {
		fmt.Println("Usage: git-user [add|delete|list|switch|import]")
		os.Exit(0)
	}

	cmd := os.Args[1]
	switch cmd {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		name := addCmd.String("name", "", "Git user name")
		email := addCmd.String("email", "", "Git user email")
		addCmd.Parse(os.Args[2:])
		if *name == "" || *email == "" {
			fmt.Println("Both --name and --email flags are required")
			os.Exit(1)
		}
		if err := addProfile(*name, *email); err != nil {
			log.Fatalf("Error adding profile: %v", err)
		}
		fmt.Println("Profile added successfully.")

	case "delete":
		delCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		name := delCmd.String("name", "", "Git user name to delete")
		delCmd.Parse(os.Args[2:])
		if *name == "" {
			fmt.Println("--name flag is required")
			os.Exit(1)
		}
		if err := deleteProfile(*name); err != nil {
			log.Fatalf("Error deleting profile: %v", err)
		}
		fmt.Println("Profile deleted successfully.")

	case "list":
		if err := listProfiles(); err != nil {
			log.Fatalf("Error listing profiles: %v", err)
		}

	case "switch":
		if err := switchProfile(); err != nil {
			log.Fatalf("Error switching profile: %v", err)
		}

	case "import":
		if err := importProfile(); err != nil {
			log.Fatalf("Error importing current profile: %v", err)
		}

	default:
		fmt.Println("Unknown command. Usage: git-user [add|delete|list|switch|import]")
		os.Exit(1)
	}
}

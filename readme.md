# git-user

**git-user** is a simple command-line utility written in Go that allows users to manage multiple Git profiles. It provides functionality to switch between different Git users, list available profiles, and import the current Git user configuration.

## Features

- Add a new Git user profile
- Delete an existing profile
- List all stored Git profiles (with the currently active profile highlighted)
- Switch between profiles interactively
- Import the currently active Git configuration into stored profiles

## Installation

### Prerequisites
- **Go 1.20+** is required to build from source.
- Ensure that **Git** is installed on your system.

### Install via Script
To install the program, run the following command:

```sh
curl -sL https://raw.githubusercontent.com/consulitsk/git-user/refs/heads/main/setup.sh | bash
```

### Build from Source
To build the project manually, run the following commands:

```sh
git clone https://github.com/consulitsk/git-user.git
cd git-user
go build -ldflags="-s -w" -o git-user main.go
```

### Using Prebuilt Binaries
If you prefer not to build from source, you can download the latest prebuilt binary from the [GitHub Releases](https://github.com/consulitsk/git-user/releases) section.

## Usage

### List Available Subcommands
```sh
git-user
```
This will display the current active Git profile and the available subcommands.

### Add a New Profile
```sh
git-user add --name "John Doe" --email "john@example.com"
```

### Delete a Profile
```sh
git-user delete --name "John Doe"
```

### List Profiles
```sh
git-user list
```
This command will show all stored profiles, with the currently active profile marked.

### Switch Git Profile
```sh
git-user switch
```
An interactive menu will appear where you can select the profile to switch to.

### Import Current Git Profile
```sh
git-user import
```
If your current Git configuration (`git config --global`) is not in the stored profiles, this command will add it.

## Configuration Storage
Profiles are stored in:
```
~/.config/git-user/profiles.json
```
This file contains all stored Git profiles in JSON format.

## Running Tests
To run the test suite, use:
```sh
go test ./...
```

## GitHub Actions Pipeline
The repository includes a GitHub Actions pipeline that:
- Runs tests
- Builds the binary
- Uploads the built artifact

You can find the workflow file at `.github/workflows/ci.yml`.

## Contributing
1. Fork this repository
2. Create a new branch: `git checkout -b feature-branch`
3. Commit your changes: `git commit -m "Add new feature"`
4. Push the branch: `git push origin feature-branch`
5. Open a Pull Request

## License
This project is licensed under the **MIT License**. See [LICENSE](LICENSE) for details.

package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
)

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// readToken reads a token from a file.
func readToken(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		path = strings.Replace(path, "~", usr.HomeDir, 1)
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read token from %s: %w", path, err)
	}
	return strings.TrimSpace(string(content)), nil
}

// GetTokenFromFile reads a token from a file.
// This is a duplicate of readToken for now to satisfy the test, but should be removed.
func GetTokenFromFile(path string) (string, error) {
	return readToken(path)
}

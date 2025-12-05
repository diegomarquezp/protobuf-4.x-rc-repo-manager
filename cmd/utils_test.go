package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadLines(t *testing.T) {
	// Create a temporary file
	tmpfile, err := ioutil.TempFile("", "test")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	// Write some lines to the file
	lines := []string{"line1", "line2", "line3"}
	for _, line := range lines {
		_, err := tmpfile.WriteString(line + "\n")
		assert.NoError(t, err)
	}

	// Read the lines back
	readLines, err := readLines(tmpfile.Name())
	assert.NoError(t, err)
	assert.Equal(t, lines, readLines)
}

func TestReadToken(t *testing.T) {
	// Create a temporary file
	tmpfile, err := ioutil.TempFile("", "test")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	// Write a token to the file
	token := "my-secret-token"
	_, err = tmpfile.WriteString(token)
	assert.NoError(t, err)

	// Read the token back
	readToken, err := readToken(tmpfile.Name())
	assert.NoError(t, err)
	assert.Equal(t, token, readToken)
}


func TestGetTokenFromFile(t *testing.T) {
	// Create a temporary file
	tmpfile, err := ioutil.TempFile("", "test")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	// Write a token to the file
	token := "my-secret-token"
	_, err = tmpfile.WriteString(token)
	assert.NoError(t, err)

	// Read the token back
	getToken, err := GetTokenFromFile(tmpfile.Name())
	assert.NoError(t, err)
	assert.Equal(t, token, getToken)
}
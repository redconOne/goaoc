package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetInput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "test input data")
	}))
	defer server.Close()

	year, day := 2000, 1
	dirName, err := GetInput(year, day, "test-session-cookie", server.URL)
	if err != nil {
		t.Errorf("GetInput returned an error: %v", err)
	}

	defer func() {
		dir := fmt.Sprint(year)
		if err := os.RemoveAll(dir); err != nil {
			t.Errorf("Error cleaning up directory: %v", err)
		}
	}()

	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		t.Errorf("Directory was not created")
	}

	filePath := fmt.Sprintf("%s/input.txt", dirName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("Input file was not created")
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Error reading input.txt file: %v", err)
	}
	expected := "test input data\n"
	if string(content) != expected {
		t.Errorf("Content of input.txt does not match. Got: %s, Expected: %s", string(content), expected)
	}
}

func TestCreateFile(t *testing.T) {
	filePath := "tempFile.txt"
	body := "test body data"

	err := CreateFile(filePath, body)
	if err != nil {
		t.Errorf("File creation failed: %s", err)
	}

	defer func() {
		if err := os.Remove(filePath); err != nil {
			t.Errorf("Error removing file: %s", err)
		}
	}()

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Error reading tempFile.txt file: %v", err)
	}
	expected := "test body data"
	if string(content) != expected {
		t.Errorf("COntent of tempFile.txt does not match. Got: %s, Expected %s", string(content), expected)
	}
}

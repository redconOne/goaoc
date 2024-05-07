package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func GetInput(year, day int, sessionCookie, server string) (string, error) {
	if sessionCookie == "" {
		return "", fmt.Errorf("session cookie missing")
	}

	dirName := fmt.Sprintf("%d/%02d", year, day)
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/%d/day/%d/input", server, year, day)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.AddCookie(&http.Cookie{Name: "session", Value: sessionCookie})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to retrieve input: %s", resp.Status)
	}

	filePath := fmt.Sprintf("%s/input.txt", dirName)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %s", err)
	}

	err = CreateFile(filePath, string(body))
	if err != nil {
		return "", fmt.Errorf("failed to create file: %s", err)
	}

	return dirName, nil
}

func CreateFile(filePath, content string) error {
	outputFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = outputFile.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

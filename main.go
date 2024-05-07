package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/joho/godotenv"
)

func main() {
	sessionCookie := ""
	if err := godotenv.Load(); err != nil {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Please enter your Advent of Code session cookie: ").
					Description("This can be found in your browser's cookies section as 'session'").
					Value(&sessionCookie),
			),
		)
		err = form.Run()
		if err != nil {
			fmt.Println("Error: ", err)
		}

		envBody := fmt.Sprintf("SESSION_COOKIE=%v", &sessionCookie)
		CreateFile(".env", envBody)
	} else {
		sessionCookie = os.Getenv("SESSION_COOKIE")
	}

	day := 1
	year := 2015

	dirName, err := GetInput(year, day, sessionCookie, "https://adventofcode.com/")
	if err != nil {
		fmt.Printf("error creating input: %s", err)
		return
	}

	fmt.Printf("Successfully created input file at %s/input.txt\n", dirName)
}

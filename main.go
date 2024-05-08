package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/joho/godotenv"
)

// TODO: Make years input more robust, EG allow only actually available years
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

	dayStr := "1"
	year := 2015
	years := []int{}

	for i := 0; i < 8; i++ {
		years = append(years, 2015+i)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Options(huh.NewOptions(years...)...).
				Title("Year").
				Description("Choose a year of advent of code to participate in.").
				Value(&year),
			huh.NewInput().
				Title("Day").
				Description("Choose a day of advent of code to participate in between 1-31").
				Placeholder("1").
				Validate(func(str string) error {
					_, err := strconv.Atoi(str)
					if err != nil {
						return errors.New("please enter a valid day")
					}
					return nil
				}).
				Value(&dayStr),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatalf("error running form: %s", err)
	}

	day, err := strconv.Atoi(dayStr)
	if err != nil {
		log.Fatalf("error converted day string to day int: %s", err)
	}

	_, err = GetInput(year, day, sessionCookie, "https://adventofcode.com")
	if err != nil {
		log.Fatalf("error creating input: %s", err)
	}
}

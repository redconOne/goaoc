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

// TODO: Move all form logic into separate file

// TODO: Instead of prompting user for cookie pull directly from browser
func getSessionCookie() (string, error) {
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
			return "", err
		}

		envBody := fmt.Sprintf("SESSION_COOKIE=%v", &sessionCookie)
		CreateFile(".env", envBody)
	} else {
		sessionCookie = os.Getenv("SESSION_COOKIE")
	}

	return sessionCookie, nil
}

// TODO: Make years input more robust, EG allow only actually available years
// TODO: Consider creating an object to return instead of 3 variables @ getDate()
func getDate() (int, int, error) {
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
		return -1, -1, err
	}

	day, err := strconv.Atoi(dayStr)
	if err != nil {
		return -1, -1, err
	}

	return year, day, nil
}

func createSolutionFile(year, day int, dirName string) error {
	solutionPath := fmt.Sprintf("%s/solution.go", dirName)
	solutionBody := []byte(fmt.Sprintf(`package year%v_day%02v

  func solve1(input string) string {
	  // Implement solution for part 1
  }

  func solve2(input string) string {
	  // Implement solution for part 2
  }`, year, day))

	err := CreateFile(solutionPath, string(solutionBody))
	if err != nil {
		return err
	}

	return nil
}

func getExamples() (map[string]string, error) {
	examples := make(map[string]string)
	moreExamples := true

	for moreExamples {
		var exampleInput, exampleResult string

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Example").
					Description("Please enter an example input").
					Value(&exampleInput).
					Validate(func(str string) error {
						if str == "" {
							return errors.New("please enter an example input")
						}
						return nil
					}),
				huh.NewInput().
					Title("Expected").
					Description("Please enter the expected result").
					Validate(func(str string) error {
						if str == "" {
							return errors.New("please enter a result")
						}
						return nil
					}).
					Value(&exampleResult),
				huh.NewConfirm().
					Title("Add example?").
					Description("Would you like to add this example and answer to your test cases?").
					Affirmative("ADD").
					Negative("CANCEL").
					Value(&moreExamples),
			),
		)

		err := form.Run()
		if err != nil {
			return examples, err
		}

		if moreExamples {
			examples[exampleInput] = exampleResult

			form := huh.NewForm(
				huh.NewGroup(
					huh.NewConfirm().
						Title("Add another?").
						Affirmative("Add another").
						Negative("I'm done here").
						Value(&moreExamples),
				),
			)

			err := form.Run()
			if err != nil {
				return examples, err
			}
		}
	}

	return examples, nil
}

func createTestFiles(year, day int, dirName string, p1Examples, p2Examples map[string]string) error {
	p1TestCases, p2TestCases := "", ""

	for input, want := range p1Examples {
		p1TestCases += fmt.Sprintf(`{
      input: "%s",
      want: "%s",
    },`, input, want)
	}

	for input, want := range p2Examples {
		p2TestCases += fmt.Sprintf(`{
      input: "%s",
      want: "%s",
    },`, input, want)
	}

	testFilePath := fmt.Sprintf("%s/solution_test.go", dirName)
	testFileBody := fmt.Sprintf(`package year%v_day%02v

import "testing"

func TestSolve1(t *testing.T) {
  testcases := []struct{
    input, want string
  }{
  %s
  }
    
  for _, tc := range testcases {
    result := solve1(tc.input)
    if result != tc.want {
      t.Errorf("Part1 failed - Result: %%s Expected: %%s", result, tc.want)
    }
  }
}

func TestSolve2(t *testing.T) {
  testcases := []struct{
    input, want string
  }{
  %s
  }
    
  for _, tc := range testcases {
    result := solve2(tc.input)
    if result != tc.want {
      t.Errorf("Part2 failed - Result: %%s Expected: %%s", result, tc.want)
    }
  }
}
    `, year, day, p1TestCases, p2TestCases)

	err := CreateFile(testFilePath, string(testFileBody))
	if err != nil {
		return err
	}

	return nil
}

// TODO: Create a main menu (Fetch input, run tests, solve day, about, update)
func main() {
	sessionCookie, err := getSessionCookie()
	if err != nil {
		log.Fatalf("error getting session cookie: %s", err)
	}

	year, day, err := getDate()
	if err != nil {
		log.Fatalf("error getting year/day: %s", err)
	}

	dirName, err := GetInput(year, day, sessionCookie, "https://adventofcode.com")
	if err != nil {
		log.Fatalf("error creating input: %s", err)
	}

	err = createSolutionFile(year, day, dirName)
	if err != nil {
		log.Fatalf("error creating solutions file: %s", err)
	}

	p1Examples, err := getExamples()
	if err != nil {
		log.Fatalf("error getting examples: %s", err)
	}

	p2Examples, err := getExamples()
	if err != nil {
		log.Fatalf("error getting examples: %s", err)
	}

	err = createTestFiles(year, day, dirName, p1Examples, p2Examples)
	if err != nil {
		log.Fatalf("error creating test files: %s", err)
	}
}

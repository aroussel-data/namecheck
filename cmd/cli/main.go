package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/aroussel-data/namecheck"
	"github.com/aroussel-data/namecheck/github"
	"github.com/aroussel-data/namecheck/twitter"
)

type Result struct {
	Username  string
	Platform  string
	Valid     bool
	Available bool
	Err       error
}

func check(checker namecheck.Checker, username string, wg *sync.WaitGroup, rc chan Result) {
	defer wg.Done()

	res := Result{Username: username, Platform: checker.String()}

	if checker.IsValid(username) {
		res.Valid = true
		res.Available, res.Err = checker.IsAvailable(username)
	}
	rc <- res
}

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "expected a username")
		os.Exit(1)
	}

	username := args[1]

	tw := twitter.Twitter{
		Getter: http.DefaultClient,
	}
	gh := github.Github{
		Getter: http.DefaultClient,
	}

	checkers := []namecheck.Checker{&tw, &gh}

	resultChan := make(chan Result)

	var wg sync.WaitGroup

	wg.Add(2)

	for _, checker := range checkers {
		go check(checker, username, &wg, resultChan)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var results []Result
	for res := range resultChan {
		results = append(results, res)
	}
	fmt.Println(results)
}

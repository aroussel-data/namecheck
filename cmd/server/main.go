package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/aroussel-data/namecheck"
	"github.com/aroussel-data/namecheck/github"
	"github.com/aroussel-data/namecheck/twitter"
	"github.com/julienschmidt/httprouter"
)

type Result struct {
	Username  string `json:"username"`
	Platform  string `json:"platform"`
	Valid     bool   `json:"valid"`
	Available bool   `json:"available"`
}

type Error struct {
	Err error `json:"err"`
}

var (
	m  = make(map[string]uint)
	mu sync.Mutex
)

func check(ctx context.Context, checker namecheck.Checker, username string, wg *sync.WaitGroup, rc chan Result, ec chan Error) {
	defer wg.Done()

	res := Result{Username: username, Platform: checker.String()}

	if checker.IsValid(username) {
		res.Valid = true
		available, errRes := checker.IsAvailable(username)
		if errRes != nil {
			select {
			case <-ctx.Done():
				return
			case ec <- Error{Err: errRes}:
				return
			}
		}
		res.Available = available
	}
	rc <- res
}

func handleStats(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	mu.Lock()
	defer mu.Unlock()
	if err := enc.Encode(m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func handleCheck(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	username := req.URL.Query().Get("username")

	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mu.Lock()
	m[username]++
	mu.Unlock()

	tw := twitter.Twitter{
		Getter: http.DefaultClient,
	}
	gh := github.Github{
		Getter: http.DefaultClient,
	}

	checkers := []namecheck.Checker{&tw, &gh}

	resultChan := make(chan Result)
	errorChan := make(chan Error)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	wg.Add(2)

	for _, checker := range checkers {
		go check(ctx, checker, username, &wg, resultChan, errorChan)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var results []Result

	var finished bool
	for !finished {
		select {
		case <-errorChan:
			w.WriteHeader(http.StatusInternalServerError)
			return
		case res, ok := <-resultChan:
			if !ok {
				finished = true
				continue
			}
			results = append(results, res)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(results); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func main() {
	router := httprouter.New()
	router.GET("/check", handleCheck)
	router.GET("/stats", handleStats)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

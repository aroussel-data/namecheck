package github_test

import (
	"net/http"
	"testing"

	"github.com/aroussel-data/namecheck/github"
	"github.com/aroussel-data/namecheck/stub"
)

func TestUsernameContainsTwoHyphens(t *testing.T) {
	const (
		username = "-Alex1990-"
		want     = false
	)
	var gh = github.Github{}
	got := gh.IsValid(username)
	if got != want {
		t.Errorf("twitter.IsValid(%q): got %t; want %t", username, got, want)
	}
}

func TestUsernameContainsOneHyphen(t *testing.T) {
	const (
		username = "Alex1990-"
		want     = false
	)
	var gh = github.Github{}
	got := gh.IsValid(username)
	if got != want {
		t.Errorf("twitter.IsValid(%q): got %t; want %t", username, got, want)
	}
}

func TestIsAvailable200(t *testing.T) {
	client := stub.SuccessfulGetter{StatusCode: http.StatusOK}
	gh := github.Github{
		Getter: &client,
	}
	avail, err := gh.IsAvailable("alex1990")
	if err != nil {
		t.Errorf("want nil got error; got %v", err)
	}
	if avail {
		t.Errorf("want false got true")
	}
}

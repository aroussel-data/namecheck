package twitter_test

import (
	"testing"

	"github.com/aroussel-data/namecheck/twitter"
)

func TestIsValid(t *testing.T) {
	type TestCase struct {
		desc     string
		username string
		want     bool
	}
	testCases := []TestCase{
		{"contains 'Twitter'", "jub0bsOnTwitter", false},
		{"does not contain 'Twitter'", "alexRoussel", true},
		{"illegal chars", "-alexRoussel-", false},
		{"too short", "al", false},
	}
	var tw = twitter.Twitter{}
	for _, tc := range testCases {
		f := func(t *testing.T) {
			got := tw.IsValid(tc.username)
			if got != tc.want {
				t.Errorf("twitter.IsValid(%q): got %t; want %t", tc.username, got, tc.want)
			}
		}

		t.Run(tc.desc, f)
	}
}

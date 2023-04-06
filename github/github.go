package github

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/aroussel-data/namecheck"
)

var re = regexp.MustCompile("^[a-zA-Z0-9-]{3,39}$")

type Github struct {
	Getter namecheck.Getter
}

func (t *Github) String() string {
	return "Github"
}

func noConsecutiveHyphens(username string) bool {
	return !strings.Contains(strings.ToLower(username), "--")
}

func looksGood(username string) bool {
	return re.MatchString(username)
}

func noHyphenPrefixSuffix(username string) bool {
	return !strings.HasPrefix(username, "-") && !strings.HasSuffix(username, "-")
}

func (*Github) IsValid(username string) bool {
	return noConsecutiveHyphens(username) && looksGood(username) && noHyphenPrefixSuffix(username)
}

func (g *Github) IsAvailable(username string) (bool, error) {
	endpoint := fmt.Sprintf("https://github.com/%s", username)
	resp, err := g.Getter.Get(endpoint)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		return false, nil
	case http.StatusNotFound:
		return true, nil
	default:
		return false, errors.New("unknown availability")
	}
}

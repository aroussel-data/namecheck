package twitter

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/aroussel-data/namecheck"
)

var re = regexp.MustCompile("^[a-zA-Z0-9_]{4,15}$")

type Twitter struct {
	Getter namecheck.Getter
}

func (t *Twitter) String() string {
	return "Twitter"
}

func containsNoIllegalPattern(username string) bool {
	return !strings.Contains(strings.ToLower(username), "twitter")
}

func looksGood(username string) bool {
	return re.MatchString(username)
}

func (*Twitter) IsValid(username string) bool {
	return containsNoIllegalPattern(username) && looksGood(username)
}

func (t *Twitter) IsAvailable(username string) (bool, error) {
	const tmpl = "https://europe-west6-namechecker-api.cloudfunctions.net/userlookup?username=%s"
	endpoint := fmt.Sprintf(tmpl, username)
	resp, err := t.Getter.Get(endpoint)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, errors.New("unknown availability")
	}
	var dto struct {
		Data any `json:"data"`
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&dto); err != nil {
		return false, err
	}
	return dto.Data == nil, nil
}

package stub

import (
	"io"
	"net/http"
	"os"
)

type SuccessfulGetter struct {
	StatusCode int
}

func (sg *SuccessfulGetter) Get(url string) (resp *http.Response, err error) {
	return &http.Response{
			StatusCode: sg.StatusCode,
			Body:       io.NopCloser(os.Stdin),
		},
		nil
}

package http

import (
	"errors"
	"io/ioutil"
	"net/http"
)

// GET implementing http get request
func GET(url string) ([]byte, error) {
	// start get request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// check http status
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	// read data
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

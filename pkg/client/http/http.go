package http

import (
	"errors"
	"io/ioutil"
	"net/http"
)

func HTTPGET(url string) ([]byte, error) {
	resp, err := http.Get(url + "/stopped")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

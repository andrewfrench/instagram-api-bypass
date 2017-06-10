package request

import (
	"errors"
	"net/http"
	"io/ioutil"
)

func Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode != 200 {
		return []byte{}, errors.New(resp.Status)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)

	return respBytes, err
}

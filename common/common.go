package common

import (
	"log"
	"regexp"
	"errors"
)

func ExtractJson(input []byte) ([]byte, error) {
	pattern := "<script type=\"text/javascript\">window._sharedData = (.*);</script>"
	r, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatal(err.Error())
	}

	submatches := r.FindSubmatch(input)
	if len(submatches) != 2 {
		return []byte{}, errors.New("Failed to extract JSON from raw GET request")
	}

	return submatches[1], err
}

package utils

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("authorization header is missing")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("authorization header is invalid")
	}

	if vals[0] != "Bearer" {
		return "", errors.New("authorization header is invalid")
	}

	return vals[1], nil
}

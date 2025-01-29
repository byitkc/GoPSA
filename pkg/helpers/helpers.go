package helpers

import (
	"errors"
	"net"
	"net/http"
)

var (
	ErrUnexpectedStatusCode = errors.New("unexpected status code returned")
)

func DoWithRetry(c *http.Client, req *http.Request, expectedStatusCode, retryCount, timeout int) (http.Response, error) {
	var resp http.Response
	for i := 1; i < retryCount; i++ {
		resp, err := c.Do(req)
		if i != retryCount {
			if err != nil {
				continue
			}
			if resp.StatusCode != expectedStatusCode {
				continue
			}
		}
		if i >= retryCount {
			if err != nil {
				return http.Response{}, err
			}
			if resp.StatusCode != expectedStatusCode {
				return http.Response{}, ErrUnexpectedStatusCode
			}
		}
	}
	return resp, nil
}

func IsResolvable(host string) bool {
	addrs, err := net.LookupHost(host)
	if err != nil {
		return false
	}
	return len(addrs) > 0
}
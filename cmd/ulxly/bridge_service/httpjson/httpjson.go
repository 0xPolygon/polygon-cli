package httpjson

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// Helper function to get the appropriate HTTP client
func NewHTTPClient(insecure bool) *http.Client {
	if insecure {
		log.Warn().Msg("WARNING: Using insecure HTTP client for bridge service requests")
		return &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}

// Get makes an HTTP GET request to the specified URL and unmarshals the JSON response into the provided generic type T.
func HTTPGet[T any](client *http.Client, url string) (obj T, statusCode int, err error) {
	res, err := client.Get(url)
	if err != nil {
		return obj, 0, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return obj, res.StatusCode, err
	}
	defer res.Body.Close()

	err = json.Unmarshal(body, &obj)

	return obj, res.StatusCode, err
}

// Get makes an HTTP GET request to the specified URL and unmarshals the JSON response into the provided generic type T.
func HTTPGetWithError[T any, TError any](client *http.Client, url string) (obj T, objError TError, statusCode int, err error) {
	res, err := client.Get(url)
	if err != nil {
		return obj, objError, 0, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return obj, objError, res.StatusCode, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = json.Unmarshal(body, &obj)
	} else {
		err = json.Unmarshal(body, &objError)
	}

	return obj, objError, res.StatusCode, err
}

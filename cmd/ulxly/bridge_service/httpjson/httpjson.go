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
func HTTPGet[T any](client *http.Client, url string) (T, error) {
	var obj T
	res, err := client.Get(url)
	if err != nil {
		return obj, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return obj, err
	}
	defer res.Body.Close()

	err = json.Unmarshal(body, &obj)
	if err != nil {
		return obj, err
	}

	return obj, nil
}

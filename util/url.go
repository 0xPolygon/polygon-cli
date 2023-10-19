package util

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/rs/zerolog/log"
)

// ValidateUrl checks if a string URL can be parsed and if it has a valid scheme.
func ValidateUrl(input string) error {
	url, err := url.Parse(input)
	if err != nil {
		log.Error().Err(err).Msg("Unable to parse url input error")
		return err
	}

	if url.Scheme == "" {
		return errors.New("the scheme has not been specified")
	}
	switch url.Scheme {
	case "http", "https", "ws", "wss":
		return nil
	default:
		return fmt.Errorf("the scheme '%s' is not supported", url.Scheme)
	}
}

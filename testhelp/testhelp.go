// Package testhelp are re-usable test helpers for all modules in reach. 
package testhelp

import (
	"fmt"
	"net/url"
	"testing"
)

// Reporter is this package's name for the standard Formatter.
type Reporter func(message string, values ...interface{})

// Errmsg returns a function that works like fmt.Errorf, but logs the test information
// and runs it through t's Errorf.
func Errmsg(t *testing.T, instance string) Reporter {
	return func(message string, values ...interface{}) {
		t.Errorf(fmt.Sprintf("%s: %s", instance, message), values...)
	}
}

// NewURL creates a new url, and fails the test if it's invalid.
func NewURL(t *testing.T, u string) *url.URL {
	url, err := url.Parse(u)
	if err != nil {
		t.Error(err)
	}
	return url
}

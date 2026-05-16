// Package jwt will implement JWT signing and verification in M7.
// Until then, Sign and Verify return stub errors.
package jwt

import "errors"

// Claims holds the JWT payload. Populated in M7.
type Claims struct {
	UserID int64  `json:"sub"`
	Role   string `json:"role"`
}

var errNotImplemented = errors.New("JWT auth is not implemented until M7")

// Sign returns a stub error. Implemented in M7.
func Sign(_ Claims, _ string) (string, error) {
	return "", errNotImplemented
}

// Verify returns a stub error. Implemented in M7.
func Verify(_, _ string) (Claims, error) {
	return Claims{}, errNotImplemented
}

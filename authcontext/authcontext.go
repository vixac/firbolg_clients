// Package authcontext implements the HTTP authentication-context contract used
// between the Firbolg gateway and its upstream services.
package authcontext

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// HTTP header names in the Firbolg gateway authentication-context contract.
const (
	HeaderAuthVersion     = "X-Firbolg-Auth-Version"
	HeaderUserID          = "X-Firbolg-User-ID"
	HeaderUserEmail       = "X-Firbolg-User-Email"
	HeaderUserGivenName   = "X-Firbolg-User-Given-Name"
	HeaderUserDisplayName = "X-Firbolg-User-Display-Name"
	HeaderRequestID       = "X-Request-ID"

	// AuthVersionV1 is the only currently supported gateway contract version.
	AuthVersionV1 = "1"
)

// AuthContext is the identity asserted by the Firbolg gateway for a request.
// UserGivenName and UserDisplayName are optional.
type AuthContext struct {
	UserID          string
	UserEmail       string
	UserGivenName   string
	UserDisplayName string
	RequestID       string
}

var (
	ErrMissingHeader      = errors.New("missing required authentication-context header")
	ErrDuplicateHeader    = errors.New("duplicate authentication-context header")
	ErrUnsupportedVersion = errors.New("unsupported authentication-context version")
	ErrInvalidUserID      = errors.New("invalid authentication-context user ID")
)

// HeaderError identifies the header that made an authentication context invalid.
type HeaderError struct {
	Header string
	Err    error
}

func (e *HeaderError) Error() string {
	return fmt.Sprintf("%s: %v", e.Header, e.Err)
}

func (e *HeaderError) Unwrap() error { return e.Err }

// ValidateV1 validates an authentication context for version 1 of the gateway
// contract. User IDs must be Google principal IDs: google:<Google subject>.
func ValidateV1(context AuthContext) error {
	if subject, ok := strings.CutPrefix(context.UserID, "google:"); !ok || strings.TrimSpace(subject) == "" {
		return &HeaderError{Header: HeaderUserID, Err: ErrInvalidUserID}
	}
	if strings.TrimSpace(context.UserEmail) == "" {
		return &HeaderError{Header: HeaderUserEmail, Err: ErrMissingHeader}
	}
	if strings.TrimSpace(context.RequestID) == "" {
		return &HeaderError{Header: HeaderRequestID, Err: ErrMissingHeader}
	}
	return nil
}

// ValidateV1 validates c for version 1 of the gateway contract.
func (c AuthContext) ValidateV1() error { return ValidateV1(c) }

// ParseRequest reads and validates the gateway authentication context from r.
// Missing, duplicate, or unsupported critical headers are rejected.
func ParseRequest(r *http.Request) (*AuthContext, error) {
	if r == nil {
		return nil, errors.New("nil request")
	}

	version, err := requiredHeader(r.Header, HeaderAuthVersion)
	if err != nil {
		return nil, err
	}
	if version != AuthVersionV1 {
		return nil, &HeaderError{Header: HeaderAuthVersion, Err: ErrUnsupportedVersion}
	}

	userID, err := requiredHeader(r.Header, HeaderUserID)
	if err != nil {
		return nil, err
	}
	email, err := requiredHeader(r.Header, HeaderUserEmail)
	if err != nil {
		return nil, err
	}
	requestID, err := requiredHeader(r.Header, HeaderRequestID)
	if err != nil {
		return nil, err
	}

	context := &AuthContext{
		UserID:          userID,
		UserEmail:       email,
		UserGivenName:   headerValue(r.Header, HeaderUserGivenName),
		UserDisplayName: headerValue(r.Header, HeaderUserDisplayName),
		RequestID:       requestID,
	}
	if err := context.ValidateV1(); err != nil {
		return nil, err
	}
	return context, nil
}

func requiredHeader(headers http.Header, name string) (string, error) {
	values := headerValues(headers, name)
	switch len(values) {
	case 0:
		return "", &HeaderError{Header: name, Err: ErrMissingHeader}
	case 1:
		if strings.TrimSpace(values[0]) == "" {
			return "", &HeaderError{Header: name, Err: ErrMissingHeader}
		}
		return values[0], nil
	default:
		return "", &HeaderError{Header: name, Err: ErrDuplicateHeader}
	}
}

func headerValues(headers http.Header, name string) []string {
	var values []string
	for key, keyValues := range headers {
		if strings.EqualFold(key, name) {
			values = append(values, keyValues...)
		}
	}
	return values
}

func headerValue(headers http.Header, name string) string {
	values := headerValues(headers, name)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

// WriteHeaders validates context and writes its version-1 representation to
// headers. It replaces any values already present for contract headers.
func WriteHeaders(headers http.Header, context AuthContext) error {
	if headers == nil {
		return errors.New("nil headers")
	}
	if err := context.ValidateV1(); err != nil {
		return err
	}

	setHeader(headers, HeaderAuthVersion, AuthVersionV1)
	setHeader(headers, HeaderUserID, context.UserID)
	setHeader(headers, HeaderUserEmail, context.UserEmail)
	setHeader(headers, HeaderRequestID, context.RequestID)
	writeOptionalHeader(headers, HeaderUserGivenName, context.UserGivenName)
	writeOptionalHeader(headers, HeaderUserDisplayName, context.UserDisplayName)
	return nil
}

// WriteHeaders writes c as version-1 authentication-context headers.
func (c AuthContext) WriteHeaders(headers http.Header) error {
	return WriteHeaders(headers, c)
}

func writeOptionalHeader(headers http.Header, name, value string) {
	if value == "" {
		deleteHeader(headers, name)
		return
	}
	setHeader(headers, name, value)
}

func setHeader(headers http.Header, name, value string) {
	deleteHeader(headers, name)
	headers[name] = []string{value}
}

func deleteHeader(headers http.Header, name string) {
	for key := range headers {
		if strings.EqualFold(key, name) {
			delete(headers, key)
		}
	}
}

package authcontext

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func validHeaders() http.Header {
	return http.Header{
		HeaderAuthVersion: {AuthVersionV1},
		HeaderUserID:      {"google:10987654321"},
		HeaderUserEmail:   {"ada@example.test"},
		HeaderRequestID:   {"request-123"},
	}
}

func TestParseRequest(t *testing.T) {
	valid := AuthContext{
		UserID: "google:10987654321", UserEmail: "ada@example.test", RequestID: "request-123",
	}
	cases := []struct {
		name    string
		mutate  func(http.Header)
		want    *AuthContext
		wantErr error
	}{
		{name: "required headers", want: &valid},
		{name: "optional headers", mutate: func(h http.Header) {
			setHeader(h, HeaderUserGivenName, "Ada")
			setHeader(h, HeaderUserDisplayName, "Ada Lovelace")
		}, want: &AuthContext{UserID: valid.UserID, UserEmail: valid.UserEmail, UserGivenName: "Ada", UserDisplayName: "Ada Lovelace", RequestID: valid.RequestID}},
		{name: "missing version", mutate: func(h http.Header) { deleteHeader(h, HeaderAuthVersion) }, wantErr: ErrMissingHeader},
		{name: "blank version", mutate: func(h http.Header) { setHeader(h, HeaderAuthVersion, " ") }, wantErr: ErrMissingHeader},
		{name: "unsupported version", mutate: func(h http.Header) { setHeader(h, HeaderAuthVersion, "2") }, wantErr: ErrUnsupportedVersion},
		{name: "duplicate version", mutate: func(h http.Header) { h.Add(HeaderAuthVersion, AuthVersionV1) }, wantErr: ErrDuplicateHeader},
		{name: "missing user ID", mutate: func(h http.Header) { deleteHeader(h, HeaderUserID) }, wantErr: ErrMissingHeader},
		{name: "duplicate user ID", mutate: func(h http.Header) { h.Add(HeaderUserID, "google:other") }, wantErr: ErrDuplicateHeader},
		{name: "wrong identity provider", mutate: func(h http.Header) { setHeader(h, HeaderUserID, "oidc:10987654321") }, wantErr: ErrInvalidUserID},
		{name: "empty Google subject", mutate: func(h http.Header) { setHeader(h, HeaderUserID, "google: ") }, wantErr: ErrInvalidUserID},
		{name: "missing email", mutate: func(h http.Header) { deleteHeader(h, HeaderUserEmail) }, wantErr: ErrMissingHeader},
		{name: "blank email", mutate: func(h http.Header) { setHeader(h, HeaderUserEmail, " ") }, wantErr: ErrMissingHeader},
		{name: "duplicate email", mutate: func(h http.Header) { h.Add(HeaderUserEmail, "other@example.test") }, wantErr: ErrDuplicateHeader},
		{name: "missing request ID", mutate: func(h http.Header) { deleteHeader(h, HeaderRequestID) }, wantErr: ErrMissingHeader},
		{name: "blank request ID", mutate: func(h http.Header) { setHeader(h, HeaderRequestID, " ") }, wantErr: ErrMissingHeader},
		{name: "duplicate request ID", mutate: func(h http.Header) { h.Add(HeaderRequestID, "request-456") }, wantErr: ErrDuplicateHeader},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			headers := validHeaders()
			if tc.mutate != nil {
				tc.mutate(headers)
			}
			got, err := ParseRequest(&http.Request{Header: headers})
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("ParseRequest() error = %v, want errors.Is(_, %v)", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseRequest() error = %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ParseRequest() = %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestParseRequestNil(t *testing.T) {
	if _, err := ParseRequest(nil); err == nil {
		t.Fatal("ParseRequest(nil) error = nil")
	}
}

func TestValidateV1(t *testing.T) {
	valid := AuthContext{UserID: "google:subject", UserEmail: "user@example.test", RequestID: "request"}
	cases := []struct {
		name    string
		context AuthContext
		wantErr error
	}{
		{name: "valid", context: valid},
		{name: "missing user ID", context: AuthContext{UserEmail: valid.UserEmail, RequestID: valid.RequestID}, wantErr: ErrInvalidUserID},
		{name: "invalid user ID prefix", context: AuthContext{UserID: "googleish:subject", UserEmail: valid.UserEmail, RequestID: valid.RequestID}, wantErr: ErrInvalidUserID},
		{name: "missing email", context: AuthContext{UserID: valid.UserID, RequestID: valid.RequestID}, wantErr: ErrMissingHeader},
		{name: "missing request ID", context: AuthContext{UserID: valid.UserID, UserEmail: valid.UserEmail}, wantErr: ErrMissingHeader},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateV1(tc.context)
			if tc.wantErr == nil && err != nil {
				t.Fatalf("ValidateV1() error = %v", err)
			}
			if tc.wantErr != nil && !errors.Is(err, tc.wantErr) {
				t.Fatalf("ValidateV1() error = %v, want errors.Is(_, %v)", err, tc.wantErr)
			}
		})
	}
}

func TestWriteHeaders(t *testing.T) {
	context := AuthContext{UserID: "google:subject", UserEmail: "user@example.test", UserGivenName: "Ada", UserDisplayName: "Ada Lovelace", RequestID: "request"}
	headers := http.Header{HeaderUserID: {"old", "duplicate"}, HeaderUserGivenName: {"old"}}
	if err := WriteHeaders(headers, context); err != nil {
		t.Fatalf("WriteHeaders() error = %v", err)
	}
	want := http.Header{
		HeaderAuthVersion: {AuthVersionV1}, HeaderUserID: {context.UserID}, HeaderUserEmail: {context.UserEmail},
		HeaderUserGivenName: {context.UserGivenName}, HeaderUserDisplayName: {context.UserDisplayName}, HeaderRequestID: {context.RequestID},
	}
	if !reflect.DeepEqual(headers, want) {
		t.Errorf("WriteHeaders() headers = %#v, want %#v", headers, want)
	}
}

func TestWriteHeadersClearsAbsentOptionalHeaders(t *testing.T) {
	headers := validHeaders()
	headers.Set(HeaderUserGivenName, "old")
	headers.Set(HeaderUserDisplayName, "old")
	context := AuthContext{UserID: "google:subject", UserEmail: "user@example.test", RequestID: "request"}
	if err := context.WriteHeaders(headers); err != nil {
		t.Fatalf("WriteHeaders() error = %v", err)
	}
	if headers.Get(HeaderUserGivenName) != "" || headers.Get(HeaderUserDisplayName) != "" {
		t.Errorf("optional headers were not cleared: %#v", headers)
	}
}

func TestWriteHeadersRejectsInvalidContext(t *testing.T) {
	headers := http.Header{}
	err := WriteHeaders(headers, AuthContext{UserID: "not-google", UserEmail: "user@example.test", RequestID: "request"})
	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("WriteHeaders() error = %v, want invalid user ID", err)
	}
	if len(headers) != 0 {
		t.Errorf("WriteHeaders() wrote headers on error: %#v", headers)
	}
	if err := WriteHeaders(nil, AuthContext{UserID: "google:subject", UserEmail: "user@example.test", RequestID: "request"}); err == nil {
		t.Error("WriteHeaders(nil, context) error = nil")
	}
}

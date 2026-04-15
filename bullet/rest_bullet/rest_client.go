package rest_bullet

import (
	"net/http"
	"strings"

	"github.com/vixac/bullet/store/store_interface"
	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	"github.com/vixac/firbolg_clients/bullet/util"
)

var _ bullet_interface.BulletClientInterface = (*RestClient)(nil)

type RestClient struct {
	*util.FirbolgClient
}

// Option configures a RestClient.
type Option func(*util.FirbolgClient)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(c *http.Client) Option {
	return func(f *util.FirbolgClient) {
		f.HTTPClient = c
	}
}

// WithLogger sets a logger that will print each request and response.
func WithLogger(l util.Logger) Option {
	return func(f *util.FirbolgClient) {
		f.Logger = l
	}
}

func NewRestClient(baseURL string, space store_interface.TenancySpace, opts ...Option) *RestClient {
	fc := util.NewFirbolgClient(strings.TrimRight(baseURL, "/"), int64(space.AppId), space.TenancyId)
	for _, o := range opts {
		o(fc)
	}
	return &RestClient{FirbolgClient: fc}
}

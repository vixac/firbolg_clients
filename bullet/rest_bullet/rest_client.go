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

func NewRestClient(baseURL string, space store_interface.TenancySpace) *RestClient {
	return &RestClient{
		FirbolgClient: util.NewFirbolgClient(strings.TrimRight(baseURL, "/"), int64(space.AppId), space.TenancyId),
	}
}

func NewRestClientWithHTTPClient(baseURL string, space store_interface.TenancySpace, httpClient *http.Client) *RestClient {
	client := NewRestClient(baseURL, space)
	client.HTTPClient = httpClient
	return client
}

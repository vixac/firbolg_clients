// Package rest_bullet provides legacy constructor names for Bullet's REST client.
package rest_bullet

import (
	bullet_rest "github.com/vixac/bullet/client/rest"
	"github.com/vixac/bullet/model"
)

type Option = bullet_rest.Option

var WithHTTPClient = bullet_rest.WithHTTPClient
var WithLogger = bullet_rest.WithLogger

func NewRestClient(baseURL string, space model.TenancySpace, opts ...Option) *bullet_rest.Client {
	return bullet_rest.New(baseURL, space, opts...)
}

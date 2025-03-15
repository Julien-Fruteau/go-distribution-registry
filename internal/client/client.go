package client

//  📢 do the http requests,
//
import (
	"encoding/base64"
	"net/http"

	reg "git.isi.nc/go/dtb-tool/pkg/registry"
)

type RegClient struct {
	registry   reg.Registry
	client     *http.Client
	authHeader string
}

func NewClient() *RegClient {
	r := &RegClient{
		registry: reg.NewRegistry(),
		client:   &http.Client{},
	}
	r.authHeader = getAuthHeader(r.registry.Conf.Username, r.registry.Conf.Password)
	return r
}

func getAuthHeader(username, password string) string {
	auth := username + ":" + password
	base64Auth := base64.StdEncoding.EncodeToString([]byte(auth))
	return "Basic " + base64Auth
}

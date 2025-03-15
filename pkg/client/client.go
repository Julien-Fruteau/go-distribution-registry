package client

import (
	"net/http"

	"git.isi.nc/go/dtb-tool/pkg/registry"
)

type RegistryClient struct {
  httpClient *http.Client
  registry   registry.Registry
}

func NewRegistryClient() *RegistryClient {
  return &RegistryClient{
    httpClient: &http.Client{},
    registry:   registry.NewRegistry(),
  }
}

func (r *RegistryClient) GetCatalog() ([]string, error){
  return r.registry.Catalog(r.httpClient)
}


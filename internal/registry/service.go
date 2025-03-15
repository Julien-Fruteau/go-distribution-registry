package registry

import (
	"net/http"
)

type RegistrySvc struct {
  httpClient *http.Client
  registry   Registry
}

func NewRegistrySvc() *RegistrySvc {
  return &RegistrySvc{
    httpClient: &http.Client{},
    registry:   NewRegistry(),
  }
}

func (r *RegistrySvc) GetCatalog() ([]string, error){
  return r.registry.Catalog(r.httpClient)
}


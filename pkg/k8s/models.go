package k8s

type Image struct {
  Repository string `json:"repository"`
  Tag        string `json:"tag"`
}

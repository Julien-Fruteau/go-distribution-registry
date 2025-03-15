package registry

type RegistryErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

type RegistryError struct {
	Errors []RegistryErrorDetail `json:"errors"`
}

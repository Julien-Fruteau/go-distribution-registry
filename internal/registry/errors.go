package registry

type RegistryErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

type RegistryError struct {
	Errors []RegistryErrorDetail `json:"errors"`
}


// 💥 exemple d erreur sur header non supporté
// HTTP/1.1 404 Not Found
// Content-Type: application/json; charset=utf-8
// Docker-Distribution-Api-Version: registry/2.0
// X-Content-Type-Options: nosniff
// Date: Wed, 29 Jan 2025 21:36:09 GMT
// Content-Length: 117
//
// {"errors":[{"code":"MANIFEST_UNKNOWN","message":"OCI index found, but accept header does not support OCI indexes"}]}

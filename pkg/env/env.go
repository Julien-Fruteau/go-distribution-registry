package env

import "os"

func GetEnvOrDefault(name, other string) string {
	found, ok := os.LookupEnv(name)
	if !ok {
		return other
	}
	return found
}

package env

import "os"

// Get returns the env var, or fallback if it is unset or empty.
func Get(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

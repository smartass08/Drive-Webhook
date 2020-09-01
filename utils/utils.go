package utils

import "os"

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "8000"
	}
	return port
}

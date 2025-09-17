package config

import "os"

type Config struct {
	Addr string //e.g. ":8080"
}

func FromEnv() Config {
	port := os.Getenv("PORT")
	if port == ""{
		port == "8080"
	}
	return Config{Addr: ":" + port}
}
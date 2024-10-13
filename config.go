package main

import "os"

type Config struct {
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	AWS_REGION            string
	K8S_NAMESPACE         string
	K8S_SECRET_NAME       string
	DOCKER_SERVER         string
	DOCKER_EMAIL          string
}

func NewConfig() *Config {
	return &Config{
		AWS_ACCESS_KEY_ID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWS_SECRET_ACCESS_KEY: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWS_REGION:            os.Getenv("AWS_REGION"),
		K8S_NAMESPACE:         os.Getenv("K8S_NAMESPACE"),
		K8S_SECRET_NAME:       os.Getenv("K8S_SECRET_NAME"),
		DOCKER_SERVER:         os.Getenv("DOCKER_SERVER"),
		DOCKER_EMAIL:          os.Getenv("DOCKER_EMAIL"),
	}
}

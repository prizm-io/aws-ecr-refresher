package main

import "os"

var (
	AWS_ACCESS_KEY_ID     = os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY = os.Getenv("AWS_SECRET_ACCESS_KEY")
	AWS_REGION            = os.Getenv("AWS_REGION")
	K8S_NAMESPACE         = os.Getenv("K8S_NAMESPACE")
	K8S_SECRET_NAME       = os.Getenv("K8S_SECRET_NAME")

	DOCKER_SERVER = os.Getenv("DOCKER_SERVER")
	DOCKER_EMAIL  = os.Getenv("DOCKER_EMAIL")
)

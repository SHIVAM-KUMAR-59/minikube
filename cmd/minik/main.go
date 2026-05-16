package main

import (
	"github.com/SHIVAM-KUMAR-59/minikube/cmd/minik/cmd"
	_ "github.com/SHIVAM-KUMAR-59/minikube/cmd/minik/cmd/delete"
	_ "github.com/SHIVAM-KUMAR-59/minikube/cmd/minik/cmd/get"
)

func main() {
	cmd.Execute()
}

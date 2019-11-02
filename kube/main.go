package main

import (
	edgeclusterlib "kube/microbusiness/edgeclusterlib"

	"log"
)

func main() {
	log.Println("start ... ")

	var isUsingKind bool = true

	var objDeployment edgeclusterlib.EdgeClusterDeploymentDetail
	var objService edgeclusterlib.EdgeClusterServiceDetail

	if isUsingKind {
		objDeployment.ConfigName = "kind-config-devEnv"
	} else {
		objDeployment.ConfigName = "config"
	}

	configContext := objDeployment.GetKubeConfig()

	objDeployment.ConnectToCluster(configContext)

	objDeployment.Name = "test"

	log.Print(objDeployment.Name)
	log.Print(objService.Name)

	log.Println("test")
}

package main

import (
	edgeclusterlib "kube/microbusiness/edgeclusterlib"
	"log"

	apiv1 "k8s.io/api/core/v1"
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

	objDeployment.NameSpace = apiv1.NamespaceDefault
	objDeployment.AppName = "Test1"
	objDeployment.MetaData = "Test1"
	objDeployment.ContainerImage = ""
	objDeployment.ContainerName = ""

	configContext := objDeployment.GetKubeConfig()

	objDeployment.ConnectToCluster(configContext)

	objDeployment.Name = "test"

	log.Print(objDeployment.Name)
	log.Print(objService.Name)

	log.Println("test")
}

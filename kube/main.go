package main

import (
	edgeclusterlib "kube/microbusiness/edgeclusterlib"
	"log"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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

	//objDeployment.Metaobject.NameSpace = "azizi-deployment"
	objDeployment.Metaobject.Name = "azizi-deployment"
	objDeployment.AppName = "azizi-myfirst"
	objDeployment.Replicas = 2
	objDeployment.ContainerName = "myfirst"
	objDeployment.ContainerImage = "mortezaazizi/myfirstcontainer"
	objDeployment.Args[0] = "/src-azizi2/main"
	objDeployment.Ports[0] = 8080

	objService.Metaobject.Name = "azizi-myfirst-service"

	objService.Ports = []v1.ServicePort{
		{
			Protocol:   "TCP",
			Port:       12345,
			TargetPort: intstr.FromInt(12345),
		},
	}

	objService.Selector = map[string]string{
		"app": "azizi-myfirst",
	}

	configContext := objDeployment.GetKubeConfig()

	objDeployment.ConnectToCluster(configContext)

	objDeployment.AppName = "test"

	log.Print(objDeployment.AppName)
	log.Print(objService.AppName)

	log.Println("test")
}

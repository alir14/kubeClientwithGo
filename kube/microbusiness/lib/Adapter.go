package adapter

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

//TestPackage TestPackage
func TestPackage(test string) {
	fmt.Println(test)
}

type KubeAdapter interface {
	GetKubeConfig(isUsingKind bool) string
	ConnectToCluster(configContext *rest.Config) kubernetes.Clientset
	CreateDeployment(clientSet *kubernetes.Clientset)
	UpdateDeployment(clientSet *kubernetes.Clientset)
	DeleteDeployment(clientSet *kubernetes.Clientset)
	CreateService(clientSet *kubernetes.Clientset)
	UpdateService(clientSet *kubernetes.Clientset)
	DeleteService(clientSet *kubernetes.Clientset)
	GetPods(clientSet *kubernetes.Clientset) []string
	GetPod(clientSet *kubernetes.Clientset, namespace string, podName string) string
}

type EdgeClusterDeploymentDetail struct {
	Name           string
	DomainName     string
	IPAddress      string
	Replicas       int32
	ContainerName  string
	ContainerImage string
}

type EdgeClusterServiceDetail struct {
	Name           string
	DomainName     string
	IPAddress      string
	Replicas       int32
	ContainerName  string
	ContainerImage string
}

package edgeclusterlib

import (
	"flag"
	microbusiness "kube/microbusiness"
	"log"
	"path/filepath"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

//EdgeClusterServiceDetail micro business adapter for service
type EdgeClusterServiceDetail struct {
	Name           string
	DomainName     string
	IPAddress      string
	Replicas       int32
	ContainerName  string
	ContainerImage string
}

//GetKubeConfig getting kube configuration from os
func (edge EdgeClusterServiceDetail) GetKubeConfig(isUsingKind bool) *rest.Config {
	var kubeconfig string

	//get hoem directory path
	homeDir := microbusiness.GetHomeDirectoryPath()
	log.Print(homeDir)

	var configName string
	if isUsingKind {
		configName = "kind-config-devEnv"
	} else {
		configName = "config"
	}

	if homeDir != "" {
		flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(homeDir, ".kube", configName), "(optional) path to config file")
	} else {
		flag.StringVar(&kubeconfig, "kubeconfig", "", "path to kube config file")
	}

	log.Println("building config ...")

	configContext, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

	microbusiness.HandleError(err)

	return configContext
}

//ConnectToCluster connecting to cluster
func (edge EdgeClusterServiceDetail) ConnectToCluster(configContext *rest.Config) kubernetes.Clientset {
	log.Print("get Client")

	client, err := kubernetes.NewForConfig(configContext)

	microbusiness.HandleError(err)

	return *client
}

//Create service
func (edge EdgeClusterServiceDetail) Create(clientSet *kubernetes.Clientset) {
	log.Println("call Create from service")

	serviceDeployment := clientSet.CoreV1().Services(apiv1.NamespaceDefault)

	serviceConfig := edge.populateDeploymentConfigValue()

	log.Println("creating ....")

	result, err := serviceDeployment.Create(serviceConfig)

	microbusiness.HandleError(err)

	log.Printf("created service %q /n", result.GetObjectMeta().GetName())
}

//Update service
func (edge EdgeClusterServiceDetail) Update(clientSet *kubernetes.Clientset) {
	log.Println("call Update from service")
}

//Delete service
func (edge EdgeClusterServiceDetail) Delete(clientSet *kubernetes.Clientset) {
	log.Println("call Delete from service")
}

func (edge EdgeClusterServiceDetail) populateDeploymentConfigValue() *apiv1.Service {
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "demo-deploymentservice",
			Namespace: apiv1.NamespaceDefault,
			Labels: map[string]string{
				"k8s-app": "kube-controller-manager",
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports:     nil,
			Selector:  nil,
			ClusterIP: "",
		},
	}

	return service
}

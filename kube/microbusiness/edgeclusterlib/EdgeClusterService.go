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
	"k8s.io/client-go/util/retry"
)

//EdgeClusterServiceDetail micro business adapter for service
type EdgeClusterServiceDetail struct {
	Metaobject     microbusiness.DeploymentMetaData
	AppName        string
	IPAddress      string
	Ports          int32
	Replicas       int32
	ContainerName  string
	ContainerImage string
	ConfigName     string
	Selector       string
}

//GetKubeConfig getting kube configuration from os
func (edge EdgeClusterServiceDetail) GetKubeConfig() *rest.Config {
	var kubeconfig string

	//get hoem directory path
	homeDir := microbusiness.GetHomeDirectoryPath()
	log.Print(homeDir)

	if homeDir != "" {
		flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(homeDir, ".kube", edge.ConfigName), "(optional) path to config file")
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
func (edge EdgeClusterServiceDetail) UpdateWithRetry(clientSet *kubernetes.Clientset) {
	log.Println("call Update from service")

	updateClient := clientSet.AppsV1().Deployments(apiv1.NamespaceDefault)

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {

		result, getErr := updateClient.Get(edge.Metaobject.Name, metav1.GetOptions{})
		if getErr != nil {
			log.Println("Failed to get the deployment for update..")
		}

		//Do what need to be updated
		//Todo complete the whole necessary fileds
		result.Spec.Replicas = microbusiness.Int32Ptr(edge.Replicas)
		result.Spec.Template.Spec.Containers[0].Image = edge.ContainerImage

		_, updateErr := updateClient.Update(result)

		return updateErr
	})

	microbusiness.HandleError(retryErr)
	log.Println("Update complete ...")
}

//Delete service
func (edge EdgeClusterServiceDetail) Delete(clientSet *kubernetes.Clientset) {
	log.Println("call Delete from service")

	deleteClient := clientSet.AppsV1().Deployments(apiv1.NamespaceDefault)
	deletePolicy := metav1.DeletePropagationForeground

	err := deleteClient.Delete(edge.Metaobject.Name, &metav1.DeleteOptions{
		DeletePropagation: &deletePolicy,
	})

	microbusiness.HandleError(err)
}

func (edge EdgeClusterServiceDetail) populateDeploymentConfigValue() *apiv1.Service {
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      edge.Metaobject.Name,
			Namespace: apiv1.NamespaceDefault,
			Labels: map[string]string{
				"k8s-app": "kube-controller-manager",
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports:     edge.Ports,
			Selector:  edge.Selector,
			ClusterIP: edge.IPAddress,
		},
	}

	return service
}

package edgeclusterlib

import (
	"flag"
	microbusiness "kube/microbusiness"
	"log"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
)

//EdgeClusterDeploymentDetail microbusiness adapter for deployment
type EdgeClusterDeploymentDetail struct {
	Name           string
	DomainName     string
	IPAddress      string
	Replicas       int32
	ContainerName  string
	ContainerImage string
}

//GetKubeConfig getting kube configuration from os
func (edge EdgeClusterDeploymentDetail) GetKubeConfig(isUsingKind bool) *rest.Config {
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
func (edge EdgeClusterDeploymentDetail) ConnectToCluster(configContext *rest.Config) kubernetes.Clientset {
	log.Print("get Client")

	client, err := kubernetes.NewForConfig(configContext)

	microbusiness.HandleError(err)

	return *client
}

//Create deployment
func (edge EdgeClusterDeploymentDetail) Create(clientSet *kubernetes.Clientset) {
	log.Println("call Create from deployment")
	deploymentClient := clientSet.AppsV1().Deployments(apiv1.NamespaceDefault)

	deploymentConfig := edge.populateDeploymentConfigValue()

	log.Println("creating ...")

	result, err := deploymentClient.Create(deploymentConfig)

	microbusiness.HandleError(err)

	log.Printf("created deployment %q /n", result.GetObjectMeta().GetName())
}

//Update deployment
func (edge EdgeClusterDeploymentDetail) UpdateWithRetry(clientSet *kubernetes.Clientset) {
	log.Println("call Update from deployment")

	updateClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// result, geterr := updateClient.Get(edge.)
	})
}

//Delete deployment
func (edge EdgeClusterDeploymentDetail) Delete(clientSet *kubernetes.Clientset) {
	log.Println("call Delete from deployment")
}

//PopulateDeploymentConfigValue create spec object for deployment
func (edge EdgeClusterDeploymentDetail) populateDeploymentConfigValue() *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: microbusiness.Int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	return deployment
}

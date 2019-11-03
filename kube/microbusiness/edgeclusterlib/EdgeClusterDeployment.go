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
	Metaobject     microbusiness.MetaData
	AppName        string
	IPAddress      string
	Replicas       int32
	ContainerName  string
	ContainerImage string
	ConfigName     string
}

//GetKubeConfig getting kube configuration from os
func (edge EdgeClusterDeploymentDetail) GetKubeConfig() *rest.Config {
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
func (edge EdgeClusterDeploymentDetail) ConnectToCluster(configContext *rest.Config) kubernetes.Clientset {
	log.Print("get Client")

	client, err := kubernetes.NewForConfig(configContext)

	microbusiness.HandleError(err)

	return *client
}

//Create deployment
func (edge EdgeClusterDeploymentDetail) Create(clientSet *kubernetes.Clientset) {
	log.Println("call Create from deployment")
	deploymentClient := clientSet.AppsV1().Deployments(edge.Metaobject.NameSpace)

	deploymentConfig := edge.populateDeploymentConfigValue()

	log.Println("creating ...")

	result, err := deploymentClient.Create(deploymentConfig)

	microbusiness.HandleError(err)

	log.Printf("created deployment %q /n", result.GetObjectMeta().GetName())
}

//UpdateWithRetry deployment
func (edge EdgeClusterDeploymentDetail) UpdateWithRetry(clientSet *kubernetes.Clientset) {
	log.Println("call Update from deployment")

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

//Delete deployment
func (edge EdgeClusterDeploymentDetail) Delete(clientSet *kubernetes.Clientset) {
	log.Println("call Delete from deployment")

	deleteClient := clientSet.AppsV1().Deployments(apiv1.NamespaceDefault)
	deletePolicy := metav1.DeletePropagationForeground

	err := deleteClient.Delete(edge.Metaobject.Name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})

	microbusiness.HandleError(err)
}

//PopulateDeploymentConfigValue create spec object for deployment
func (edge EdgeClusterDeploymentDetail) populateDeploymentConfigValue() *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: edge.Metaobject.Name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: microbusiness.Int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": edge.AppName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": edge.AppName,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  edge.ContainerName,
							Image: edge.ContainerImage,
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

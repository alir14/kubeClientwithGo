package edgeclusterlib

import (
	microbusiness "kube/microbusiness"
	"log"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// KubeAdapter a struct for grtting metadata information from kubernetes
type KubeAdapter struct {
	NameSpace string
	PodName   string
}

//GetPods getting available pods in Kubernetes
func (adapter KubeAdapter) GetPods(clientSet *kubernetes.Clientset) *v1.PodList {
	log.Print("get All Pods")

	pods, err := clientSet.CoreV1().Pods(adapter.NameSpace).List(metav1.ListOptions{})

	microbusiness.HandleError(err)

	log.Printf(" number of pods is %s /n", len(pods.Items))

	return pods
}

//GetPod getting a specific pod by name in Kubernetes
func (adapter KubeAdapter) GetPod(clientSet *kubernetes.Clientset) *v1.Pod {
	log.Printf(" get pod %s /n", adapter.PodName)

	pod, err := clientSet.CoreV1().Pods(adapter.NameSpace).Get(adapter.PodName, metav1.GetOptions{})

	statusError, isStatus := err.(*errors.StatusError)

	if errors.IsNotFound(err) {
		log.Printf("Pod %s in namespace %s not found /n", adapter.PodName, adapter.NameSpace)
	} else if err != nil && isStatus {
		log.Printf("Error getting pod %s in namespace %s: %v\n ",
			adapter.PodName, adapter.NameSpace, statusError.ErrStatus.Message)
	} else if err != nil {
		microbusiness.HandleError(err)
	}

	return pod
}

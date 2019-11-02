package edgeclusterinterfaces

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// KubeCRUDAdapter microbusiness Kubernetes interface for create,update and delete
type KubeCRUDAdapter interface {
	Create(clientSet *kubernetes.Clientset)
	UpdateWithRetry(clientSet *kubernetes.Clientset)
	Replace(clientSet *kubernetes.Clientset)
	Delete(clientSet *kubernetes.Clientset)
}

// KuebConnectionAdapter microbusiness Kubernetes interface for connecting to Kubernetes
type KuebConnectionAdapter interface {
	GetKubeConfig(isUsingKind bool) *rest.Config
	ConnectToCluster(configContext *rest.Config) kubernetes.Clientset
}

// KubeMonitor microbusiness Kubernetes interface for monitoring
type KubeMonitor interface {
	GetPods(clientSet *kubernetes.Clientset) *v1.PodList
	GetPod(clientSet *kubernetes.Clientset) *v1.Pod
}

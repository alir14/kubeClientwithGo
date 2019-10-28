package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var isUsingKind bool = true
	var kubeConfig string

	//get kubeconfig
	kubeConfig = getKubeConfig(isUsingKind)

	//use the current context in kubeconfig
	configContext, err := clientcmd.BuildConfigFromFlags("", kubeConfig)

	handleError(err)

	// create the clientset
	clientSet := getClientSetInstance(configContext)

	for {
		showMethePods(&clientSet)

		time.Sleep(10 * time.Second)
	}
}

func getKubeConfig(usingKind bool) string {
	var kubeconfig string

	//get hoem directory path
	homeDir := getHomeDirectoryPath()
	log.Print(homeDir)

	var configName string
	if usingKind {
		configName = "kind-config-devEnv"
	} else {
		configName = "config"
	}

	if homeDir != "" {
		flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(homeDir, ".kube", configName), "(optional) path to config file")
	} else {
		flag.StringVar(&kubeconfig, "kubeconfig", "", "path to kube config file")
	}
	return kubeconfig
}

func getHomeDirectoryPath() string {
	homePath, err := os.UserHomeDir()
	handleError(err)
	log.Print(homePath)
	if homePath != "" {
		log.Print("linux mode")
		return homePath
	}
	log.Print("windows mode")
	return os.Getenv("USERPROFILE")
}

func getClientSetInstance(configContext *rest.Config) kubernetes.Clientset {
	log.Print("get Client")

	client, err := kubernetes.NewForConfig(configContext)

	handleError(err)

	return *client
}

func showMethePods(clientSet *kubernetes.Clientset) {
	log.Print("get All Pods")
	pods, err := clientSet.CoreV1().Pods("").List(metav1.ListOptions{})

	handleError(err)

	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	log.Print(pods.Items[1].Spec.DNSConfig)
}

func getSpecificPod(clientSet *kubernetes.Clientset, namespace string, podName string) {
	log.Printf("get Pod %s \n", podName)
	_, err := clientSet.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})

	if errors.IsNotFound(err) {
		fmt.Printf("Pod %s in namespace %s not found\n", podName, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting pod %s in namespace %s: %v\n",
			podName, namespace, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found pod %s in namespace %s\n", podName, namespace)
		handleError(err)
	}
}

func handleError(err error) {
	if err != nil {
		log.Print(err)
		panic(err.Error())
	}
}

func getKindConfigByCommand() string {

	cmd := exec.Command("kind", "get", "kubeconfig-path", "--name=devEnv ")
	output, err := cmd.CombinedOutput()
	handleError(err)
	result := string(output)
	return result
}

package main

import (
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

func k8sClient() (*kubernetes.Clientset, error) {
	var config *rest.Config
	kuBeFle := os.Getenv("KUBECONF")
	if len(kuBeFle) == 0 {
		// creates the in-cluster config
		log.Println("[INFO] KUBECONF is not set, use default.")
		config, err = rest.InClusterConfig()
	} else {
		kuBeConf, err := ioutil.ReadFile(kuBeFle)
		if err != nil {
			return nil, err
		}
		config, err = clientcmd.RESTConfigFromKubeConfig(kuBeConf)
	}
	if err != nil {
		return nil, err
	}
    return kubernetes.NewForConfig(config)
}

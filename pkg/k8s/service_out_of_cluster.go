package k8s

import (
	"errors"
	"path/filepath"

	"git.isi.nc/go/dtb-tool/pkg/env"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8SOutSvc struct {
	clientset *kubernetes.Clientset
}

func getKubeConfig() (string, error) {
	// get env KUBE_CONFIG or default homedir kube config

	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	//📢 voluntarily named KUBE_CONFIG and not KUBECONFIG
	kubeconfig = env.GetEnvOrDefault("KUBE_CONFIG", kubeconfig)

	if kubeconfig == "" {
		return "", errors.New("kubeconfig not found, nor provided by environment variable")
	}

	return kubeconfig, nil
}

// Instanciate a new kubernetes out of cluster client
// using kubeconfig from env var KUBE_CONFIG (expect file path)
// or trying to find $HOME/.kube/config
// returns a new service with a clientset to interact with k8s
func NewK8SOutSvc() (*K8SOutSvc, error) {
	k := &K8SOutSvc{}

	kubeconfig, err := getKubeConfig()
	if err != nil {
		return nil, err
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	k.clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return k, nil
}

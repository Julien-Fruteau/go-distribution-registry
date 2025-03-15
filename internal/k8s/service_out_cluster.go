package k8s

import (
	"context"
	"errors"
	"path/filepath"
	"slices"
	"time"

	"github.com/julien-fruteau/go/distctl/internal/env"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8SOutSvc struct {
	ctx       context.Context
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
		return "", errors.New("kubeconfig path not found, nor provided by environment variable KUBE_CONFIG")
	}

	return kubeconfig, nil
}

// Instanciate a new kubernetes out of cluster client
// using kubeconfig from env var KUBE_CONFIG (expect file path)
// or trying to find $HOME/.kube/config
// returns a new service with a clientset to interact with k8s
func NewK8SOutSvc(ctx context.Context) (*K8SOutSvc, error) {
	k := &K8SOutSvc{}

	// 📢 TODO understand context, and
	// 📢 use parent context and allow some SIG TERM ctrl-c and so on
	// 📢 apply this to registry service too, so that all
	// 📢 can be cancelled
	if ctx != nil {
		k.ctx = ctx
	} else {
		k.ctx = context.Background()
	}

	// k.ctx = context.Background()
	// k.ctx, cancel = context.WithTimeout(context.Background(), time.Second*30)

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

// RETRIEVE CLUSTER IMAGES IN ORDER TO NOT DELETE THOSE IMAGES

// func (k *K8SOutSvc) GetClusterImages() ([]Image, error) {
func (k *K8SOutSvc) GetClusterImages() ([]string, error) {
	var images []string
	// var images []Image

	//  func slowOperationWithTimeout(ctx context.Context) (Result, error) {
	// 	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	// 	defer cancel()  // releases resources if slowOperation completes before timeout elapses
	// 	return slowOperation(ctx)
	// }

	ctx, cancel := context.WithTimeout(k.ctx, DEFAULT_TIMEOUT*time.Second)
	defer cancel() // releases resources if slowOperation completes before timeout elapses

	pods, err := k.clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return images, err
	}

	// select {
	// case <-k.ctx.Done():
	// 	return images, k.ctx.Err()
	// case res := <-data:
	// 	return res, nil
	// }

	return getImagesV1(pods), nil
}

func getImagesV1(pods *v1.PodList) []string {
	var images []string
	for _, pod := range pods.Items {
		for _, c := range pod.Spec.Containers {
			if c.Image != "" && !slices.Contains(images, c.Image) {
				images = append(images, c.Image)
			}
		}
	}
	return images
}

// func getImagesV2(pods *v1.PodList) []string {
// 	// var images []string
//   images := slices.Map(pods.Items, func(pod v1.Pod) string {
//     return slices.Map(pod.Spec.Containers, func(c v1.Container) string { return c.Image })
//   })
// 	return images
// }

// 📢 READ https://go.dev/blog/pipelines 🛫

// Stream generates values with DoSomething and sends them to out
// until DoSomething returns an error or ctx.Done is closed.
// func Stream(ctx context.Context, out chan<- Value) error {
// 	for {
// 		v, err := DoSomething(ctx)
// 		if err != nil {
// 			return err
// 		}
// 		select {
// 		case <-ctx.Done():
// 			return ctx.Err()
// 		case out <- v:
// 		}
// 	}
// }

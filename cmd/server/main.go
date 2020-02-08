package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/storyscript/scheduler"
	"github.com/storyscript/scheduler/fakes"
	"github.com/storyscript/scheduler/http"
	"github.com/storyscript/scheduler/kube"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func main() {
	kubeClient := &kube.Client{
		ClientSet: newKubeClientSet(),
	}

	scheduler := &scheduler.Scheduler{
		Publisher: &fakes.Publisher{},
		Deployer:  kubeClient,
	}

	_, err := scheduler.Start()
	if err != nil {
		panic(err)
	}

	server := http.Server{}

	if err := server.Start(); err != nil {
		panic(err)
	}
}

func newKubeClientSet() *kubernetes.Clientset {
	if os.Getenv("CLUSTER") != "" {
		fmt.Println("loading in cluster config")
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(errors.Wrap(err, "failed to load cluster config"))
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(errors.Wrap(err, "failed to create the clientset"))
		}

		return clientset
	}

	fmt.Println("loading out of cluster config")
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	if err != nil {
		panic(errors.Wrap(err, "failed to load config from file"))
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(errors.Wrap(err, "failed to create the clientset"))
	}

	return clientset
}

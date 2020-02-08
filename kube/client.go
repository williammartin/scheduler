package kube

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/avast/retry-go"
	"github.com/storyscript/scheduler"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

type Client struct {
	ClientSet *kubernetes.Clientset
}

func (c *Client) Deploy(story scheduler.Story) error {
	fmt.Fprintln(os.Stdout, "Deploying Story: ", story)
	if err := c.ensureNamespaceExists(story.Name); err != nil {
		return err
	}

	if err := c.ensureRuntimeExists(story.Name); err != nil {
		return err
	}

	if err := c.deployStory(story); err != nil {
		return err
	}

	return nil
}

func (c *Client) ensureNamespaceExists(name string) error {
	namespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	// TODO: This is probably not synchronous and we may need a watch after
	_, err := c.ClientSet.CoreV1().Namespaces().Create(namespace)
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

func (c *Client) ensureRuntimeExists(namespace string) error {
	deploymentsClient := c.ClientSet.AppsV1().Deployments(namespace)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      "runtime",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "runtime",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "runtime",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "runtime",
							Image: "williammartin/fakeruntime",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 3000,
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := deploymentsClient.Create(deployment)
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}

	_, err = c.ClientSet.CoreV1().Services(namespace).Create(&v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      "runtime",
			Labels: map[string]string{
				"app": "runtime",
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{{Port: 3000}},
			Selector: map[string]string{
				"app": "runtime",
			},
		},
	})
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

func (c *Client) deployStory(story scheduler.Story) error {
	// TODO: This only works in cluster yikes

	runtimeURI := fmt.Sprintf("http://runtime.%s.svc.cluster.local:3000", story.Name)
	err := retry.Do(func() error {
		resp, err := http.Get(fmt.Sprintf("%s/healthcheck", runtimeURI))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to healthcheck: expected status code 200 but got %d", resp.StatusCode)
		}

		return nil
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(fmt.Sprintf("%s/api/app/deploy", runtimeURI), "application/json", bytes.NewBufferString(story.Payload))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to deploy: expected status code 200 but got %d", resp.StatusCode)
	}

	return nil
}

func int32Ptr(i int32) *int32 { return &i }

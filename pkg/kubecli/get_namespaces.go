package kubecli

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func (cli *KubeCLI) GetNamespaces(ctx context.Context, opts metav1.ListOptions) (*v1.NamespaceList, error) {
	config, err := cli.factory.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return client.CoreV1().Namespaces().List(ctx, opts)
}

func (cli *KubeCLI) HasNamespacePermission(ctx context.Context) bool {
	_, err := cli.GetNamespaces(ctx, metav1.ListOptions{})
	if err == nil {
		return true
	}

	if errors.IsForbidden(err) {
		return false
	}

	return true
}

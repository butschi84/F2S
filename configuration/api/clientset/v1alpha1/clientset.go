package v1alpha1

import (
	"butschi84/f2s/configuration/api/types/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime/serializer"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type V1Alpha1Interface interface {
	Functions(namespace string) FunctionInterface
}

type V1Alpha1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*V1Alpha1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1alpha1.GroupName, Version: v1alpha1.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &V1Alpha1Client{restClient: client}, nil
}

func (c *V1Alpha1Client) Functions(namespace string) FunctionInterface {
	return &functionClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

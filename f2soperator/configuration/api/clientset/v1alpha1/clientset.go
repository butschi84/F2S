package v1alpha1

import (
	"butschi84/f2s/configuration/api/types/v1alpha1"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/client-go/rest"
)

type V1Alpha1Interface interface {
	Functions(namespace string) FunctionInterface
}

type V1Alpha1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config, scheme *runtime.Scheme) (*V1Alpha1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1alpha1.GroupName, Version: v1alpha1.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme)
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	logging.Info("Rest Client Group Version:", fmt.Sprintf("%s", config.GroupVersion))

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

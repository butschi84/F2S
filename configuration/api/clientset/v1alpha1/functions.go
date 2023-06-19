package v1alpha1

import (
	"butschi84/f2s/configuration/api/types/v1alpha1"
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type FunctionInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.FunctionList, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*v1alpha1.Function, error)
	Create(*v1alpha1.Function) (*v1alpha1.Function, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

type functionClient struct {
	restClient rest.Interface
	ns         string
}

func (c *functionClient) List(opts metav1.ListOptions) (*v1alpha1.FunctionList, error) {
	result := v1alpha1.FunctionList{}
	logging.Println("get function list for namespace", c.ns)
	ctx := context.TODO()
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("functions").
		//VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)
	logging.Println(err)
	return &result, err
}

func (c *functionClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1alpha1.Function, error) {
	result := v1alpha1.Function{}
	req := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("functions").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec)

	err := req.Do(ctx).Into(&result)
	return &result, err
}

func (c *functionClient) Create(project *v1alpha1.Function) (*v1alpha1.Function, error) {
	result := v1alpha1.Function{}
	ctx := context.TODO()
	req := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("functions").
		Body(project)

	err := req.Do(ctx).Into(&result)
	return &result, err
}

func (c *functionClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("functions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}

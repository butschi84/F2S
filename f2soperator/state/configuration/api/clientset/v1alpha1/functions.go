package v1alpha1

import (
	"butschi84/f2s/state/configuration/api/types/v1alpha1"
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type FunctionInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.FunctionList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.Function, error)
	Create(*v1alpha1.Function) (*v1alpha1.Function, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Delete(uid string, opts metav1.DeleteOptions) error
}

type functionClient struct {
	restClient rest.Interface
	ns         string
}

func (c *functionClient) List(opts metav1.ListOptions) (*v1alpha1.FunctionList, error) {
	result := v1alpha1.FunctionList{}
	logging.Info(fmt.Sprintf("get function list for namespace %s", c.ns))
	ctx := context.TODO()
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("functions").
		//VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)
	if err != nil {
		logging.Error(fmt.Errorf("[List] could not get function list for namespace %s: %s", c.ns, err.Error()))
	}
	return &result, err
}

func (c *functionClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.Function, error) {
	result := v1alpha1.Function{}
	ctx := context.TODO()
	req := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("functions").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec)

	err := req.Do(ctx).Into(&result)
	logging.Error(fmt.Errorf("[Get] could no get function definitions: %s", err.Error()))
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
	logging.Error(fmt.Errorf("[Create] could not create function definition: %s", err.Error()))
	return &result, err
}

func (c *functionClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	ctx := context.TODO()
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("functions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}

func (c *functionClient) Delete(uid string, opts metav1.DeleteOptions) error {
	ctx := context.TODO()

	objUID := types.UID(uid)
	opts.Preconditions = &metav1.Preconditions{
		UID: &objUID,
	}

	functions, _ := c.List(metav1.ListOptions{})
	for _, f := range functions.Items {
		if string(f.UID) == uid {
			fmt.Println("functon found")
			err := c.restClient.
				Delete().
				Namespace(c.ns).
				Resource("functions").
				Name(f.Name).
				Do(ctx).
				Error()
			logging.Error(fmt.Errorf("[Delete] could not delete function definition: %s", err.Error()))
			return err
		}
	}

	logging.Error(fmt.Errorf("[Delete] function with uid %s could not be found", uid))
	return fmt.Errorf("function could not be found")
}

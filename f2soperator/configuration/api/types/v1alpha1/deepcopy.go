package v1alpha1

import "k8s.io/apimachinery/pkg/runtime"

// DeepCopyInto copies all properties of this object into another object of the
// same type that is provided as a pointer.
func (in *Function) DeepCopyInto(out *Function) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = FunctionSpec{
		Endpoint: in.Spec.Endpoint,
		Method:   in.Spec.Method,
	}
	out.Target = FunctionTarget{
		ContainerImage: in.Target.ContainerImage,
		Endpoint:       in.Target.Endpoint,
		Port:           in.Target.Port,
		MinReplicas:    in.Target.MinReplicas,
		MaxReplicas:    in.Target.MaxReplicas,
	}
}

// DeepCopyObject returns a generically typed copy of an object
func (in *Function) DeepCopyObject() runtime.Object {
	out := Function{}
	in.DeepCopyInto(&out)

	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *FunctionList) DeepCopyObject() runtime.Object {
	out := FunctionList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]Function, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}

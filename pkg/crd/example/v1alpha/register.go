package v1alpha

import (
	crdexample "github.com/jinghzhu/KubernetesCRD/pkg/crd/example"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	Kind string = "Example"
	// GroupVersion is the version.
	GroupVersion string = "v1alpha"
	// Plural is the Plural for Example.
	Plural string = "examples"
	// Singular is the singular for Example.
	Singular string = "example"
	// CRDName is the CRD name for Example.
	CRDName string = Plural + "." + crdexample.GroupName
)

var (
	// SchemeGroupVersion is the group version used to register these objects.
	SchemeGroupVersion = schema.GroupVersion{
		Group:   crdexample.GroupName,
		Version: GroupVersion,
	}
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// addKnownTypes adds the set of types defined in this package to the supplied scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Example{},
		&ExampleList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)

	return nil
}

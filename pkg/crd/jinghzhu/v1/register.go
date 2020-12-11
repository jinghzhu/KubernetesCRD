package v1

import (
	crdjinghzhu "github.com/jinghzhu/KubernetesCRD/pkg/crd/jinghzhu"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// Kind is normally the CamelCased singular type. The resource manifest uses this.
	Kind string = "Jinghzhu"
	// GroupVersion is the version.
	GroupVersion string = "v1"
	// Plural is the plural name used in /apis/<group>/<version>/<plural>
	Plural string = "jinghzhus"
	// Singular is used as an alias on kubectl for display.
	Singular string = "jinghzhu"
	// CRDName is the CRD name for Jinghzhu.
	CRDName string = Plural + "." + crdjinghzhu.GroupName
	// ShortName is the short alias for the CRD.
	ShortName string = "jh"
)

var (
	// SchemeGroupVersion is the group version used to register these objects.
	SchemeGroupVersion = schema.GroupVersion{
		Group:   crdjinghzhu.GroupName,
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
		&Jinghzhu{},
		&JinghzhuList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)

	return nil
}

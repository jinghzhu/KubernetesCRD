package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	TestResourcePlural string = "tests"
	// GroupName is the group name used in this package.
	GroupName      string = "test0.io"
	TestCRDName    string = TestResourcePlural + "." + GroupName
	version        string = "v1"
	StateCreated   string = "Created"
	StateProcessed string = "Processed"
)

var (
	// SchemeGroupVersion is the group version used to register these objects.
	SchemeGroupVersion = schema.GroupVersion{
		Group:   GroupName,
		Version: version,
	}
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// Test is the CRD.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Test struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              TestSpec   `json:"spec"`
	Status            TestStatus `json:"status,omitempty"`
}

type TestSpec struct {
	Foo string `json:"foo"`
	Bar bool   `json:"bar"`
}

type TestStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Test `json:"items"`
}

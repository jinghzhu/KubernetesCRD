package v1alpha

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Example is the CRD. Use this command to generate deepcopy for it:
// ./k8s.io/code-generator/generate-groups.sh all github.com/jinghzhu/KubernetesCRD/pkg/crd/example/v1alpha/apis github.com/jinghzhu/KubernetesCRD/pkg/crd "example:v1alpha"
// For more details of code-generator, please visit https://github.com/kubernetes/code-generator
// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Example struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata"`
	// Specification of the desired behavior of Example.
	Spec ExampleSpec `json:"spec"`
	// Observed status of Example.
	Status ExampleStatus `json:"status"`
}

// ExampleSpec is a desired state description of Example.
// +k8s:deepcopy-gen=true
type ExampleSpec struct {
	Foo string `json:"foo"`
	Bar bool   `json:"bar"`
}

// ExampleStatus describes the lifecycle status of Example.
// +k8s:deepcopy-gen=true
type ExampleStatus struct {
	State   string `json:"state"`
	Message string `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExampleList is the list of Examples.
type ExampleList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	metav1.ListMeta `json:"metadata"`
	// List of Examples.
	Items []Example `json:"items"`
}

package v1alpha

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Jinghzhu is the CRD. Use this command to generate deepcopy for it:
// ./k8s.io/code-generator/generate-groups.sh all github.com/jinghzhu/KubernetesCRD/pkg/crd/jinghzhu/v1alpha/apis github.com/jinghzhu/KubernetesCRD/pkg/crd "jinghzhu:v1alpha"
// For more details of code-generator, please visit https://github.com/kubernetes/code-generator
// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Jinghzhu struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata"`
	// Specification of the desired behavior of Jinghzhu.
	Spec JinghzhuSpec `json:"spec"`
	// Observed status of Jinghzhu.
	Status JinghzhuStatus `json:"status"`
}

// JinghzhuSpec is a desired state description of Jinghzhu.
// +k8s:deepcopy-gen=true
type JinghzhuSpec struct {
	Foo string `json:"foo"`
	Bar bool   `json:"bar"`
}

// JinghzhuStatus describes the lifecycle status of Jinghzhu.
// +k8s:deepcopy-gen=true
type JinghzhuStatus struct {
	State   string `json:"state"`
	Message string `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// JinghzhuList is the list of Jinghzhus.
type JinghzhuList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	metav1.ListMeta `json:"metadata"`
	// List of Jinghzhus.
	Items []Jinghzhu `json:"items"`
}

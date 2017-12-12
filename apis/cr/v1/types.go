package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const ExampleResourcePlural = "examples"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Example struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ExampleSpec   `json:"spec"`
	Status            ExampleStatus `json:"status,omitempty"`
}

type ExampleSpec struct {
	Foo string `json:"foo"`
	Bar bool   `json:"bar"`
}

type ExampleStatus struct {
	State   ExampleState `json:"state,omitempty"`
	Message string       `json:"message,omitempty"`
}

type ExampleState string

const (
	ExampleStateCreated   ExampleState = "Created"
	ExampleStateProcessed ExampleState = "Processed"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ExampleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Example `json:"items"`
}

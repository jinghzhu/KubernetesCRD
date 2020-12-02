package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// StatePending means CRD instance is created; Pod info has been updated into CRD instance;
	// Pod has been accepted by the system, but one or more of the containers has not been started.
	StatePending string = "Pending"
	// StateRunning means Pod has been bound to a node and all of the containers have been started.
	StateRunning string = "Running"
	// StateSucceeded means that all containers in the Pod have voluntarily terminated with a container
	// exit code of 0, and the system is not going to restart any of these containers.
	StateSucceeded string = "Succeeded"
	// StateFailed means that all containers in the Pod have terminated, and at least one container has
	// terminated in a failure (exited with a non-zero exit code or was stopped by the system).
	StateFailed string = "Failed"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=jinghzhu

// Jinghzhu is the CRD. Use this command to generate deepcopy for it:
// ./k8s.io/code-generator/generate-groups.sh all github.com/jinghzhu/KubernetesCRD/pkg/crd/jinghzhu/v1/apis github.com/jinghzhu/KubernetesCRD/pkg/crd "jinghzhu:v1"
// For more details of code-generator, please visit https://github.com/kubernetes/code-generator
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
type JinghzhuSpec struct {
	// Desired is the desired Pod number.
	Desired string `json:"desired"`
	// Current is the number of Pod currently running.
	Current bool `json:"current"`
	// PodList is the name list of current Pods.
	PodList []string `json:"podList"`
}

// JinghzhuStatus describes the lifecycle status of Jinghzhu.
type JinghzhuStatus struct {
	State   string `json:"state"`
	Message string `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=jinghzhu

// JinghzhuList is the list of Jinghzhus.
type JinghzhuList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	metav1.ListMeta `json:"metadata"`
	// List of Jinghzhus.
	Items []Jinghzhu `json:"items"`
}

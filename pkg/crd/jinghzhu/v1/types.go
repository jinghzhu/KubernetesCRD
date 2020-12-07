package v1

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	Desired int `json:"desired"`
	// Current is the number of Pod currently running.
	Current int `json:"current"`
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

func (j *Jinghzhu) String() string {
	return fmt.Sprintf(
		"\tName = %s\n\tResource Version = %s\n\tDesired = %d\n\tCurrent = %d\n\tPodList = %s\n\tState = %s\n\tMessage = %s\n\t",
		j.GetName(),
		j.GetResourceVersion(),
		j.Spec.Desired,
		j.Spec.Current,
		strings.Join(j.Spec.PodList, ", "),
		j.Status.State,
		j.Status.Message,
	)
}

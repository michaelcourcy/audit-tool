package v1alpha1

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type BackupActionSpec struct {
	IgnoreExceptions bool                `json:"ignoreExceptions"`
	ScheduledTime    time.Time           `json:"scheduledTime"`
	Subject          BackupActionSubject `json:"subject"`
}

type BackupActionSubject struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type BackupActionStatus struct {
	EndTime  time.Time `json:"endTime"`
	Progress int64     `json:"progress"`
	State    string    `json:"state"`
}

type BackupAction struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackupActionSpec   `json:"spec"`
	Status BackupActionStatus `json:"status"`
}

type BackupActionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []BackupAction `json:"items"`
}

// DeepCopyInto copies all properties of this object into another object of the
// same type that is provided as a pointer.
func (in *BackupAction) DeepCopyInto(out *BackupAction) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Status = BackupActionStatus{
		Progress: in.Status.Progress,
		EndTime:  in.Status.EndTime,
		State:    in.Status.State,
	}
	out.Spec = BackupActionSpec{
		ScheduledTime:    in.Spec.ScheduledTime,
		IgnoreExceptions: in.Spec.IgnoreExceptions,
		Subject: BackupActionSubject{
			Name:      in.Spec.Subject.Name,
			Namespace: in.Spec.Subject.Namespace,
		},
	}
}

// DeepCopyObject returns a generically typed copy of an object
func (in *BackupAction) DeepCopyObject() runtime.Object {
	out := BackupAction{}
	in.DeepCopyInto(&out)

	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *BackupActionList) DeepCopyObject() runtime.Object {
	out := BackupActionList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]BackupAction, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}

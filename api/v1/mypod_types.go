/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MyPodSpec defines the desired state of MyPod
type MyPodSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	PodAnnotations map[string]string `json:"podAnnotations,omitempty"`
	PodLabels      map[string]string `json:"podLabels,omitempty"`
	PodSpec        v1.PodSpec        `json:"podSpec"`
}

// MyPodStatus defines the observed state of MyPod
type MyPodStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	PodPhase         string `json:"podPhase,omitempty"`
	PodIp            string `json:"podIp,omitEmpty"`
	NodeIp           string `json:"nodeIp,omitEmpty"`
	CreatedTimestamp int64  `json:"createdTimestamp,omitEmpty"`
}

//+kubebuilder:printcolumn:JSONPath=".status.phase",name=Phase,type=string
//+kubebuilder:printcolumn:JSONPath=".status.podIp",name=PodIp,type=string

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MyPod is the Schema for the mypods API
type MyPod struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MyPodSpec   `json:"spec,omitempty"`
	Status MyPodStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MyPodList contains a list of MyPod
type MyPodList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MyPod `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MyPod{}, &MyPodList{})
}

/*
Copyright 2021.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LenSpec defines the desired state of Len
type LenSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Len. Edit len_types.go to remove/update
	UpdateInterval int64 `json:"updateInterval,omitempty"`
}

// LenStatus defines the observed state of Len
type LenStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	NodeName    string       `json:"nodeName,omitempty"`
	LatencyList LatencyList  `json:"latencyList,omitempty"`
	UpdateTime  *metav1.Time `json:"updateTime,omitempty"`
}

type LatencyList []Latency

type Latency struct {
	NodeName string `json:"nodeName,omitempty"`
	Latency  int64  `json:"latency,omitempty"`
	// Bandwidth uint   `json:"bandwidth,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster

// Len is the Schema for the lens API
type Len struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              LenSpec   `json:"spec,omitempty"`
	Status            LenStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LenList contains a list of Len
type LenList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Len `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Len{}, &LenList{})
}

/*
Copyright 2021 Stenic BV.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SqlHostSpec defines the desired state of SqlHost
type SqlHostSpec struct {
	// Engine of the external endpoint (like Mysql)
	Engine EngineType `json:"engine"`

	// Endpoint to manage
	Endpoint SqlHostEndpoint `json:"endpoint"`

	// Credentials to use when connecting to the endpoint
	Credentials SqlCredentials `json:"credentials"`
}

// +kubebuilder:validation:Enum=Mysql
type EngineType string

const (
	// Keep
	EngineTypeMysql EngineType = "Mysql"
)

// SqlHostStatus defines the observed state of SqlHost
type SqlHostStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SqlHost is the Schema for the sqlhosts API
type SqlHost struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SqlHostSpec   `json:"spec,omitempty"`
	Status SqlHostStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SqlHostList contains a list of SqlHost
type SqlHostList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SqlHost `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SqlHost{}, &SqlHostList{})
}

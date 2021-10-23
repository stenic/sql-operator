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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SqlGrantsSpec defines the desired state of SqlGrants
type SqlGrantsSpec struct {
	// Reference to the SqlHost
	HostRef SqlObjectRef `json:"hostRef"`

	// Reference to the SqlUser
	UserRef SqlObjectRef `json:"userRef"`

	// Reference to the SqlUser
	DatabaseRef SqlObjectRef `json:"databaseRef"`

	// List of grants
	Grants []string `json:"grants"`

	// Specifies how to handle deletion of a SqlUser.
	// Valid values are:
	// - "Retain" (default): keeps the external resource when the object is deleted;
	// - "Delete": deletes the external resource when the object is deleted;
	// +optional
	CleanupPolicy CleanupPolicy `json:"cleanupPolicy,omitempty"`
}

// SqlGrantsStatus defines the observed state of SqlGrants
type SqlGrantsStatus struct {
	// Boolean indicating the creation process has started
	Created bool `json:"created"`

	// Timestamp when the user was first created.
	CreationTimestamp *metav1.Time `json:"creationTimestamp,omitempty"`

	// Timestamp when the user was last updated/checked.
	LastModifiedTimestamp *metav1.Time `json:"lastModifiedTimestamp,omitempty"`

	CurrentGrants []string `json:"currentGrants,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Host",type="string",JSONPath=".spec.hostRef.name",description="Name of the host"
//+kubebuilder:printcolumn:name="User",type="string",JSONPath=".spec.userRef.name",description="Name of the user"
//+kubebuilder:printcolumn:name="Database",type="string",JSONPath=".spec.databaseRef.name",description="Name of the database"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// SqlGrants is the Schema for the sqlgrants API
type SqlGrants struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SqlGrantsSpec   `json:"spec,omitempty"`
	Status SqlGrantsStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SqlGrantsList contains a list of SqlGrants
type SqlGrantsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SqlGrants `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SqlGrants{}, &SqlGrantsList{})
}

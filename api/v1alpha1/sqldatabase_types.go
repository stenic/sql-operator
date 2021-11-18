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

// SqlDatabaseSpec defines the desired state of SqlDatabase
type SqlDatabaseSpec struct {
	// Reference to the SqlHost
	HostRef SqlObjectRef `json:"hostRef"`

	// Name of the external database
	DatabaseName string `json:"databaseName"`

	// Specifies how to handle deletion of a SqlUser.
	// Valid values are:
	// - "Retain" (default): keeps the external resource when the object is deleted;
	// - "Delete": deletes the external resource when the object is deleted;
	// +optional
	CleanupPolicy CleanupPolicy `json:"cleanupPolicy,omitempty"`
}

// SqlDatabaseStatus defines the observed state of SqlDatabase
type SqlDatabaseStatus struct {
	// Boolean indicating the creation process has started
	Created bool `json:"created"`

	// Timestamp when the user was first created.
	CreationTimestamp *metav1.Time `json:"creationTimestamp,omitempty"`

	// Timestamp when the user was last updated/checked.
	LastModifiedTimestamp *metav1.Time `json:"lastModifiedTimestamp,omitempty"`

	// String used to identify owership
	OwnerID OwnerID `json:"ownerID,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Host",type="string",JSONPath=".spec.hostRef.name",description="Name of the host"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// SqlDatabase is the Schema for the sqldatabases API
type SqlDatabase struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SqlDatabaseSpec   `json:"spec,omitempty"`
	Status SqlDatabaseStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SqlDatabaseList contains a list of SqlDatabase
type SqlDatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SqlDatabase `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SqlDatabase{}, &SqlDatabaseList{})
}

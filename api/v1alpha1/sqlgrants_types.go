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

// SqlGrantSpec defines the desired state of SqlGrant
type SqlGrantSpec struct {
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

// SqlGrantStatus defines the observed state of SqlGrant
type SqlGrantStatus struct {
	// Boolean indicating the creation process has started
	Created bool `json:"created"`

	// Timestamp when the user was first created.
	CreationTimestamp *metav1.Time `json:"creationTimestamp,omitempty"`

	// Timestamp when the user was last updated/checked.
	LastModifiedTimestamp *metav1.Time `json:"lastModifiedTimestamp,omitempty"`

	CurrentGrants []string `json:"currentGrants,omitempty"`

	// String used to identify owership
	OwnerID OwnerID `json:"ownerID,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="User",type="string",JSONPath=".spec.userRef.name",description="Name of the user"
//+kubebuilder:printcolumn:name="Database",type="string",JSONPath=".spec.databaseRef.name",description="Name of the database"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// SqlGrant is the Schema for the sqlgrant API
type SqlGrant struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SqlGrantSpec   `json:"spec,omitempty"`
	Status SqlGrantStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SqlGrantList contains a list of SqlGrant
type SqlGrantList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SqlGrant `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SqlGrant{}, &SqlGrantList{})
}

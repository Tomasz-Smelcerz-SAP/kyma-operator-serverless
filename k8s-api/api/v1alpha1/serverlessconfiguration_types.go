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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type GithubRepository struct {
	AuthKey string `json:"authKey,omitempty"`
	URL     string `json:"url,omitempty"`
}

// ServerlessConfigurationSpec defines the desired state of ServerlessConfiguration
type ServerlessConfigurationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	FunctionPrefix   string           `json:"commonPrefix,omitempty"`
	GithubRepository GithubRepository `json:"githubRepository,omitempty"`
}

// ServerlessConfigurationStatus defines the observed state of ServerlessConfiguration
type ServerlessConfigurationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ServerlessConfiguration is the Schema for the serverlessconfigurations API
type ServerlessConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServerlessConfigurationSpec   `json:"spec,omitempty"`
	Status ServerlessConfigurationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ServerlessConfigurationList contains a list of ServerlessConfiguration
type ServerlessConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServerlessConfiguration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ServerlessConfiguration{}, &ServerlessConfigurationList{})
}

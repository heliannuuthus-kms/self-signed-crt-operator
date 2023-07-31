/*
Copyright 2023 The cert-manager Authors.

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

// CAIssuerSpec defines the desired state of CAIssuer
type CAIssuerSpec struct {
	// A reference to a Secret in the same namespace as the referent. If the
	// referent is a ClusterCAIssuer, the reference instead refers to the resource
	// with the given name in the configured 'cluster resource namespace', which
	// is set as a flag on the controllers component (and defaults to the
	// namespace that the controllers runs in).
	CASecretName string `json:"CASecretName" default:"idaas-caissuer-secret"`
}

// CAIssuerStatus defines the observed state of CAIssuer
type CAIssuerStatus struct {
	// List of status conditions to indicate the status of a CertificateRequest.
	// Known condition types are `Ready`.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

const ConditionTypeReady = "Ready"

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CAIssuer is the Schema for the issuers API
type CAIssuer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CAIssuerSpec   `json:"spec,omitempty"`
	Status CAIssuerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CAIssuerList contains a list of CAIssuer
type CAIssuerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CAIssuer `json:"items"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ClusterCAIssuer is the Schema for the clusterissuers API
// +kubebuilder:resource:path=clustercaissuers,scope=Cluster
type ClusterCAIssuer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CAIssuerSpec   `json:"spec,omitempty"`
	Status CAIssuerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterCAIssuerList contains a list of ClusterCAIssuer
type ClusterCAIssuerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterCAIssuer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CAIssuer{}, &CAIssuerList{})
	SchemeBuilder.Register(&ClusterCAIssuer{}, &ClusterCAIssuerList{})
}

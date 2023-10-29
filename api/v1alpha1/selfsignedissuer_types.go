package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SelfSignedIssuerSpec defines the desired state of SelfSignedIssuer
type SelfSignedIssuerSpec struct {
}

// SelfSignedIssuerStatus defines the observed state of SelfSignedIssuer
type SelfSignedIssuerStatus struct {
	// List of status conditions to indicate the status of a CertificateRequest.
	// Known condition types are `Ready`.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

const ConditionTypeReady = "Ready"

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SelfSignedIssuer is the Schema for the issuers API
type SelfSignedIssuer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SelfSignedIssuerSpec   `json:"spec,omitempty"`
	Status SelfSignedIssuerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SelfSignedIssuerList contains a list of SelfSignedIssuer
type SelfSignedIssuerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SelfSignedIssuer `json:"items"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ClusterSelfSignedIssuer is the Schema for the clusterissuers API
// +kubebuilder:resource:path=clustercaissuers,scope=Cluster
type ClusterSelfSignedIssuer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SelfSignedIssuerSpec   `json:"spec,omitempty"`
	Status SelfSignedIssuerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterSelfSignedIssuerList contains a list of ClusterSelfSignedIssuer
type ClusterSelfSignedIssuerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterSelfSignedIssuer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SelfSignedIssuer{}, &SelfSignedIssuerList{})
	SchemeBuilder.Register(&ClusterSelfSignedIssuer{}, &ClusterSelfSignedIssuerList{})
}

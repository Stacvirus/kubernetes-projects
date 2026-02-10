package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DummySiteSpec struct {
	WebsiteURL string `json:"website_url"`
}

type DummySiteStatus struct {
	Ready bool `json:"ready,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type DummySite struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DummySiteSpec   `json:"spec,omitempty"`
	Status DummySiteStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type DummySiteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DummySite `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DummySite{}, &DummySiteList{})
}

/*
Copyright 2023.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Route struct {
	//+kubebuilder:validation:MinLength=1
	//+kubebuilder:validation:MaxLength=255
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	prefix string `json:"prefix"`

	//+kubebuilder:validation:MinLength=1
	//+kubebuilder:validation:MaxLength=255
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	service string `json:"service"`
}

// HttpListenerSpec defines the desired state of HttpListener
type HttpListenerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//+kubebuilder:validation:MinLength=7
	//+kubebuilder:validation:MaxLength=15
	//+kubebuilder:validation:Required
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ipAddress string `json:"address"`

	//+kubebuilder:validation:Minimum=1
	//+kubebuilder:validation:Maximum=65535
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	port int32 `json:"port"`

	//+kubebuilder:validation:MinItems=0
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	routes []Route `json:"routes"`
}

// HttpListenerStatus defines the observed state of HttpListener
type HttpListenerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// HttpListener is the Schema for the httplisteners API
type HttpListener struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HttpListenerSpec   `json:"spec,omitempty"`
	Status HttpListenerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HttpListenerList contains a list of HttpListener
type HttpListenerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HttpListener `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HttpListener{}, &HttpListenerList{})
}

/*
Copyright 2023.

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

// LatencyMeasurementSpec nested structs
type Server struct {
	ClientNodeName  string `json:"clientNodeName,omitempty"`
	ServerNodeName  string `json:"serverNodeName,omitempty"`
	ServerIPAddress string `json:"serverIPAddress,omitempty"`
	ServerPort      int    `json:"serverPort,omitempty"`
}

type Client struct {
	IPAddress                string `json:"ipAddress,omitempty"`
	ServerPort               int    `json:"serverPort,omitempty"`
	Interval                 int    `json:"interval,omitempty"`
	Duration                 int    `json:"duration,omitempty"`
	MetricsAggregatorAddress string `json:"metricsAggregatorAddress,omitempty"`
	ClientNodeName           string `json:"clientNodeName,omitempty"`
	ServerNodeName           string `json:"serverNodeName,omitempty"`
	ClientSideClusterName    string `json:"clientSideClusterName,omitempty"`
	ServerSideClusterName    string `json:"serverSideClusterName,omitempty"`
}

// LatencyMeasurementSpec defines the desired state of LatencyMeasurement
type LatencyMeasurementSpec struct {
	Side    string   `json:"side,omitempty"`
	Servers []Server `json:"servers,omitempty"`
	Clients []Client `json:"clients,omitempty"`
}

// LatencyMeasurementStatus defines the observed state of LatencyMeasurement
type LatencyMeasurementStatus struct {
	State   string `json:"state,omitempty"`
	Details string `json:"details,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LatencyMeasurement is the Schema for the latencymeasurements API
type LatencyMeasurement struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LatencyMeasurementSpec   `json:"spec,omitempty"`
	Status LatencyMeasurementStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LatencyMeasurementList contains a list of LatencyMeasurement
type LatencyMeasurementList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LatencyMeasurement `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LatencyMeasurement{}, &LatencyMeasurementList{})
}

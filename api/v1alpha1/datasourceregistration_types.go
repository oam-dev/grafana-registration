/*


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

// DatasourceRegistrationSpec defines the desired state of DatasourceRegistration
type DatasourceRegistrationSpec struct {
	Grafana    Grafana    `json:"grafana"`
	Datasource Datasource `json:"datasource"`
}

// Grafana defines the access information for Grafana
type Grafana struct {
	Service string `json:"service"`
	// +kubebuilder:default:=default
	Namespace        string `json:"namespace"`
	CredentialSecret string `json:"credentialSecret"`
	// +kubebuilder:default:=default
	CredentialsSecretNamespace string `json:"credentialSecretNamespace,omitempty"`
}

// Datasource defines the information of a DataSource, like Loki, Prometheus
type Datasource struct {
	Name    string `json:"name"`
	Service string `json:"service"`
	// +kubebuilder:default:=default
	Namespace string `json:"namespace"`
	// +kubebuilder:default:=proxy
	Access string `json:"access,omitempty"`
	Type   string `json:"type"`
}

// DatasourceRegistrationStatus defines the observed state of DatasourceRegistration
type DatasourceRegistrationStatus struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
}

// DatasourceRegistration is the Schema for the DatasourceRegistration API

// +kubebuilder:object:root=true
type DatasourceRegistration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DatasourceRegistrationSpec   `json:"spec,omitempty"`
	Status DatasourceRegistrationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DatasourceRegistrationList contains a list of DatasourceRegistration
type DatasourceRegistrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DatasourceRegistration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DatasourceRegistration{}, &DatasourceRegistrationList{})
}

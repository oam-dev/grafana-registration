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

// DatasourceRegistrationDeprecatedSpec defines the desired state of DatasourceRegistrationDeprecated
type DatasourceRegistrationDeprecatedSpec struct {
	GrafanaURL       string `json:"grafanaUrl"`
	CredentialSecret string `json:"credentialSecret"`
	// +kubebuilder:default:=default
	CredentialsSecretNamespace string `json:"credentialSecretNamespace,omitempty"`
	Name                       string `json:"name"`
	URL                        string `json:"url"`
	// +kubebuilder:default:=proxy
	Access string `json:"access,omitempty"`
	Type   string `json:"type"`
}

// DatasourceRegistrationDeprecatedStatus defines the observed state of DatasourceRegistrationDeprecated
type DatasourceRegistrationDeprecatedStatus struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
}

// DatasourceRegistrationDeprecated is the Schema for the DatasourceRegistrationDeprecated API
// +kubebuilder:object:root=true
type DatasourceRegistrationDeprecated struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DatasourceRegistrationDeprecatedSpec   `json:"spec,omitempty"`
	Status DatasourceRegistrationDeprecatedStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DatasourceRegistrationDeprecatedList contains a list of DatasourceRegistrationDeprecated
type DatasourceRegistrationDeprecatedList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DatasourceRegistrationDeprecated `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DatasourceRegistrationDeprecated{}, &DatasourceRegistrationDeprecatedList{})
}

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

// ImportDashboardSpec defines the desired state of ImportDashboard
type ImportDashboardSpec struct {
	Grafana Grafana `json:"grafana"`
	URL     string  `json:"url"`
}

// ImportDashboardStatus defines the observed state of ImportDashboard
type ImportDashboardStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ImportDashboard is the Schema for the importdashboards API
type ImportDashboard struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImportDashboardSpec   `json:"spec,omitempty"`
	Status ImportDashboardStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ImportDashboardList contains a list of ImportDashboard
type ImportDashboardList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ImportDashboard `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ImportDashboard{}, &ImportDashboardList{})
}

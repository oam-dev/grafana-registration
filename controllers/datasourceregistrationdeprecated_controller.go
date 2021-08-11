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

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"github.com/zzxwill/grafana-datasource-registration/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	grafanaUser     = "admin-user"
	grafanaPassword = "admin-password"
)

// DatasourceRegistrationDeprecatedReconciler reconciles a DatasourceRegistrationDeprecated object
type DatasourceRegistrationDeprecatedReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

const datasourceRegistrationfinalizer = "grafana.extension.oam.dev/datasource"

// +kubebuilder:rbac:groups=grafana.extension.oam.dev,resources=datasourceregistrationdeprecated,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=grafana.extension.oam.dev,resources=datasourceregistrationdeprecateds/status,verbs=get;update;patch

func (r *DatasourceRegistrationDeprecatedReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	var (
		ctx = context.Background()
		dsr v1alpha1.DatasourceRegistrationDeprecated
	)
	if err := r.Get(ctx, req.NamespacedName, &dsr); err != nil {
		if kerrors.IsNotFound(err) {
			klog.ErrorS(err, "unable to fetch DatasourceRegistrationDeprecated", "NamespacedName", req.NamespacedName)
			err = nil
		}
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *DatasourceRegistrationDeprecatedReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.DatasourceRegistrationDeprecated{}).
		Complete(r)
}

func getGrafanaAuth(ctx context.Context, k8sClient client.Client, namespace, name string) (string, error) {
	var cred v1.Secret
	if err := k8sClient.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, &cred); err != nil {
		return "", errors.Wrap(err, "Grafana credential is not provided")
	}

	if cred.Data[grafanaUser] == nil || cred.Data[grafanaPassword] == nil {
		return "", errors.Errorf("%s or %s isn't in Grafana credential", grafanaUser, grafanaPassword)
	}
	return fmt.Sprintf("%s:%s", string(cred.Data[grafanaUser]), string(cred.Data[grafanaPassword])), nil
}

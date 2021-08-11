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
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/zzxwill/grafana-datasource-registration/api/v1alpha1"
)

// ImportDashboardReconciler reconciles a ImportDashboard object
type ImportDashboardReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=grafana.extension.oam.dev,resources=importdashboards,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=grafana.extension.oam.dev,resources=importdashboards/status,verbs=get;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ImportDashboard object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.6.4/pkg/reconcile
func (r *ImportDashboardReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	var (
		ctx       = context.Background()
		err       error
		dashboard v1alpha1.ImportDashboard
	)

	if err = r.Get(ctx, req.NamespacedName, &dashboard); err != nil {
		if kerrors.IsNotFound(err) {
			klog.ErrorS(err, "unable to fetch ImportDashboard", "NamespacedName", req.NamespacedName)
			err = nil
		}
		return ctrl.Result{}, err
	}
	klog.InfoS("Trying to retrieve Grafana Service", "Namespace", dashboard.Spec.Grafana.Namespace,
		"Name", dashboard.Spec.Grafana.Service)
	grafanaURL, err := getServiceURL(ctx, r.Client, dashboard.Spec.Grafana.Namespace, dashboard.Spec.Grafana.Service)
	if err != nil {
		return ctrl.Result{}, err
	}

	basicAuth, err := getGrafanaAuth(ctx, r.Client, dashboard.Spec.Grafana.CredentialsSecretNamespace, dashboard.Spec.Grafana.CredentialSecret)
	if err != nil {
		return ctrl.Result{}, err
	}

	board, err := downloadGrafanaDashboard(dashboard.Spec.URL)
	if err != nil {
		klog.InfoS("failed to downloaded Grafana Dashboard from RUL", "RUL", dashboard.Spec.URL)
		return ctrl.Result{}, err
	}
	klog.InfoS("successfully downloaded Grafana Dashboard and converted it to a *sdk.Board object", "Title", board.Title)

	if err := importDashboard(ctx, grafanaURL, basicAuth, board); err != nil {
		return ctrl.Result{}, err
	}
	klog.InfoS("successfully imported a Grafana Dashboard", "Title", board.Title)

	// TODO(zzxwill) add a finalizer and delete Dashboard when this CR is deleted.

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ImportDashboardReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.ImportDashboard{}).
		Complete(r)
}

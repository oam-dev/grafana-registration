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
	"github.com/grafana-tools/sdk"
	"github.com/pkg/errors"
	"github.com/zzxwill/grafana-datasource-registration/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

// DatasourceRegistrationReconciler reconciles a DatasourceRegistration object
type DatasourceRegistrationReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

const datasourceRegistrationfinalizer = "grafana.extension.oam.dev/datasource"

// +kubebuilder:rbac:groups=grafana.extension.oam.dev,resources=datasourceregistrations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=grafana.extension.oam.dev,resources=datasourceregistrations/status,verbs=get;update;patch

func (r *DatasourceRegistrationReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	var (
		ctx  = context.Background()
		err  error
		dsr  v1alpha1.DatasourceRegistration
		cred v1.Secret
	)
	const (
		grafanaUser     = "admin-user"
		grafanaPassword = "admin-password"
	)
	if err := r.Get(ctx, req.NamespacedName, &dsr); err != nil {
		if kerrors.IsNotFound(err) {
			klog.ErrorS(err, "unable to fetch DatasourceRegistration", "NamespacedName", req.NamespacedName)
			err = nil
		}
		return ctrl.Result{}, err
	}

	dataSourceName := dsr.Spec.Name
	dataSourceURL := dsr.Spec.URL
	dataSourceAccess := dsr.Spec.Access
	dataSourceType := dsr.Spec.Type
	grafanaURL := dsr.Spec.GrafanaURL

	if err := r.Client.Get(ctx, client.ObjectKey{Namespace: dsr.Spec.CredentialsSecretNamespace, Name: dsr.Spec.CredentialSecret}, &cred); err != nil {
		return ctrl.Result{}, errors.Wrap(err, "Grafana credential is not provided")
	}

	if cred.Data[grafanaUser] == nil || cred.Data[grafanaPassword] == nil {
		return ctrl.Result{}, errors.Wrap(err, fmt.Sprintf("%s or %s isn't in Grafana credential", grafanaUser, grafanaPassword))
	}

	basicAuth := fmt.Sprintf("%s:%s", string(cred.Data[grafanaUser]), string(cred.Data[grafanaPassword]))

	c, err := sdk.NewClient(grafanaURL, basicAuth, sdk.DefaultHTTPClient)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "Failed to create a client")
	}
	dataSources, err := c.GetAllDatasources(ctx)
	if err != nil {
		return ctrl.Result{}, err
	}
	var dataSourceExisted bool
	var existedDataSource *sdk.Datasource
	for _, existingDS := range dataSources {
		if existingDS.Name == dataSourceName {
			dataSourceExisted = true
			existedDataSource = &existingDS
			break
		}
	}

	if dsr.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(&dsr, datasourceRegistrationfinalizer) {
			controllerutil.AddFinalizer(&dsr, datasourceRegistrationfinalizer)
			if err := r.Client.Update(ctx, &dsr); err != nil {
				return ctrl.Result{}, errors.Wrap(err, "failed to add a finalizer")
			}
		}
	} else {
		if dataSourceExisted {
			if _, err = c.DeleteDatasource(ctx, existedDataSource.ID); err != nil {
				return ctrl.Result{}, errors.Wrap(err, "error on deleting datasource")
			}
		}
		controllerutil.RemoveFinalizer(&dsr, datasourceRegistrationfinalizer)
		if err := r.Update(ctx, &dsr); err != nil {
			return ctrl.Result{RequeueAfter: 3 * time.Second}, errors.Wrap(err, "failed to remove finalizer")
		}
		return ctrl.Result{}, nil
	}

	if dataSourceExisted {
		existedDataSource.Name = dataSourceName
		existedDataSource.URL = dataSourceURL
		existedDataSource.Type = dataSourceType
		existedDataSource.Access = dataSourceAccess
		_, err = c.UpdateDatasource(ctx, *existedDataSource)
		if err != nil {
			err = errors.Wrap(err, "error on deleting datasource")
		}
	} else {
		ds := sdk.Datasource{
			Name:   dataSourceName,
			URL:    dataSourceURL,
			Type:   dataSourceType,
			Access: dataSourceAccess,
		}
		_, err = c.CreateDatasource(ctx, ds)
		if err != nil {
			err = errors.Wrap(err, "failed to add datasource to Grafana")
		}
	}

	if err != nil {
		dsr.Status = v1alpha1.DatasourceRegistrationStatus{
			Success: false,
			Message: err.Error(),
		}
	} else {
		dsr.Status = v1alpha1.DatasourceRegistrationStatus{
			Success: true,
		}
	}

	if err := r.Client.Update(ctx, &dsr); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, errors.Wrap(err, "error on importing datasource")
}

func (r *DatasourceRegistrationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.DatasourceRegistration{}).
		Complete(r)
}

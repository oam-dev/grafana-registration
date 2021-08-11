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
	"github.com/grafana-tools/sdk"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/zzxwill/grafana-datasource-registration/api/v1alpha1"
)

// DatasourceRegistrationReconciler reconciles a DatasourceRegistration object
type DatasourceRegistrationReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=grafana.extension.oam.dev,resources=datasourceregistrations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=grafana.extension.oam.dev,resources=datasourceregistrations/status,verbs=get;update;patch

func (r *DatasourceRegistrationReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	var (
		ctx = context.Background()
		err error
		dsr v1alpha1.DatasourceRegistration
	)

	if err = r.Get(ctx, req.NamespacedName, &dsr); err != nil {
		if kerrors.IsNotFound(err) {
			klog.ErrorS(err, "unable to fetch DatasourceRegistration", "NamespacedName", req.NamespacedName)
			err = nil
		}
		return ctrl.Result{}, err
	}

	klog.InfoS("Trying to retrieve Grafana Service", "Namespace", dsr.Spec.Grafana.Namespace,
		"Name", dsr.Spec.Grafana.Service)
	grafanaURL, err := getServiceURL(ctx, r.Client, dsr.Spec.Grafana.Namespace, dsr.Spec.Grafana.Service)
	if err != nil {
		return ctrl.Result{}, err
	}

	klog.InfoS("Trying to retrieve Datasource Service", "Namespace", dsr.Spec.Grafana.Namespace,
		"Name", dsr.Spec.Grafana.Service)
	dataSourceURL, err := getServiceURL(ctx, r.Client, dsr.Spec.Datasource.Namespace, dsr.Spec.Datasource.Service)
	if err != nil {
		return ctrl.Result{}, err
	}

	if err := datasourceOperation(ctx, r.Client, dsr, grafanaURL, dataSourceURL); err != nil {
		klog.ErrorS(err, "failed to add or update datasource to Grafana", "Name", dsr.Spec.Datasource.Name)
		dsr.Status = v1alpha1.DatasourceRegistrationStatus{
			Success: false,
			Message: err.Error(),
		}
		if updateErr := r.Client.Update(ctx, &dsr); updateErr != nil {
			return ctrl.Result{}, updateErr
		}
		return ctrl.Result{}, err
	}

	klog.InfoS("successfully added or updated datasource to Grafana", "Name", dsr.Spec.Datasource.Name)
	dsr.Status = v1alpha1.DatasourceRegistrationStatus{
		Success: true,
	}
	if err := r.Client.Update(ctx, &dsr); err != nil {
		return ctrl.Result{}, err
	}
	klog.InfoS("successfully updated the status of DataSourceRegistration", "Name", dsr.Name,
		"Namespace", dsr.Namespace)
	return ctrl.Result{}, errors.Wrap(err, "error on importing datasource")
}

func (r *DatasourceRegistrationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.DatasourceRegistration{}).
		Complete(r)
}

func datasourceOperation(ctx context.Context, k8sClient client.Client, dsr v1alpha1.DatasourceRegistration, grafanaURL, dataSourceURL string) error {
	dataSourceName := dsr.Spec.Datasource.Name
	dataSourceAccess := dsr.Spec.Datasource.Access
	dataSourceType := dsr.Spec.Datasource.Type

	klog.InfoS("adding Datasource to Grafana", "Name", dataSourceName, "URL", dataSourceURL,
		"Access", dataSourceAccess, "Type", dataSourceType, "GrafanaURL", grafanaURL)

	basicAuth, err := getGrafanaAuth(ctx, k8sClient, dsr.Spec.Grafana.CredentialsSecretNamespace, dsr.Spec.Grafana.CredentialSecret)
	if err != nil {
		return err
	}

	c, err := sdk.NewClient(grafanaURL, basicAuth, sdk.DefaultHTTPClient)
	if err != nil {
		return errors.Wrap(err, "Failed to create a client")
	}
	dataSources, err := c.GetAllDatasources(ctx)
	if err != nil {
		return err
	}
	klog.InfoS("All datasources of Grafana", "Datasources", dataSources)
	var dataSourceExisted bool
	var existedDataSource *sdk.Datasource
	for _, existingDS := range dataSources {
		if existingDS.Name == dataSourceName {
			dataSourceExisted = true
			existedDataSource = &existingDS
			break
		}
	}

	if dataSourceExisted {
		klog.InfoS("The target Datasource doesn't exist", "Name", dataSourceName, "URL")
	}

	if dsr.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(&dsr, datasourceRegistrationfinalizer) {
			controllerutil.AddFinalizer(&dsr, datasourceRegistrationfinalizer)
			if err := k8sClient.Update(ctx, &dsr); err != nil {
				klog.InfoS("failed to add a finalizer to DataSourceRegistration", "Name", dsr.Name,
					"Namespace", dsr.Namespace)
				return errors.Wrap(err, "failed to add a finalizer")
			}
			klog.InfoS("successfully added a finalizer to DataSourceRegistration", "Name", dsr.Name,
				"Namespace", dsr.Namespace)
		}
	} else {
		if dataSourceExisted {
			if _, err = c.DeleteDatasource(ctx, existedDataSource.ID); err != nil {
				klog.InfoS("failed to delete DataSourceRegistration", "Name", dsr.Name,
					"Namespace", dsr.Namespace)
				return errors.Wrap(err, "error on deleting datasource")
			}
			klog.InfoS("successfully deleted DataSourceRegistration", "Name", dsr.Name,
				"Namespace", dsr.Namespace)
		}
		controllerutil.RemoveFinalizer(&dsr, datasourceRegistrationfinalizer)
		if err := k8sClient.Update(ctx, &dsr); err != nil {
			return errors.Wrap(err, "failed to remove finalizer")
		}
		klog.InfoS("successfully delete the finalizer from DataSourceRegistration", "Name", dsr.Name,
			"Namespace", dsr.Namespace)
		return nil
	}

	var operationErr error
	if dataSourceExisted {
		existedDataSource.Name = dataSourceName
		existedDataSource.URL = dataSourceURL
		existedDataSource.Type = dataSourceType
		existedDataSource.Access = dataSourceAccess
		_, operationErr = c.UpdateDatasource(ctx, *existedDataSource)
		if operationErr != nil {
			operationErr = errors.Wrap(operationErr, "error on deleting datasource")
		}
	} else {
		ds := sdk.Datasource{
			Name:   dataSourceName,
			URL:    dataSourceURL,
			Type:   dataSourceType,
			Access: dataSourceAccess,
		}
		_, operationErr = c.CreateDatasource(ctx, ds)
		if operationErr != nil {
			operationErr = errors.Wrap(operationErr, "failed to add datasource to Grafana")
		}
	}

	return errors.Wrap(operationErr, "error on importing datasource")
}

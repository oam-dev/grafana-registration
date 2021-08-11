package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/grafana-tools/sdk"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func importDashboard(ctx context.Context, grafanaURL, basicAuth string, board *sdk.Board) error {
	c, err := sdk.NewClient(grafanaURL, basicAuth, sdk.DefaultHTTPClient)
	if err != nil {
		return errors.Wrap(err, "Failed to create a client")
	}
	params := sdk.SetDashboardParams{
		FolderID:  sdk.DefaultFolderId,
		Overwrite: true,
	}

	if _, _, err := c.GetDashboardByUID(ctx, board.UID); err != nil {
		board.ID = 0
	}

	_, err = c.SetDashboard(ctx, *board, params)
	if err != nil {
		errMsg := "hit an error to import a dashboard"
		klog.ErrorS(err, errMsg, "GrafanaURL", grafanaURL, "Title", board.Title)
		return errors.Wrap(err, errMsg)
	}
	return nil
}

func downloadGrafanaDashboard(url string) (*sdk.Board, error) {
	var board sdk.Board
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&board); err != nil {
		return nil, err
	}
	return &board, nil
}

func getServiceURL(ctx context.Context, k8sClient client.Client, namespace, name string) (string, error) {
	var svc v1.Service
	if err := k8sClient.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, &svc); err != nil {
		klog.ErrorS(err, "failed to get service")
		return "", errors.Wrap(err, "failed to get service")
	}
	klog.InfoS("successfully retrieved Service")
	if svc.Spec.ClusterIP == "" || len(svc.Spec.Ports) == 0 {
		errMsg := "The ClusterIP or Port of the Service is not rendered"
		klog.Info(errMsg)
		return "", errors.New(errMsg)
	}
	return fmt.Sprintf("http://%s:%d", svc.Spec.ClusterIP, svc.Spec.Ports[0].Port), nil
}

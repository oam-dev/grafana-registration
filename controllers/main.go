package controllers

import (
	"context"
	"fmt"
	"os"

	"github.com/grafana-tools/sdk"
)

func main() {
	const (
		grafanaURL    = "GrafanaUrl"
		adminUser     = "admin-user"
		adminPassword = "admin-password"
		name          = "Name"
		url           = "URL"
		proxy         = "Proxy"
		dsType        = "Type"
	)

	dataSourceName := os.Getenv(name)
	dataSourceURL := os.Getenv(url)
	dataSourceProxy := os.Getenv(proxy)
	if dataSourceProxy == "" {
		dataSourceProxy = "proxy"
	}
	dataSourceType := os.Getenv(dsType)
	if os.Getenv(grafanaURL) == "" || os.Getenv(adminUser) == "" || os.Getenv(adminPassword) == "" ||
		dataSourceName == "" || dataSourceURL == "" || dataSourceType == "" {
		fmt.Println("environments are not complete")
		os.Exit(1)
	}
	basicAuth := fmt.Sprintf("%s:%s", os.Getenv(adminUser), os.Getenv(adminPassword))

	ctx := context.Background()
	c, err := sdk.NewClient(os.Getenv(grafanaURL), basicAuth, sdk.DefaultHTTPClient)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create a client: %s\n", err)
		os.Exit(1)
	}
	dataSources, err := c.GetAllDatasources(ctx)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	ds := sdk.Datasource{
		Name:   dataSourceName,
		URL:    dataSourceURL,
		Type:   dataSourceType,
		Access: dataSourceProxy,
	}
	for _, existingDS := range dataSources {
		if existingDS.Name == dataSourceName {
			if _, err = c.DeleteDatasource(ctx, existingDS.ID); err != nil {
				fmt.Fprintf(os.Stderr, "error on deleting datasource %s with %s", dataSourceName, err)
			}
			break
		}
	}
	if status, err := c.CreateDatasource(ctx, ds); err != nil {
		fmt.Fprintf(os.Stderr, "error on importing datasource %s with %s: %v", dataSourceName, err, status)
	}
}

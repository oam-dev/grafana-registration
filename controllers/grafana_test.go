package controllers

//func TestImportDashboard(t *testing.T) {
//	ctx := context.Background()
//	url := "https://raw.githubusercontent.com/zzxwill/grafana-dashboards/master/Kubernetes_for_Prometheus_Dashboard_CN_20201209_1628671317889.json"
//	board, _ := downloadGrafanaDashboard(url)
//	grafanaURL := "http://grafana.c276f4dac730c47b8b8988905e3c68fcf.cn-hongkong.alicontainer.com/"
//	basicAuth := "admin:vspgFXHScKDNr0sEdmOaLRGP7EKUsuEcCFCYraUi"
//	err := importDashboard(ctx, grafanaURL, basicAuth, board)
//	assert.NoError(t, err)
//}
//
//func TestDownloadGrafanaDashboard(t *testing.T) {
//	url := "https://raw.githubusercontent.com/zzxwill/grafana-dashboards/master/Kubernetes_for_Prometheus_Dashboard_CN_20201209_1628671317889.json"
//	board, err := downloadGrafanaDashboard(url)
//	assert.NoError(t, err)
//	assert.NotNil(t, board)
//}

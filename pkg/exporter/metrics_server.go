package exporter

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type MetricsServer struct {
	address string
	path    string
}

func NewMetricsServer(address, path string) *MetricsServer {
	return &MetricsServer{
		address,
		path,
	}
}

func (m *MetricsServer) ListenAndServe() error {
	http.Handle(m.path, promhttp.Handler())
	return http.ListenAndServe(m.address, nil)
}

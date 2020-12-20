package main

import (
	"github.com/patrick246/jicofo-exporter/pkg/exporter"
	"github.com/patrick246/jicofo-exporter/pkg/jicofo"
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

func main() {
	statsClient := jicofo.NewStatsClient("http://localhost:8888")
	server := exporter.NewMetricsServer("127.0.0.1:9094", "/metrics")
	jicofoExporter := exporter.NewExporter(statsClient)
	prometheus.MustRegister(jicofoExporter)
	log.Fatal(server.ListenAndServe())
}

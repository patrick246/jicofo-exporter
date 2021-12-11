package main

import (
	"flag"
	"github.com/patrick246/jicofo-exporter/pkg/exporter"
	"github.com/patrick246/jicofo-exporter/pkg/jicofo"
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

func main() {
	jicofoAddr := flag.String("jicofo.url", "http://localhost:8888", "Jicofo API URL (protocol://host:port")
	metricsAddr := flag.String("metrics.addr", "127.0.0.1:9094", "Listener Address (host:port)")
	metricsPath := flag.String("metrics.path", "/metrics", "Metrics Path")
	flag.Parse()

	statsClient := jicofo.NewStatsClient(*jicofoAddr)
	server := exporter.NewMetricsServer(*metricsAddr, *metricsPath)
	jicofoExporter := exporter.NewExporter(statsClient)
	prometheus.MustRegister(jicofoExporter)
	log.Fatal(server.ListenAndServe())
}

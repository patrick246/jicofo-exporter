package exporter

import (
	"github.com/patrick246/jicofo-exporter/pkg/jicofo"
	"github.com/prometheus/client_golang/prometheus"
	"math"
	"strings"
)

type JicofoStatsExporter struct {
	client *jicofo.StatsClient
}

func NewExporter(client *jicofo.StatsClient) *JicofoStatsExporter {
	return &JicofoStatsExporter{client: client}
}

func (j *JicofoStatsExporter) Describe(descs chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(j, descs)
}

func (j *JicofoStatsExporter) Collect(metrics chan<- prometheus.Metric) {
	stats, histograms, err := j.client.GetStats()
	if err != nil {
		metrics <- prometheus.NewInvalidMetric(prometheus.NewInvalidDesc(err), err)
	}

	for metricName, metricValue := range stats {
		metricType := prometheus.GaugeValue
		if strings.Contains(metricName, "_total_") {
			metricType = prometheus.CounterValue
		}

		metrics <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(metricName, "", []string{}, make(map[string]string)),
			metricType,
			metricValue,
		)
	}

	for histogramName, histogramValue := range histograms {
		sum := histogramValue[math.Inf(1)]
		metrics <- prometheus.MustNewConstHistogram(
			prometheus.NewDesc(histogramName, "", []string{}, make(map[string]string)),
			sum,
			float64(sum),
			histogramValue,
		)
	}
}

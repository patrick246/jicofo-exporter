package jicofo

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"math"
	"net/http"
	"time"
)

type StatsClient struct {
	url string
}

var (
	jicofoScrapeDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "jicofoexporter_scrape_duration_seconds",
		Help:    "request latencies for the jicofo /stats endpoint",
		Buckets: prometheus.DefBuckets,
	})
)

func init() {
	prometheus.MustRegister(jicofoScrapeDuration)
}

func NewStatsClient(url string) *StatsClient {
	return &StatsClient{url: url}
}

func (s *StatsClient) GetStats() (map[string]float64, map[string]map[float64]uint64, error) {
	resp, err := s.fetch()
	if err != nil {
		return nil, nil, err
	}

	decoder := json.NewDecoder(resp.Body)

	var statsResponse map[string]interface{}
	err = decoder.Decode(&statsResponse)
	if err != nil {
		return nil, nil, err
	}

	stats, histograms := convertStats(statsResponse, "jicofo_")

	return stats, histograms, nil
}

func (s *StatsClient) fetch() (*http.Response, error) {
	start := time.Now()
	defer func() {
		duration := time.Now().Sub(start).Seconds()
		jicofoScrapeDuration.Observe(duration)
	}()

	return http.Get(s.url + "/stats")
}

func convertStats(stats map[string]interface{}, prefix string) (map[string]float64, map[string]map[float64]uint64) {
	singleLevelStats := make(map[string]float64)
	histograms := make(map[string]map[float64]uint64)

	for key, value := range stats {
		switch v := value.(type) {
		case float64:
			singleLevelStats[prefix+key] = v
		case bool:
			if v {
				singleLevelStats[prefix+key] = 1
			} else {
				singleLevelStats[prefix+key] = 0
			}
		case map[string]interface{}:
			nestedSingle, nestedHistogram := convertStats(v, prefix+key+"_")
			singleLevelStats = mergeMapsFloat(singleLevelStats, nestedSingle)
			histograms = mergeMapsHistogram(histograms, nestedHistogram)

		case []interface{}:
			histograms[prefix+key] = convertHistogram(v)
		default:
			fmt.Printf("skipping %s", key)
		}
	}
	return singleLevelStats, histograms
}

func convertHistogram(data []interface{}) map[float64]uint64 {
	histogram := make(map[float64]uint64)
	runningCount := uint64(0)
	for i, val := range data {
		runningCount += uint64(val.(float64))
		histogram[float64(i)] = runningCount
		if i == len(data)-1 {
			histogram[math.Inf(1)] = runningCount
		}
	}
	return histogram
}

func mergeMapsFloat(map1, map2 map[string]float64) map[string]float64 {
	merged := make(map[string]float64)
	for k, v := range map1 {
		merged[k] = v
	}

	for k, v := range map2 {
		merged[k] = v
	}
	return merged
}

func mergeMapsHistogram(map1, map2 map[string]map[float64]uint64) map[string]map[float64]uint64 {
	merged := make(map[string]map[float64]uint64)
	for k, v := range map1 {
		merged[k] = v
	}

	for k, v := range map2 {
		merged[k] = v
	}
	return merged
}

/* 
Copyright 2021 Acacio Cruz acacio@acacio.coom

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

package grpcutil

import (
	"net/http"
	_ "net/http/pprof" // "/debug/pprof/trace" handler

	// curl http://localhost:9999/debug/pprof/trace?seconds=5 -o trace.out
	// go tool trace trace.out  ( /userregions  & /usertasks )

	// "github.com/prometheus/client_golang/prometheus"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
)

// EnablePrometheus - sets up gRPC metrics collection and launch web collection endpoint
func EnablePrometheus(s *grpc.Server, PORT string) http.Handler {
	// After your registrations, Prometheus metrics are initialized.
	println("TASKSERVER: Registering Prometheus...")
	grpc_prometheus.Register(s)
	grpc_prometheus.EnableHandlingTimeHistogram()
	return promhttp.Handler()
}

// GetgRPCMetrics returns the grpc_prometheus metrics object
func GetgRPCMetrics() map[string]float64 {
	metrics := grpc_prometheus.DefaultServerMetrics

	lats := make(map[string]float64)

	// Collect all gRPC metrics and build map
	c := make(chan prometheus.Metric)
	go func() {
		metrics.Collect(c)
		close(c)
	}()
	// iterate through sent metrics (via channel)
	for metric := range c {
		data := dto.Metric{}
		metric.Write(&data)
		if data.Histogram != nil {
			count := *data.Histogram.SampleCount
			sum := *data.Histogram.SampleSum
			latency := (1000.0 * sum) / float64(count) // s --> milliseconds
			labels := data.Label
			method := getMethod(labels)
			lats[method] = latency
			// fmt.Println("LATENCY:", method, latency)
		}
	}
	return lats
}

// GetgRPCHistograms grabs the data for all methods
func GetgRPCHistograms() map[string]map[float64]uint64 {
	metrics := grpc_prometheus.DefaultServerMetrics
	// Collect all gRPC metrics and build map
	c := make(chan prometheus.Metric)
	go func() {
		metrics.Collect(c)
		close(c)
	}()
	// iterate through sent metrics (via channel)
	histPerMethod := make(map[string]map[float64]uint64)
	for metric := range c {
		data := dto.Metric{}
		metric.Write(&data)
		if data.Histogram != nil {
			buckets := data.Histogram.GetBucket()
			labels := data.Label
			method := getMethod(labels)
			hist := make(map[float64]uint64)
			var prev uint64 = 0
			for _, v := range buckets {
				max := v.GetUpperBound()
				cnt := v.GetCumulativeCount()
				// z := v.String()
				hist[max] = cnt - prev // record just the differential
				prev = cnt
				// fmt.Println("BUCKET:", method, max, cnt, z)
			}
			histPerMethod[method] = hist
		}
	}
	return histPerMethod
}

func getMethod(labels []*dto.LabelPair) string {
	for _, v := range labels {
		if *v.Name == "grpc_method" {
			return *v.Value
		}
	}
	return "UNKNOWN"
}

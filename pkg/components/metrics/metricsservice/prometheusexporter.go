/*
Copyright 2023 The Radius Authors.

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

package metricsservice

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// PrometheusExporter is the struct that holds the metrics reklated data
type PrometheusExporter struct {
	// MeterProvider is used in the creation and coordination of Meters
	MeterProvider *sdkmetric.MeterProvider

	// Handler is the HTTP handler with basic metrics
	Handler http.Handler
}

// NewPrometheusExporter creates a PrometheusExporter instance with a MeterProvider and a Handler,
// and returns it along with any errors.
func NewPrometheusExporter(options *Options) (*PrometheusExporter, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
		sdkmetric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(options.ServiceName),
		)))

	otel.SetMeterProvider(mp)

	return &PrometheusExporter{
		MeterProvider: mp,
		Handler:       promhttp.Handler(),
	}, nil
}

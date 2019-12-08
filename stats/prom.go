package stats

import "github.com/prometheus/client_golang/prometheus"

const (
	namespaceMsx = "msx"
)

var standardBuckets = prometheus.ExponentialBuckets(10, 2, 16)

func NewGauge(subsystem, name string) prometheus.Gauge {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespaceMsx,
		Subsystem: subsystem,
		Name:      name,
	})

	prometheus.MustRegister(gauge)
	return gauge
}

func NewGaugeVec(subsystem, name string, labels ...string) *prometheus.GaugeVec {
	gaugeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespaceMsx,
		Subsystem: subsystem,
		Name:      name,
	}, labels)

	prometheus.MustRegister(gaugeVec)
	return gaugeVec
}

func NewHistogram(subsystem, name string, buckets []float64) prometheus.Histogram {
	if buckets == nil {
		buckets = standardBuckets
	}

	histogram := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace:   namespaceMsx,
		Subsystem:   subsystem,
		Name:        name,
		ConstLabels: nil,
		Buckets:     buckets,
	})

	prometheus.MustRegister(histogram)
	return histogram
}

func NewHistogramVec(subsystem, name string, buckets []float64, labels ...string) *prometheus.HistogramVec {
	if buckets == nil {
		buckets = standardBuckets
	}

	histogramVec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespaceMsx,
		Subsystem: subsystem,
		Name:      name,
		Buckets:   buckets,
	}, labels)

	prometheus.MustRegister(histogramVec)
	return histogramVec
}

func NewCounter(subsystem, name string) prometheus.Counter {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespaceMsx,
		Subsystem: subsystem,
		Name:      name,
	})

	prometheus.MustRegister(counter)
	return counter
}

func NewCounterVec(subsystem, name string, labels ...string) *prometheus.CounterVec {
	counterVec := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespaceMsx,
		Subsystem: subsystem,
		Name:      name,
	}, labels)

	prometheus.MustRegister(counterVec)
	return counterVec
}

var ExponentialBuckets = prometheus.ExponentialBuckets

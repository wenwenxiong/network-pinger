package pinger

import "github.com/prometheus/client_golang/prometheus"

var (
	apiserverHealthyGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "pinger_apiserver_healthy",
			Help: "If the apiserver request is healthy on this node",
		},
		[]string{
			"nodeName",
		})
	apiserverUnhealthyGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "pinger_apiserver_unhealthy",
			Help: "If the apiserver request is unhealthy on this node",
		},
		[]string{
			"nodeName",
		})
	apiserverRequestLatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "pinger_apiserver_latency_ms",
			Help:    "The latency ms histogram the node request apiserver",
			Buckets: []float64{2, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50},
		},
		[]string{
			"nodeName",
		})
	internalDNSHealthyGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "pinger_internal_dns_healthy",
			Help: "If the internal dns request is healthy on this node",
		},
		[]string{
			"nodeName",
		})
	internalDNSUnhealthyGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "pinger_internal_dns_unhealthy",
			Help: "If the internal dns request is unhealthy on this node",
		},
		[]string{
			"nodeName",
		})
	internalDNSRequestLatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "pinger_internal_dns_latency_ms",
			Help:    "The latency ms histogram the node request internal dns",
			Buckets: []float64{2, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50},
		},
		[]string{
			"nodeName",
		})
)

func InitPingerMetrics() {

	prometheus.MustRegister(apiserverHealthyGauge)
	prometheus.MustRegister(apiserverUnhealthyGauge)
	prometheus.MustRegister(apiserverRequestLatencyHistogram)
	prometheus.MustRegister(internalDNSHealthyGauge)
	prometheus.MustRegister(internalDNSUnhealthyGauge)
	prometheus.MustRegister(internalDNSRequestLatencyHistogram)
}

func SetApiserverUnhealthyMetrics(nodeName string) {
	apiserverHealthyGauge.WithLabelValues(nodeName).Set(0)
	apiserverUnhealthyGauge.WithLabelValues(nodeName).Set(1)
}

func SetApiserverHealthyMetrics(nodeName string, latency float64) {
	apiserverHealthyGauge.WithLabelValues(nodeName).Set(1)
	apiserverRequestLatencyHistogram.WithLabelValues(nodeName).Observe(latency)
	apiserverUnhealthyGauge.WithLabelValues(nodeName).Set(0)
}

func SetInternalDNSHealthyMetrics(nodeName string, latency float64) {
	internalDNSHealthyGauge.WithLabelValues(nodeName).Set(1)
	internalDNSRequestLatencyHistogram.WithLabelValues(nodeName).Observe(latency)
	internalDNSUnhealthyGauge.WithLabelValues(nodeName).Set(0)
}

func SetInternalDNSUnhealthyMetrics(nodeName string) {
	internalDNSHealthyGauge.WithLabelValues(nodeName).Set(0)
	internalDNSUnhealthyGauge.WithLabelValues(nodeName).Set(1)
}
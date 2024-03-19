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
	podPingLatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "pinger_pod_ping_latency_ms",
			Help:    "The latency ms histogram for pod peer ping",
			Buckets: []float64{.25, .5, 1, 2, 5, 10, 30},
		},
		[]string{
			"src_node_name",
			"src_node_ip",
			"src_pod_ip",
			"target_node_name",
			"target_node_ip",
			"target_pod_ip",
		})
	podPingLostCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pinger_pod_ping_lost_total",
			Help: "The lost count for pod peer ping",
		}, []string{
			"src_node_name",
			"src_node_ip",
			"src_pod_ip",
			"target_node_name",
			"target_node_ip",
			"target_pod_ip",
		})
	podPingTotalCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pinger_pod_ping_count_total",
			Help: "The total count for pod peer ping",
		}, []string{
			"src_node_name",
			"src_node_ip",
			"src_pod_ip",
			"target_node_name",
			"target_node_ip",
			"target_pod_ip",
		})
	nodePingLatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "pinger_node_ping_latency_ms",
			Help:    "The latency ms histogram for pod ping node",
			Buckets: []float64{.25, .5, 1, 2, 5, 10, 30},
		},
		[]string{
			"src_node_name",
			"src_node_ip",
			"src_pod_ip",
			"target_node_name",
			"target_node_ip",
		})
	nodePingLostCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pinger_node_ping_lost_total",
			Help: "The lost count for pod ping node",
		}, []string{
			"src_node_name",
			"src_node_ip",
			"src_pod_ip",
			"target_node_name",
			"target_node_ip",
		})
	nodePingTotalCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pinger_node_ping_count_total",
			Help: "The total count for pod ping node",
		}, []string{
			"src_node_name",
			"src_node_ip",
			"src_pod_ip",
			"target_node_name",
			"target_node_ip",
		})
	IpPingLatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "pinger_ip_ping_latency_ms",
			Help:    "The latency ms histogram for ip peer ping",
			Buckets: []float64{.25, .5, 1, 2, 5, 10, 30},
		},
		[]string{
			"src_node_name",
			"src_node_ip",
			"src_pod_ip",
			"target_ip",
		})
	IpPingLostCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pinger_ip_ping_lost_total",
			Help: "The lost count for ip peer ping",
		}, []string{
			"src_node_name",
			"src_node_ip",
			"src_pod_ip",
			"target_ip",
		})
	IpPingTotalCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pinger_ip_ping_count_total",
			Help: "The total count for ip peer ping",
		}, []string{
			"src_node_name",
			"src_node_ip",
			"src_pod_ip",
			"target_ip",
		})
)

func InitPingerMetrics() {
	prometheus.MustRegister(apiserverHealthyGauge)
	prometheus.MustRegister(apiserverUnhealthyGauge)
	prometheus.MustRegister(apiserverRequestLatencyHistogram)
	prometheus.MustRegister(internalDNSHealthyGauge)
	prometheus.MustRegister(internalDNSUnhealthyGauge)
	prometheus.MustRegister(internalDNSRequestLatencyHistogram)
	prometheus.MustRegister(podPingLatencyHistogram)
	prometheus.MustRegister(podPingLostCounter)
	prometheus.MustRegister(podPingTotalCounter)
	prometheus.MustRegister(nodePingLatencyHistogram)
	prometheus.MustRegister(nodePingLostCounter)
	prometheus.MustRegister(nodePingTotalCounter)
	prometheus.MustRegister(IpPingLatencyHistogram)
	prometheus.MustRegister(IpPingLostCounter)
	prometheus.MustRegister(IpPingTotalCounter)

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
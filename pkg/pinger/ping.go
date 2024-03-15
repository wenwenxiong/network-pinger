package pinger

import (
	"context"
	"fmt"
	goping "github.com/prometheus-community/pro-bing"
	"github.com/wenwenxiong/network-pinger/pkg/util"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog/v2"
	"math"
	"net"
	"os"
	"time"
)

func StartPinger(config *Configuration) {
	errHappens := false
	for {
		if ping(config) != nil {
			errHappens = true
		}

		if config.Mode != "server" {
			break
		}

		time.Sleep(time.Duration(config.Interval) * time.Second)
	}
	if errHappens && config.ExitCode != 0 {
		os.Exit(config.ExitCode)
	}
}

func ping(config *Configuration) error {
	errHappens := false
	if checkAPIServer(config) != nil {
		errHappens = true
	}
	if pingPods(config) != nil {
		errHappens = true
	}
	if pingNodes(config) != nil {
		errHappens = true
	}
	if internalNslookup(config) != nil {
		errHappens = true
	}

	if errHappens {
		return fmt.Errorf("ping failed")
	}
	return nil
}

func checkAPIServer(config *Configuration) error {
	klog.Infof("start to check apiserver connectivity")
	t1 := time.Now()
	_, err := config.KubeClient.Discovery().ServerVersion()
	elapsed := time.Since(t1)
	if err != nil {
		klog.Errorf("failed to connect to apiserver: %v", err)
		SetApiserverUnhealthyMetrics(config.NodeName)
		return err
	}
	klog.Infof("connect to apiserver success in %.2fms", float64(elapsed)/float64(time.Millisecond))
	SetApiserverHealthyMetrics(config.NodeName, float64(elapsed)/float64(time.Millisecond))
	return nil
}

func pingPods(config *Configuration) error {
	klog.Infof("start to check pod connectivity")
	pods, err := config.KubeClient.CoreV1().Pods(config.DaemonSetNamespace).List(context.Background(), metav1.ListOptions{LabelSelector: labels.Set(config.MatchLabels).String()})
	if err != nil {
		klog.Errorf("failed to list peer pods: %v", err)
		return err
	}

	var pingErr error
	for _, pod := range pods.Items {
		for _, podIP := range pod.Status.PodIPs {
			if util.ContainsString(config.PodProtocols, util.CheckProtocol(podIP.IP)) {
				func(podIP, podName, nodeIP, nodeName string) {

					pinger, err := goping.NewPinger(podIP)
					if err != nil {
						klog.Errorf("failed to init pinger, %v", err)
						pingErr = err
						return
					}
					pinger.SetPrivileged(true)
					pinger.Timeout = 1 * time.Second
					pinger.Debug = true
					pinger.Count = 3
					pinger.Interval = 100 * time.Millisecond
					if err = pinger.Run(); err != nil {
						klog.Errorf("failed to run pinger for destination %s: %v", podIP, err)
						pingErr = err
						return
					}

					stats := pinger.Statistics()
					klog.Infof("ping pod: %s %s, count: %d, loss count %d, average rtt %.2fms",
						podName, podIP, pinger.Count, int(math.Abs(float64(stats.PacketsSent-stats.PacketsRecv))), float64(stats.AvgRtt)/float64(time.Millisecond))
					if int(math.Abs(float64(stats.PacketsSent-stats.PacketsRecv))) != 0 {
						pingErr = fmt.Errorf("ping failed")
					}
					SetPodPingMetrics(
						config.NodeName,
						config.HostIP,
						config.PodName,
						nodeName,
						nodeIP,
						podIP,
						float64(stats.AvgRtt)/float64(time.Millisecond),
						int(math.Abs(float64(stats.PacketsSent-stats.PacketsRecv))),
						int(float64(stats.PacketsSent)))
				}(podIP.IP, pod.Name, pod.Status.HostIP, pod.Spec.NodeName)
			}
		}
	}
	return pingErr
}

func pingNodes(config *Configuration) error {
	klog.Infof("start to check node connectivity")
	nodes, err := config.KubeClient.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		klog.Errorf("failed to list nodes, %v", err)
		return err
	}

	var pingErr error
	for _, no := range nodes.Items {
		for _, addr := range no.Status.Addresses {
			if addr.Type == v1.NodeInternalIP && util.ContainsString(config.PodProtocols, util.CheckProtocol(addr.Address)) {
				func(nodeIP, nodeName string) {

					pinger, err := goping.NewPinger(nodeIP)
					if err != nil {
						klog.Errorf("failed to init pinger, %v", err)
						pingErr = err
						return
					}
					pinger.SetPrivileged(true)
					pinger.Timeout = 30 * time.Second
					pinger.Count = 3
					pinger.Interval = 100 * time.Millisecond
					pinger.Debug = true
					if err = pinger.Run(); err != nil {
						klog.Errorf("failed to run pinger for destination %s: %v", nodeIP, err)
						pingErr = err
						return
					}

					stats := pinger.Statistics()
					klog.Infof("ping node: %s %s, count: %d, loss count %d, average rtt %.2fms",
						nodeName, nodeIP, pinger.Count, int(math.Abs(float64(stats.PacketsSent-stats.PacketsRecv))), float64(stats.AvgRtt)/float64(time.Millisecond))
					if int(math.Abs(float64(stats.PacketsSent-stats.PacketsRecv))) != 0 {
						pingErr = fmt.Errorf("ping failed")
					}
					SetNodePingMetrics(
						config.NodeName,
						config.HostIP,
						config.PodName,
						no.Name, addr.Address,
						float64(stats.AvgRtt)/float64(time.Millisecond),
						int(math.Abs(float64(stats.PacketsSent-stats.PacketsRecv))),
						int(float64(stats.PacketsSent)))
				}(addr.Address, no.Name)
			}
		}
	}
	return pingErr
}

func internalNslookup(config *Configuration) error {
	klog.Infof("start to check dns connectivity")
	t1 := time.Now()
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	var r net.Resolver
	addrs, err := r.LookupHost(ctx, config.InternalDNS)
	elapsed := time.Since(t1)
	if err != nil {
		klog.Errorf("failed to resolve dns %s, %v", config.InternalDNS, err)
		SetInternalDNSUnhealthyMetrics(config.NodeName)
		return err
	}
	SetInternalDNSHealthyMetrics(config.NodeName, float64(elapsed)/float64(time.Millisecond))
	klog.Infof("resolve dns %s to %v in %.2fms", config.InternalDNS, addrs, float64(elapsed)/float64(time.Millisecond))
	return nil
}

func SetPodPingMetrics(srcNodeName, srcNodeIP, srcPodIP, targetNodeName, targetNodeIP, targetPodIP string, latency float64, lost, total int) {
	podPingLatencyHistogram.WithLabelValues(
		srcNodeName,
		srcNodeIP,
		srcPodIP,
		targetNodeName,
		targetNodeIP,
		targetPodIP,
	).Observe(latency)
	podPingLostCounter.WithLabelValues(
		srcNodeName,
		srcNodeIP,
		srcPodIP,
		targetNodeName,
		targetNodeIP,
		targetPodIP,
	).Add(float64(lost))
	podPingTotalCounter.WithLabelValues(
		srcNodeName,
		srcNodeIP,
		srcPodIP,
		targetNodeName,
		targetNodeIP,
		targetPodIP,
	).Add(float64(total))
}

func SetNodePingMetrics(srcNodeName, srcNodeIP, srcPodIP, targetNodeName, targetNodeIP string, latency float64, lost, total int) {
	nodePingLatencyHistogram.WithLabelValues(
		srcNodeName,
		srcNodeIP,
		srcPodIP,
		targetNodeName,
		targetNodeIP,
	).Observe(latency)
	nodePingLostCounter.WithLabelValues(
		srcNodeName,
		srcNodeIP,
		srcPodIP,
		targetNodeName,
		targetNodeIP,
	).Add(float64(lost))
	nodePingTotalCounter.WithLabelValues(
		srcNodeName,
		srcNodeIP,
		srcPodIP,
		targetNodeName,
		targetNodeIP,
	).Add(float64(total))
}
功能
network-pinger 检测指定label的pod ips的网络指标 可达性，时延
并且检测所在的k8s集群的apiserver、dns、nodes是否健康
可以访问应用的/metrics获取Prometheus标准的指标数值
支持的metric
```
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
```

部署前提
部署网络可达探测IP，目前解决方式加多网卡

编译
```bigquery
make build-go
```
打镜像
```bigquery
make image-network-pinger
```
deployment目录下保存部署yaml文件
部署
```bigquery
kubectl apply -f deployment/network-pinger.yaml
```

#TODO
1、目前master分支上 只能ping pod的calico ip
开2个分支
分支mec_dev：
  先解决外部网络IP的探测，然后解决子网ip的探测
分支5gc_dev：
  解决whereabouts dhcp分配的ip 探测 
  5gc static ip的探测

2、helm部署
支持helm部署，支持makefile自动生成部署yaml文件
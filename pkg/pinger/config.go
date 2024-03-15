package pinger

import (
	"flag"
	"k8s.io/client-go/kubernetes"
	"os"
	"time"
)

type Configuration struct {
	KubeConfigFile     string
	KubeClient         kubernetes.Interface
	Port               int
	MatchLabels	       string
	Interval           int
	Mode               string
	ExitCode           int
	InternalDNS        string
	ExternalDNS        string
	NodeName           string
	HostIP             string
	PodName            string
	PodIP              string
	PodProtocols       []string
	ExternalAddress    string
	NetworkMode        string
	EnableMetrics      bool
}

func ParseFlags() (*Configuration, error) {
	var (
		argPort = pflag.Int("port", 8080, "metrics port")

		argKubeConfigFile     = pflag.String("kubeconfig", "", "Path to kubeconfig file with authorization and master location information. If not set use the inCluster token.")
		argInterval           = pflag.Int("interval", 5, "interval seconds between consecutive pings")
		argMode               = pflag.String("mode", "server", "server or job Mode")
		argExitCode           = pflag.Int("exit-code", 0, "exit code when failure happens")
		argInternalDNS        = pflag.String("internal-dns", "kubernetes.default", "check dns from pod")
		argExternalDNS        = pflag.String("external-dns", "", "check external dns resolve from pod")
		argExternalAddress    = pflag.String("external-address", "", "check ping connection to an external address, default: 114.114.114.114")

		argNetworkMode        = pflag.String("network-mode", "kube-ovn", "The cni plugin current cluster used, default: kube-ovn")
		argEnableMetrics      = pflag.Bool("enable-metrics", true, "Whether to support metrics query")
	)
	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(klogFlags)

	// Sync the glog and klog flags.
	pflag.CommandLine.VisitAll(func(f1 *pflag.Flag) {
		f2 := klogFlags.Lookup(f1.Name)
		if f2 != nil {
			value := f1.Value.String()
			if err := f2.Value.Set(value); err != nil {
				util.LogFatalAndExit(err, "failed to set flag")
			}
		}
	})

	pflag.CommandLine.AddGoFlagSet(klogFlags)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	config := &Configuration{
		KubeConfigFile:     *argKubeConfigFile,
		KubeClient:         nil,
		Port:               *argPort,
		Interval:           *argInterval,
		Mode:               *argMode,
		ExitCode:           *argExitCode,
		InternalDNS:        *argInternalDNS,
		ExternalDNS:        *argExternalDNS,
		PodIP:              os.Getenv("POD_IP"),
		HostIP:             os.Getenv("HOST_IP"),
		NodeName:           os.Getenv("NODE_NAME"),
		PodName:            os.Getenv("POD_NAME"),
		ExternalAddress:    *argExternalAddress,
		NetworkMode:        *argNetworkMode,
		EnableMetrics:      *argEnableMetrics,
	}
	if err := config.initKubeClient(); err != nil {
		return nil, err
	}

	podName := os.Getenv("POD_NAME")
	for i := 0; i < 3; i++ {
		pod, err := config.KubeClient.CoreV1().Pods(config.DaemonSetNamespace).Get(context.Background(), podName, metav1.GetOptions{})
		if err != nil {
			klog.Errorf("failed to get self pod %s/%s: %v", config.DaemonSetNamespace, podName, err)
			return nil, err
		}

		if len(pod.Status.PodIPs) != 0 {
			config.PodProtocols = make([]string, len(pod.Status.PodIPs))
			for i, podIP := range pod.Status.PodIPs {
				config.PodProtocols[i] = util.CheckProtocol(podIP.IP)
			}
			break
		}

		if pod.Status.ContainerStatuses[0].Ready {
			util.LogFatalAndExit(nil, "failed to get IPs of Pod %s/%s", config.DaemonSetNamespace, podName)
		}

		klog.Infof("cannot get Pod IPs now, waiting Pod to be ready")
		time.Sleep(time.Second)
	}

	if len(config.PodProtocols) == 0 {
		util.LogFatalAndExit(nil, "failed to get IPs of Pod %s/%s after 3 attempts", config.DaemonSetNamespace, podName)
	}

	klog.Infof("pinger config is %+v", config)
	return config, nil
}

func (config *Configuration) initKubeClient() error {
	var cfg *rest.Config
	var err error
	if config.KubeConfigFile == "" {
		cfg, err = rest.InClusterConfig()
		if err != nil {
			klog.Errorf("use in cluster config failed %v", err)
			return err
		}
	} else {
		cfg, err = clientcmd.BuildConfigFromFlags("", config.KubeConfigFile)
		if err != nil {
			klog.Errorf("use --kubeconfig %s failed %v", config.KubeConfigFile, err)
			return err
		}
	}
	cfg.Timeout = 15 * time.Second
	cfg.QPS = 1000
	cfg.Burst = 2000
	cfg.ContentType = "application/vnd.kubernetes.protobuf"
	cfg.AcceptContentTypes = "application/vnd.kubernetes.protobuf,application/json"
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Errorf("init kubernetes client failed %v", err)
		return err
	}
	config.KubeClient = kubeClient
	return nil
}
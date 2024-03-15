package pinger

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // #nosec
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/klog/v2"

	"github.com/wenwenxiong/network-pinger/pkg/pinger"
	"github.com/wenwenxiong/network-pinger/pkg/util"
	"github.com/wenwenxiong/network-pinger/versions"
)

func CmdMain() {
	defer klog.Flush()

	klog.Infof(versions.String())
	pinger.InitPingerMetrics()
	util.InitKlogMetrics()
	config, err := pinger.ParseFlags()
	if err != nil {
		util.LogFatalAndExit(err, "failed to parse config")
	}
	if config.Mode == "server" && config.EnableMetrics {
		http.Handle("/metrics", promhttp.Handler())

		go func() {
			// conform to Gosec G114
			// https://github.com/securego/gosec#available-rules
			server := &http.Server{
				Addr:              fmt.Sprintf("0.0.0.0:%d", config.Port),
				ReadHeaderTimeout: 3 * time.Second,
			}
			util.LogFatalAndExit(server.ListenAndServe(), "failed to listen and serve on %s", server.Addr)
		}()

	}
	pinger.StartPinger(config)
}

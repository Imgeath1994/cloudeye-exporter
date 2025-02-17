package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Imgeath1994/cloudeye-exporter/collector"
	"github.com/Imgeath1994/cloudeye-exporter/logs"
)

var (
	clientConfig = flag.String("config", "./clouds.yml", "Path to the cloud configuration file")
)

func handler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("services")
	if target == "" {
		http.Error(w, "'target' parameter must be specified", 400)
		return
	}

	targets := strings.Split(target, ",")
	registry := prometheus.NewRegistry()

	logs.Logger.Infof("Start to monitor services: %s", targets)
	exporter := collector.GetMonitoringCollector(targets)
	registry.MustRegister(exporter)
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
	logs.Logger.Infof("End to monitor services: %s", targets)
}

func epHandler(w http.ResponseWriter, r *http.Request) {
	epsInfo, err := collector.GetEPSInfo()
	if err != nil {
		http.Error(w, fmt.Sprintf("get eps info error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte(epsInfo))
	if err != nil {
		logs.Logger.Errorf("Response to caller error: %s", err.Error())
	}
}

func main() {
	flag.Parse()
	logs.InitLog()
	err := collector.InitCloudConf(*clientConfig)
	if err != nil {
		logs.Logger.Error("Init Cloud Config From File error: ", err.Error())
		logs.FlushLogAndExit(1)
	}
	err = collector.InitMetricConf()
	if err != nil {
		logs.Logger.Error("Init metric Config error: ", err.Error())
		logs.FlushLogAndExit(1)
	}

	http.HandleFunc(collector.CloudConf.Global.MetricPath, handler)
	http.HandleFunc(collector.CloudConf.Global.EpsInfoPath, epHandler)
	server := &http.Server{
		Addr:         collector.CloudConf.Global.Port,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	logs.Logger.Info("Start server at ", collector.CloudConf.Global.Port)
	if err := server.ListenAndServe(); err != nil {
		logs.Logger.Errorf("Error occur when start server %s", err.Error())
		logs.FlushLogAndExit(1)
	}
}

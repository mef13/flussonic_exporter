/*
 *    Copyright 2020 Yury Makarov
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 *
 */

package main

import (
	"flag"
	"fmt"
	"github.com/mef13/flussonic_exporter/collector"
	"github.com/mef13/flussonic_exporter/flussonic"
	"github.com/mef13/flussonic_exporter/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

//init config from file
func initViper(confPath string) {
	if confPath == "" {
		viper.SetConfigName("settings")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("/etc/flussonic_exporter")
		viper.AddConfigPath("./conf")
		viper.AddConfigPath(".")
	} else {
		viper.SetConfigFile(confPath)
	}
	err := viper.ReadInConfig() // Find and read the config file
	viper.SetDefault("log-path", "/var/log/flussonic_exporter")
	viper.SetDefault("log-level", "info")
	viper.SetDefault("listen-address", ":9113")
	viper.SetDefault("metrics-path", "/metrics")
	viper.SetDefault("exporter-metrics", true)
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

var (
	config = flag.String("config", "", "Path to config file")
)

func newHandler(includeExporterMetrics bool, c *collector.FlussonicCollector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		registry := prometheus.NewRegistry()
		if err := registry.Register(c); err != nil {
			logger.Error("Couldn't register collector", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			if _, err = w.Write([]byte(fmt.Sprintf("Couldn't register collector: %s", err))); err != nil {
				logger.Warn("Couldn't write response", zap.Error(err))
			}
			return
		}
		gatherers := prometheus.Gatherers{
			registry,
		}
		if includeExporterMetrics {
			gatherers = append(gatherers, prometheus.DefaultGatherer)
		}

		h := promhttp.InstrumentMetricHandler(
			registry,
			promhttp.HandlerFor(gatherers,
				promhttp.HandlerOpts{
					ErrorLog:      log.NewErrorLogger(),
					ErrorHandling: promhttp.ContinueOnError,
				}),
		)
		h.ServeHTTP(w, r)

	}
}

func main() {
	flag.CommandLine.SetOutput(os.Stdout)
	flag.Usage = usage
	flag.Parse()
	initViper(*config)
	logger.InitLogger(viper.GetString("log-path"), viper.GetString("log-level"), viper.GetString("sentryDSN"), version)
	logger.Info("Starting Flussonic exporter.", zap.String("version", version))
	defer logger.Sync()

	fluss, err := flussonic.ParseConfig(viper.GetViper(), "flussonics")
	if err != nil {
		logger.Error("error parse flussonics section in config", zap.Error(err))
		os.Exit(1)
	}

	flussonicCollector := collector.NewCollector()

	c := cron.New()
	c.Start()

	for _, flus := range fluss {
		jobName := fmt.Sprintf("Scrape %s", flus.InstanceName)
		funcJob := flussonicCollector.GetCronJob(*flus)
		job := cron.NewChain(cron.SkipIfStillRunning(logger.GetLoggerForCron(jobName))).Then(funcJob)
		duration := fmt.Sprintf("@every %s", flus.ScrapeInterval)
		id, err := c.AddJob(duration, job)
		if err != nil {
			logger.Error(fmt.Sprintf("error register task %s", jobName), zap.Error(err))
			os.Exit(1)
		}
		logger.Info(fmt.Sprintf("register task %s. Next run: %s", jobName, c.Entry(id).Next.Format(time.RubyDate)))

	}

	http.Handle(viper.GetString("metrics-path"), newHandler(viper.GetBool("exporter-metrics"), flussonicCollector))
	server := &http.Server{Addr: viper.GetString("listen-address")}
	logger.Info(fmt.Sprintf("listening on %s", viper.GetString("listen-address")))
	if err := server.ListenAndServe(); err != nil {
		logger.Error("error listen", zap.Error(err))
		os.Exit(1)
	}
}

func usage() {
	const s = `
flussonic_exporter is Prometheus exporter for flussonic.

See the docs at https://github.com/mef13/flussonic_exporter .
`

	f := flag.CommandLine.Output()
	fmt.Fprintf(f, "%s\n", s)
	flag.PrintDefaults()
}

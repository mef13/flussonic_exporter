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

package collector

import (
	"github.com/mef13/flussonic_exporter/flussonic"
	"github.com/mef13/flussonic_exporter/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"sync"
	"time"
)

const namespace = "flussonic"

var (
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, `scrape`, `collector_duration_seconds`),
		`flussonic_exporter: Duration of a collector scrape.`,
		[]string{`server`},
		nil,
	)
	scrapeSuccessDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, `scrape`, `collector_success`),
		`flussonic_exporter: Whether a collector succeeded.`,
		[]string{`server`},
		nil,
	)

	totalClientsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, `clients`, `total`),
		`flussonic_exporter: Total clients count.`,
		[]string{`server`},
		prometheus.Labels{"type": "total"},
	)
	requestDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, `scrape`, `api_request_duration_sec`),
		`flussonic_exporter: API request duration.`,
		[]string{`server`, `url`},
		nil,
	)
)

type FlussonicCollector struct {
	cache map[string]*flussonicCollectorCache
}

type flussonicCollectorCache struct {
	cache []prometheus.Metric
	sync  sync.RWMutex
}

// Describe implements the prometheus.Collector interface.
func (c *FlussonicCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeSuccessDesc
	ch <- scrapeDurationDesc
}

// Collect implements the prometheus.Collector interface.
func (c *FlussonicCollector) Collect(ch chan<- prometheus.Metric) {
	for flussonicUrl := range c.cache {
		send(ch, c.cache[flussonicUrl])
	}
}

func send(ch chan<- prometheus.Metric, cache *flussonicCollectorCache) {
	cache.sync.RLock()
	defer cache.sync.RUnlock()
	for _, metric := range cache.cache {
		ch <- metric
	}
}

func NewCollector() *FlussonicCollector {
	return &FlussonicCollector{cache: make(map[string]*flussonicCollectorCache)}
}

func (c *flussonicCollectorCache) addMetric(m prometheus.Metric) {
	c.cache = append(c.cache, m)
}

func (c *FlussonicCollector) save(flussonicUrl string, cache *flussonicCollectorCache) {
	_, ok := c.cache[flussonicUrl]
	if ok {
		c.cache[flussonicUrl].sync.Lock()
		defer c.cache[flussonicUrl].sync.Unlock()
	}
	c.cache[flussonicUrl] = cache
}

func (c *FlussonicCollector) failScrape(flussConf flussonic.Flussonic, startTime time.Time) {
	cache := &flussonicCollectorCache{}
	cache.addMetric(prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, float64(0),
		flussConf.InstanceName))
	duration := time.Since(startTime)
	cache.addMetric(prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(),
		flussConf.InstanceName))
	c.save(flussConf.Url.String(), cache)
}

func (c *FlussonicCollector) Scrape(flussConf flussonic.Flussonic) {
	logger.Debug("start scrapping", zap.String("instance", flussConf.InstanceName))
	startTime := time.Now()
	cache := &flussonicCollectorCache{}

	//get metrics
	serv, err := flussConf.GetServer()
	if err != nil {
		logger.Error("error scrape from flussonic api",
			zap.String("server", flussConf.Url.String()), zap.String("method", "GetServer"), zap.Error(err))
		c.failScrape(flussConf, startTime)
		return
	}

	//add metrics to cache
	cache.addMetric(prometheus.MustNewConstMetric(
		totalClientsDesc,
		prometheus.GaugeValue,
		serv.TotalClients,
		flussConf.InstanceName,
	))
	cache.addMetric(prometheus.MustNewConstMetric(
		requestDurationDesc,
		prometheus.GaugeValue,
		serv.RequestDuration,
		flussConf.InstanceName,
		serv.Url,
	))

	//end scrape & save cache
	cache.addMetric(prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, float64(1),
		flussConf.InstanceName))
	duration := time.Since(startTime)
	cache.addMetric(prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(),
		flussConf.InstanceName))
	c.save(flussConf.Url.String(), cache)
}

// FuncJob is a wrapper that turns a func() into a cron.Job
type FuncJob struct {
	f         func(flussConf flussonic.Flussonic)
	flussConf flussonic.Flussonic
}

func (f FuncJob) Run() { f.f(f.flussConf) }

func (c *FlussonicCollector) GetCronJob(flussConf flussonic.Flussonic) cron.Job {
	return FuncJob{
		f:         c.Scrape,
		flussConf: flussConf,
	}
}

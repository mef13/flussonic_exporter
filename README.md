# Flussonic exporter
[![Latest Version](https://img.shields.io/github/release/mef13/flussonic_exporter.svg?maxAge=8600)](https://github.com/mef13/flussonic_exporter/releases/latest)
[![License](https://img.shields.io/github/license/janeczku/rancher-letsencrypt.svg?maxAge=8600)](https://github.com/mef13/flussonic_exporter/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/mef13/flussonic_exporter)](https://goreportcard.com/report/github.com/mef13/flussonic_exporter)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmef13%2Fflussonic_exporter.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fmef13%2Fflussonic_exporter?ref=badge_shield)

Prometheus exporter for Flussonic media server

## What is collecting
* Server
    * Total clients count
    * Dvr clients count
* Streams 
    * Bitrate
    * Alive
    * Retry count
    * Input error rate
    * Total clients count
    * Dvr clients count
    * Tracks count

## Config
Specify config file by `-config` flag.
```shell script
./flussonic_exporter -config /etc/flussonic_exporter/settings.yaml
```

settings.yaml 
```yaml
log-path: "/var/log/flussonic_exporter"
log-level: info               
listen-address: ":9113"
metrics-path: "/metrics"
exporter-metrics: false
flussonics:
  - user: "api_user"
    password: "pass"
    url: "http://example.com:8081"
    scrape-interval: "60s"
    instance-name: "my-flussonic"
```

## Prometheus
```
  - job_name: 'flussonic'
    metrics_path: /metrics
    scrape_interval: 60s
    static_configs:
            - targets: [ 'localhost:9113']

```

## Useful alerts
Server api not response(Flussonic down):
```
  - alert: FlussonicServerNotResponse
    expr: avg_over_time(flussonic_scrape_collector_success[5m]) * 100 < 50
    labels:
      severity: critical
    annotations:
      summary: "Flussonic api not response (server {{ $labels.server }})"
      description: "Flussonic server '{{ $labels.server }}' not response."
```

Stream down more than 5 minutes:
```
  - alert: FlussonicStreamDown
    expr: flussonic_stream_retry_count > 20
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "Flussonic stream down (server {{ $labels.server }})"
      description: "Flussonic stream '{{ $labels.name }}' down. Server {{ $labels.server }}"
```

The number of tracks on a stream is more than 2:
```
  - alert: FlussonicStreamTracksCount
    expr: flussonic_stream_tracks_count > 2
    labels:
      severity: warning
    annotations:
      summary: "Flussonic stream tracks count mismatch (server {{ $labels.server }})"
      description: "Flussonic stream '{{ $labels.name }}' tracks count mismatch(tracks count = {{ $value }}). Server {{ $labels.server }}"

``` 

## Community
* [gitter](https://gitter.im/flussonic_exporter/community)

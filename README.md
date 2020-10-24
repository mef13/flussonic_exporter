# Flussonic exporter
[![Latest Version](https://img.shields.io/github/release/mef13/flussonic_exporter.svg?maxAge=8600)]()
[![License](https://img.shields.io/github/license/janeczku/rancher-letsencrypt.svg?maxAge=8600)]()
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmef13%2Fflussonic_exporter.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fmef13%2Fflussonic_exporter?ref=badge_shield)

Prometheus exporter for Flussonic media server

## Collected metrics
* Total and dvr clients count
* Streams bitrate, alive, retry count, input error rate

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

## Community
* [gitter](https://gitter.im/flussonic_exporter/community)


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmef13%2Fflussonic_exporter.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fmef13%2Fflussonic_exporter?ref=badge_large)
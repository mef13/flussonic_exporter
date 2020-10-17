# Flussonic exporter
Prometheus exporter for Flussonic media server

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

# push_exporter

Prometheus exporter for push mode application. The api and behavior are similar to the agent of open-falcon. After collecting data, export to Prometheus.


## How to run

First, compile and launch the server:
```bash
go build
./push_exporter -p 1999
```

Then, visit the home and metrics page using curl or a browser:
```bash
curl http://localhost:1999/
curl http://localhost:1999/metrics
```

## How to push data

```cURL
curl http://localhost:1999/push \
  -H "Content-Type: application/json" \
  -d '{
    "metric": "send_frequency",
    "endpoint": "localhost",
    "job": "mail",
    "tags": "a=b,c=d",
    "type": "G",
    "step": 15,
    "value": 0.005
  }'

Of which:
metric: metric name
endpoint: server where the service is located
job: service name
tags: additional tags corresponding to the metric (can be empty)
type: metric type (G for Gauge, C for Counter)
step: collection interval for the metric, in seconds
value: value of the metric (float type)
```

In addition to pushing a single JSON object, it is also possible to push an array of JSON objects, such as: [{...},{...}]

# CMDB Lite Monitoring and Troubleshooting Guide

This document provides comprehensive guidance on monitoring, troubleshooting, and maintaining the health of CMDB Lite in production environments.

## Table of Contents

- [Monitoring Overview](#monitoring-overview)
- [Key Metrics to Monitor](#key-metrics-to-monitor)
  - [Application Metrics](#application-metrics)
  - [Database Metrics](#database-metrics)
  - [System Metrics](#system-metrics)
- [Monitoring Tools](#monitoring-tools)
  - [Prometheus and Grafana](#prometheus-and-grafana)
  - [ELK Stack](#elk-stack)
  - [Datadog](#datadog)
  - [Built-in Monitoring](#built-in-monitoring)
- [Alerting](#alerting)
  - [Alerting Best Practices](#alerting-best-practices)
  - [Common Alerts](#common-alerts)
  - [Alerting Configuration](#alerting-configuration)
- [Logging](#logging)
  - [Application Logging](#application-logging)
  - [Database Logging](#database-logging)
  - [Log Aggregation](#log-aggregation)
- [Troubleshooting](#troubleshooting)
  - [Common Issues](#common-issues)
  - [Diagnostic Tools](#diagnostic-tools)
  - [Troubleshooting Procedures](#troubleshooting-procedures)
- [Performance Tuning](#performance-tuning)
  - [Application Performance](#application-performance)
  - [Database Performance](#database-performance)
  - [Infrastructure Performance](#infrastructure-performance)
- [Health Checks](#health-checks)
  - [Application Health Checks](#application-health-checks)
  - [Database Health Checks](#database-health-checks)
  - [Infrastructure Health Checks](#infrastructure-health-checks)

## Monitoring Overview

Effective monitoring is essential for maintaining the health, performance, and reliability of CMDB Lite in production environments. A comprehensive monitoring strategy should include:

1. **Metrics Collection**: Collecting and storing metrics from various components
2. **Visualization**: Creating dashboards to visualize metrics and trends
3. **Alerting**: Setting up alerts to notify operators of potential issues
4. **Logging**: Collecting and analyzing logs for troubleshooting
5. **Health Checks**: Implementing health checks to verify component status

### Monitoring Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend       │    │   Database      │
│   (Vue.js)      │    │   (Go)          │    │   (PostgreSQL)  │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          │ Metrics              │ Metrics              │ Metrics
          │ Logs                 │ Logs                 │ Logs
          │ Health Checks        │ Health Checks        │ Health Checks
          ▼                      ▼                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Monitoring System                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │ Prometheus  │  │   ELK       │  │   Datadog   │              │
│  │             │  │   Stack     │  │             │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
│                                                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │   Grafana   │  │   Kibana    │  │ Datadog UI  │              │
│  │             │  │             │  │             │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
│                                                                 │
│  ┌─────────────┐                                                │
│  │ Alertmanager│                                                │
│  │             │                                                │
│  └─────────────┘                                                │
└─────────────────────────────────────────────────────────────────┘
          │
          │ Alerts
          ▼
┌─────────────────┐
│  Notification   │
│  Channels       │
│  (Email, Slack, │
│   PagerDuty)    │
└─────────────────┘
```

## Key Metrics to Monitor

### Application Metrics

#### Frontend Metrics

- **Page Load Time**: Time taken to load the application
- **First Contentful Paint (FCP)**: Time when the first content is painted on the screen
- **Largest Contentful Paint (LCP)**: Time when the largest content is painted on the screen
- **First Input Delay (FID)**: Time from user interaction to browser response
- **Cumulative Layout Shift (CLS)**: Measure of visual stability
- **JavaScript Errors**: Number of JavaScript errors
- **API Request Time**: Time taken for API requests
- **API Error Rate**: Percentage of API requests that result in errors

#### Backend Metrics

- **Request Rate**: Number of requests per second
- **Response Time**: Average and percentiles (P50, P90, P95, P99) of response times
- **Error Rate**: Percentage of requests that result in errors (4xx, 5xx)
- **Memory Usage**: Application memory consumption
- **CPU Usage**: Application CPU consumption
- **Goroutine Count**: Number of goroutines (for Go backend)
- **Database Connection Pool Usage**: Number of active and idle database connections
- **JWT Token Validation Time**: Time taken to validate JWT tokens

### Database Metrics

- **Connection Count**: Number of active database connections
- **Connection Pool Usage**: Percentage of connection pool in use
- **Query Performance**: Average query time, number of slow queries
- **Database Size**: Size of the database and growth rate
- **Table Sizes**: Sizes of individual tables
- **Index Usage**: Usage statistics for indexes
- **Lock Waits**: Number and duration of lock waits
- **Checkpoint Activity**: Frequency and duration of checkpoints
- **WAL (Write-Ahead Log) Activity**: WAL generation rate and archival status
- **Replication Lag**: For replicated databases, the lag between primary and replicas

### System Metrics

- **CPU Usage**: System CPU usage (user, system, iowait, idle)
- **Memory Usage**: System memory usage (used, free, cached, buffers)
- **Swap Usage**: Swap memory usage
- **Disk Usage**: Available disk space and I/O operations per second (IOPS)
- **Disk Latency**: Average disk read/write latency
- **Network Traffic**: Incoming and outgoing network traffic
- **Network Latency**: Network latency between components
- **System Load**: System load average (1min, 5min, 15min)
- **Context Switches**: Number of context switches per second
- **Process Count**: Number of running processes

## Monitoring Tools

### Prometheus and Grafana

Prometheus and Grafana are popular open-source tools for monitoring and visualization.

#### Setting up Prometheus

1. Install Prometheus:

```bash
# For Ubuntu/Debian
sudo apt update
sudo apt install prometheus

# For CentOS/RHEL
sudo yum install prometheus
```

2. Configure Prometheus to scrape CMDB Lite metrics:

```yaml
# /etc/prometheus/prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'cmdb-lite-backend'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 15s

  - job_name: 'cmdb-lite-frontend'
    static_configs:
      - targets: ['localhost:3000']
    metrics_path: '/metrics'
    scrape_interval: 15s

  - job_name: 'postgres'
    static_configs:
      - targets: ['localhost:9187']  # postgres_exporter
    scrape_interval: 15s
```

3. Start Prometheus:

```bash
sudo systemctl start prometheus
sudo systemctl enable prometheus
```

#### Setting up Grafana

1. Install Grafana:

```bash
# For Ubuntu/Debian
sudo apt install -y apt-transport-https
sudo apt install -y software-properties-common wget
wget -q -O - https://packages.grafana.com/gpg.key | sudo apt-key add -
echo "deb https://packages.grafana.com/oss/deb stable main" | sudo tee -a /etc/apt/sources.list.d/grafana.list
sudo apt update
sudo apt install grafana

# For CentOS/RHEL
sudo yum install -y https://dl.grafana.com/oss/release/grafana-7.5.7-1.x86_64.rpm
```

2. Start Grafana:

```bash
sudo systemctl start grafana-server
sudo systemctl enable grafana-server
```

3. Configure Prometheus as a data source in Grafana:
   - Open Grafana at http://localhost:3000
   - Log in with the default credentials (admin/admin)
   - Navigate to Configuration > Data Sources
   - Add Prometheus as a data source with URL http://localhost:9090

4. Import or create dashboards for CMDB Lite metrics.

#### Example Grafana Dashboard

Here's an example of a Grafana dashboard for monitoring CMDB Lite:

```json
{
  "dashboard": {
    "id": null,
    "title": "CMDB Lite Overview",
    "tags": ["cmdb-lite"],
    "timezone": "browser",
    "panels": [
      {
        "id": 1,
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{method}} {{status}}"
          }
        ],
        "yAxes": [
          {
            "label": "Requests per second"
          }
        ]
      },
      {
        "id": 2,
        "title": "Response Time",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "P95"
          },
          {
            "expr": "histogram_quantile(0.50, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "P50"
          }
        ],
        "yAxes": [
          {
            "label": "Seconds"
          }
        ]
      },
      {
        "id": 3,
        "title": "Error Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m]) * 100",
            "legendFormat": "5xx Errors"
          },
          {
            "expr": "rate(http_requests_total{status=~"4.."}[5m]) / rate(http_requests_total[5m]) * 100",
            "legendFormat": "4xx Errors"
          }
        ],
        "yAxes": [
          {
            "label": "Percentage"
          }
        ]
      },
      {
        "id": 4,
        "title": "Database Connections",
        "type": "graph",
        "targets": [
          {
            "expr": "pg_stat_database_numbackends",
            "legendFormat": "{{datname}}"
          }
        ],
        "yAxes": [
          {
            "label": "Connections"
          }
        ]
      }
    ],
    "time": {
      "from": "now-1h",
      "to": "now"
    },
    "refresh": "30s"
  }
}
```

### ELK Stack

The ELK Stack (Elasticsearch, Logstash, Kibana) is a popular choice for log aggregation and analysis.

#### Setting up Elasticsearch

1. Install Elasticsearch:

```bash
# For Ubuntu/Debian
wget -qO - https://artifacts.elastic.co/GPG-KEY-elasticsearch | sudo apt-key add -
echo "deb https://artifacts.elastic.co/packages/7.x/apt stable main" | sudo tee -a /etc/apt/sources.list.d/elastic-7.x.list
sudo apt update
sudo apt install elasticsearch

# For CentOS/RHEL
sudo rpm --import https://artifacts.elastic.co/GPG-KEY-elasticsearch
echo "[elasticsearch-7.x]
name=Elasticsearch repository for 7.x packages
baseurl=https://artifacts.elastic.co/packages/7.x/yum
gpgcheck=1
gpgkey=https://artifacts.elastic.co/GPG-KEY-elasticsearch
enabled=1
autorefresh=1
type=rpm-md" | sudo tee /etc/yum.repos.d/elasticsearch.repo
sudo yum install elasticsearch
```

2. Configure Elasticsearch:

```yaml
# /etc/elasticsearch/elasticsearch.yml
cluster.name: cmdb-lite
node.name: node-1
network.host: 0.0.0.0
discovery.type: single-node
```

3. Start Elasticsearch:

```bash
sudo systemctl start elasticsearch
sudo systemctl enable elasticsearch
```

#### Setting up Logstash

1. Install Logstash:

```bash
# For Ubuntu/Debian
sudo apt install logstash

# For CentOS/RHEL
sudo yum install logstash
```

2. Configure Logstash to process CMDB Lite logs:

```ruby
# /etc/logstash/conf.d/cmdb-lite.conf
input {
  file {
    path => "/var/log/cmdb-lite/backend.log"
    start_position => "beginning"
    codec => "json"
  }
  
  file {
    path => "/var/log/cmdb-lite/frontend.log"
    start_position => "beginning"
    codec => "json"
  }
}

filter {
  if [level] == "error" {
    mutate {
      add_tag => ["error"]
    }
  }
  
  date {
    match => [ "timestamp", "ISO8601" ]
  }
}

output {
  elasticsearch {
    hosts => ["localhost:9200"]
    index => "cmdb-lite-%{+YYYY.MM.dd}"
  }
}
```

3. Start Logstash:

```bash
sudo systemctl start logstash
sudo systemctl enable logstash
```

#### Setting up Kibana

1. Install Kibana:

```bash
# For Ubuntu/Debian
sudo apt install kibana

# For CentOS/RHEL
sudo yum install kibana
```

2. Configure Kibana:

```yaml
# /etc/kibana/kibana.yml
server.host: "0.0.0.0"
elasticsearch.hosts: ["http://localhost:9200"]
```

3. Start Kibana:

```bash
sudo systemctl start kibana
sudo systemctl enable kibana
```

4. Access Kibana at http://localhost:5601 and create index patterns for CMDB Lite logs.

### Datadog

Datadog is a cloud-based monitoring service that provides infrastructure monitoring, application performance monitoring (APM), and log management.

#### Setting up Datadog

1. Sign up for a Datadog account.
2. Install the Datadog Agent:

```bash
# For Ubuntu/Debian
sudo apt update
sudo apt install -y apt-transport-https
sudo sh -c "echo 'deb [signed-by=/usr/share/keyrings/datadog-archive-keyring.gpg] https://apt.datadoghq.com/ stable 7' > /etc/apt/sources.list.d/datadog.list"
sudo apt install -y datadog-agent

# For CentOS/RHEL
sudo yum install -y datadog-agent
```

3. Configure the Datadog Agent:

```yaml
# /etc/datadog-agent/datadog.yaml
api_key: YOUR_DATADOG_API_KEY
site: datadoghq.com
```

4. Enable integrations for PostgreSQL, Nginx, and other components:

```yaml
# /etc/datadog-agent/conf.d/postgres.d/conf.yaml
init_config:

instances:
  - host: localhost
    port: 5432
    username: datadog
    password: DATADOG_POSTGRES_PASSWORD
    dbname: cmdb_lite
    ssl: false
    tags:
      - service:cmdb-lite
      - env:production

# /etc/datadog-agent/conf.d/nginx.d/conf.yaml
init_config:

instances:
  - nginx_status_url: http://localhost/nginx_status
    tags:
      - service:cmdb-lite
      - env:production
```

5. Start the Datadog Agent:

```bash
sudo systemctl start datadog-agent
sudo systemctl enable datadog-agent
```

6. Configure APM tracing for the backend:

```yaml
# /etc/datadog-agent/conf.d/go.d/conf.yaml
init_config:

instances:
  - env: production
    service: cmdb-lite-backend
    agent_port: 8126
```

7. Configure the backend to send traces to Datadog:

```go
import (
    "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
    tracer.Start(
        tracer.WithService("cmdb-lite-backend"),
        tracer.WithEnv("production"),
    )
    defer tracer.Stop()
    
    // Your application code
}
```

### Built-in Monitoring

CMDB Lite includes built-in monitoring endpoints that can be used with any monitoring system:

#### Backend Metrics Endpoint

The backend exposes metrics at `/metrics` in Prometheus format:

```bash
curl http://localhost:8080/metrics
```

Example output:

```
# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",path="/api/v1/cis",status="200"} 42
http_requests_total{method="POST",path="/api/v1/cis",status="201"} 5

# HELP http_request_duration_seconds Histogram of HTTP request durations
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{method="GET",path="/api/v1/cis",status="200",le="0.005"} 10
http_request_duration_seconds_bucket{method="GET",path="/api/v1/cis",status="200",le="0.01"} 20
http_request_duration_seconds_bucket{method="GET",path="/api/v1/cis",status="200",le="0.025"} 30
http_request_duration_seconds_bucket{method="GET",path="/api/v1/cis",status="200",le="0.05"} 35
http_request_duration_seconds_bucket{method="GET",path="/api/v1/cis",status="200",le="0.1"} 40
http_request_duration_seconds_bucket{method="GET",path="/api/v1/cis",status="200",le="0.25"} 42
http_request_duration_seconds_bucket{method="GET",path="/api/v1/cis",status="200",le="0.5"} 42
http_request_duration_seconds_bucket{method="GET",path="/api/v1/cis",status="200",le="1"} 42
http_request_duration_seconds_bucket{method="GET",path="/api/v1/cis",status="200",le="2.5"} 42
http_request_duration_seconds_bucket{method="GET",path="/api/v1/cis",status="200",le="5"} 42
http_request_duration_seconds_bucket{method="GET",path="/api/v1/cis",status="200",le="10"} 42
http_request_duration_seconds_bucket{method="GET",path="/api/v1/cis",status="200",le="+Inf"} 42
http_request_duration_seconds_sum{method="GET",path="/api/v1/cis",status="200"} 1.23
http_request_duration_seconds_count{method="GET",path="/api/v1/cis",status="200"} 42

# HELP go_goroutines Number of goroutines that currently exist
# TYPE go_goroutines gauge
go_goroutines 15

# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 2.54321e+06
```

#### Health Check Endpoint

The backend provides a health check endpoint at `/health`:

```bash
curl http://localhost:8080/health
```

Example output:

```json
{
  "status": "healthy",
  "timestamp": "2023-01-01T12:00:00Z",
  "version": "1.0.0",
  "components": {
    "database": {
      "status": "healthy",
      "details": {
        "connection_time_ms": 5,
        "version": "13.4"
      }
    },
    "cache": {
      "status": "healthy",
      "details": {
        "connection_time_ms": 1
      }
    }
  }
}
```

## Alerting

Alerting is a critical component of monitoring, allowing operators to be notified of potential issues before they impact users.

### Alerting Best Practices

1. **Set Meaningful Thresholds**: Set alert thresholds based on normal operating conditions and user experience requirements.
2. **Avoid Alert Fatigue**: Only alert on issues that require immediate attention.
3. **Use Escalation Policies**: Implement escalation policies for critical alerts.
4. **Include Context**: Include relevant context in alert notifications to help with troubleshooting.
5. **Test Alerts**: Regularly test alerting mechanisms to ensure they work correctly.
6. **Review and Adjust**: Regularly review and adjust alert thresholds and rules based on changing conditions.

### Common Alerts

#### Application Alerts

1. **High Error Rate**:
   - Condition: Error rate exceeds 5% for 5 minutes
   - Severity: Critical
   - Action: Investigate and fix the cause of errors

2. **High Response Time**:
   - Condition: 95th percentile response time exceeds 1 second for 5 minutes
   - Severity: Warning
   - Action: Investigate performance bottlenecks

3. **Service Unavailable**:
   - Condition: Service is unavailable for more than 1 minute
   - Severity: Critical
   - Action: Restart service and investigate cause

4. **High Memory Usage**:
   - Condition: Memory usage exceeds 90% for 5 minutes
   - Severity: Warning
   - Action: Investigate memory leaks or increase memory allocation

5. **High CPU Usage**:
   - Condition: CPU usage exceeds 90% for 5 minutes
   - Severity: Warning
   - Action: Investigate CPU-intensive operations or scale up

#### Database Alerts

1. **Database Connection Issues**:
   - Condition: Database connection errors exceed 1% for 5 minutes
   - Severity: Critical
   - Action: Investigate database connectivity and performance

2. **Slow Queries**:
   - Condition: Number of slow queries exceeds 10 per minute
   - Severity: Warning
   - Action: Optimize queries and indexes

3. **High Database Connections**:
   - Condition: Database connections exceed 80% of maximum for 5 minutes
   - Severity: Warning
   - Action: Investigate connection leaks or increase connection pool size

4. **Database Replication Lag**:
   - Condition: Replication lag exceeds 30 seconds
   - Severity: Critical
   - Action: Investigate replication issues

#### System Alerts

1. **Disk Space Low**:
   - Condition: Disk space is below 20%
   - Severity: Warning
   - Action: Clean up disk space or increase storage capacity

2. **High System Load**:
   - Condition: System load average exceeds number of CPU cores for 5 minutes
   - Severity: Warning
   - Action: Investigate resource-intensive processes

3. **Network Issues**:
   - Condition: Network latency exceeds 100ms or packet loss exceeds 1%
   - Severity: Warning
   - Action: Investigate network connectivity issues

### Alerting Configuration

#### Prometheus Alertmanager

Alertmanager is the alerting component of the Prometheus ecosystem.

1. Install Alertmanager:

```bash
# For Ubuntu/Debian
sudo apt install prometheus-alertmanager

# For CentOS/RHEL
sudo yum install prometheus-alertmanager
```

2. Configure Alertmanager:

```yaml
# /etc/prometheus/alertmanager.yml
global:
  smtp_smarthost: 'localhost:587'
  smtp_from: 'alerts@your-domain.com'
  smtp_auth_username: 'alerts@your-domain.com'
  smtp_auth_password: 'your-smtp-password'

route:
  group_by: ['alertname', 'severity']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 1h
  receiver: 'web.hook'

receivers:
- name: 'web.hook'
  email_configs:
  - to: 'ops-team@your-domain.com'
    subject: 'CMDB Lite Alert: {{ .GroupLabels.alertname }}'
    body: |
      {{ range .Alerts }}
      Alert: {{ .Annotations.summary }}
      Description: {{ .Annotations.description }}
      Labels: {{ .Labels }}
      {{ end }}
  webhook_configs:
  - url: 'http://127.0.0.1:5001/'
```

3. Configure Prometheus to use Alertmanager:

```yaml
# /etc/prometheus/prometheus.yml
alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - localhost:9093
```

4. Create alert rules:

```yaml
# /etc/prometheus/alerts.yml
groups:
  - name: cmdb-lite
    rules:
    - alert: HighErrorRate
      expr: rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m]) * 100 > 5
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "High error rate detected"
        description: "CMDB Lite error rate is {{ $value }}% for more than 5 minutes"
    
    - alert: HighResponseTime
      expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "High response time detected"
        description: "CMDB Lite 95th percentile response time is {{ $value }} seconds for more than 5 minutes"
    
    - alert: ServiceUnavailable
      expr: up{job="cmdb-lite-backend"} == 0
      for: 1m
      labels:
        severity: critical
      annotations:
        summary: "Service unavailable"
        description: "CMDB Lite backend service is unavailable"
    
    - alert: HighMemoryUsage
      expr: process_resident_memory_bytes{job="cmdb-lite-backend"} / process_virtual_memory_max_bytes{job="cmdb-lite-backend"} * 100 > 90
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "High memory usage detected"
        description: "CMDB Lite backend memory usage is {{ $value }}% for more than 5 minutes"
```

5. Start Alertmanager:

```bash
sudo systemctl start prometheus-alertmanager
sudo systemctl enable prometheus-alertmanager
```

#### Datadog Monitors

Datadog provides a web-based interface for creating monitors:

1. Navigate to Monitors > New Monitor in the Datadog UI.
2. Select the monitor type (e.g., Metric, APM, Log).
3. Configure the monitor definition:
   - For a metric monitor, select the metric and set the threshold.
   - For an APM monitor, select the service and operation.
   - For a log monitor, define the search query and pattern.
4. Configure the notification settings:
   - Set the notification message.
   - Select notification channels (email, Slack, PagerDuty, etc.).
   - Set the notification frequency and re-notification settings.
5. Save the monitor.

Example Datadog monitor configuration for high error rate:

```yaml
name: CMDB Lite High Error Rate
type: metric alert
query: avg:last_5m):sum:trace.http.request.errors{service:cmdb-lite-backend} / sum:trace.http.request.hits{service:cmdb-lite-backend} * 100 > 5
message: "CMDB Lite error rate is {{value}}% for more than 5 minutes"
tags: ["service:cmdb-lite", "env:production"]
options:
  notify_no_data: false
  no_data_timeframe: null
  renotify_interval: 60
  timeout_h: 0
  include_tags: true
  require_full_window: true
  new_host_delay: 300
  evaluation_delay: 900
  thresholds:
    critical: 5
    warning: 2
```

## Logging

Logging is essential for troubleshooting issues and understanding application behavior.

### Application Logging

#### Backend Logging

The backend uses structured logging with the following fields:

- **level**: Log level (debug, info, warn, error)
- **timestamp**: Timestamp of the log entry
- **request_id**: Unique identifier for the request
- **user_id**: User ID (if available)
- **method**: HTTP method
- **path**: HTTP path
- **status_code**: HTTP status code
- **duration_ms**: Request duration in milliseconds
- **message**: Log message
- **error**: Error details (if applicable)

Example log entry:

```json
{
  "level": "info",
  "timestamp": "2023-01-01T12:00:00Z",
  "request_id": "abc123",
  "user_id": "user123",
  "method": "GET",
  "path": "/api/v1/cis",
  "status_code": 200,
  "duration_ms": 45,
  "message": "Request completed successfully"
}
```

#### Frontend Logging

The frontend logs errors and warnings to the browser console, and can be configured to send logs to a logging service:

```javascript
// Configure logging
const logger = {
  debug: (message, data) => {
    if (process.env.NODE_ENV === 'development') {
      console.debug(message, data);
    }
    // Send to logging service
    sendToLoggingService('debug', message, data);
  },
  info: (message, data) => {
    console.info(message, data);
    // Send to logging service
    sendToLoggingService('info', message, data);
  },
  warn: (message, data) => {
    console.warn(message, data);
    // Send to logging service
    sendToLoggingService('warn', message, data);
  },
  error: (message, error) => {
    console.error(message, error);
    // Send to logging service
    sendToLoggingService('error', message, { error: error.message, stack: error.stack });
  }
};

function sendToLoggingService(level, message, data) {
  if (process.env.NODE_ENV === 'production') {
    // Send to logging service
    fetch('/api/v1/logs', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        level,
        message,
        data,
        timestamp: new Date().toISOString(),
        url: window.location.href,
        userAgent: navigator.userAgent,
      }),
    });
  }
}
```

### Database Logging

PostgreSQL provides extensive logging capabilities:

```ini
# /etc/postgresql/13/main/postgresql.conf
# Log all statements
log_statement = 'all'

# Log duration of statements
log_duration = on

# Log queries that take longer than 1 second
log_min_duration_statement = 1000

# Log checkpoints
log_checkpoints = on

# Log connections
log_connections = on
log_disconnections = on

# Log lock waits
log_lock_waits = on

# Log temporary files
log_temp_files = 0

# Log format
log_line_prefix = '%t [%p]: [%l-1] user=%u,db=%d,app=%a,client=%h '

# Log destination
log_destination = 'stderr'

# Log directory
logging_collector = on
log_directory = 'pg_log'
log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log'
log_file_mode = 0600
log_rotation_age = 1d
log_rotation_size = 100MB
```

### Log Aggregation

Log aggregation is essential for centralized log management and analysis. Here are some popular options:

#### ELK Stack

As described in the [Monitoring Tools](#monitoring-tools) section, the ELK Stack can be used for log aggregation and analysis.

#### Fluentd

Fluentd is an open-source data collector that can be used for log aggregation:

1. Install Fluentd:

```bash
# For Ubuntu/Debian
curl -L https://toolbelt.treasuredata.com/sh/install-ubuntu-bionic-td-agent3.sh | sh

# For CentOS/RHEL
curl -L https://toolbelt.treasuredata.com/sh/install-redhat-td-agent3.sh | sh
```

2. Configure Fluentd to collect CMDB Lite logs:

```xml
<!-- /etc/td-agent/td-agent.conf -->
<source>
  @type tail
  path /var/log/cmdb-lite/backend.log
  pos_file /var/log/td-agent/cmdb-lite-backend.log.pos
  tag cmdb-lite.backend
  format json
  time_format %Y-%m-%dT%H:%M:%S.%L%z
</source>

<source>
  @type tail
  path /var/log/cmdb-lite/frontend.log
  pos_file /var/log/td-agent/cmdb-lite-frontend.log.pos
  tag cmdb-lite.frontend
  format json
  time_format %Y-%m-%dT%H:%M:%S.%L%z
</source>

<match cmdb-lite.**>
  @type elasticsearch
  host localhost
  port 9200
  index_name cmdb-lite
  type_name _doc
</match>
```

3. Start Fluentd:

```bash
sudo systemctl start td-agent
sudo systemctl enable td-agent
```

#### Datadog Logs

Datadog provides log management as part of its monitoring platform:

1. Configure the Datadog Agent to collect logs:

```yaml
# /etc/datadog-agent/conf.d/cmdb-lite.d/conf.yaml
logs:
  - type: file
    path: /var/log/cmdb-lite/backend.log
    service: cmdb-lite-backend
    source: go
    log_processing_rules:
      - type: multi_line
        name: start_with_date
        pattern: \d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}

  - type: file
    path: /var/log/cmdb-lite/frontend.log
    service: cmdb-lite-frontend
    source: javascript
    log_processing_rules:
      - type: multi_line
        name: start_with_date
        pattern: \d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}
```

2. Restart the Datadog Agent:

```bash
sudo systemctl restart datadog-agent
```

## Troubleshooting

This section provides guidance on troubleshooting common issues with CMDB Lite.

### Common Issues

#### Application Fails to Start

**Symptoms**: The application fails to start with error messages

**Possible Causes**:
- Missing or invalid configuration
- Database connection issues
- Port already in use
- File permission issues

**Troubleshooting Steps**:
1. Check the application logs for error messages
2. Verify that all required environment variables are set correctly
3. Check that the database is accessible and the schema is up to date
4. Verify that the specified port is available
5. Check file permissions for the application directory

#### High Response Times

**Symptoms**: The application responds slowly to requests

**Possible Causes**:
- Database performance issues
- High resource usage (CPU, memory)
- Network latency
- Inefficient queries

**Troubleshooting Steps**:
1. Check resource usage (CPU, memory, disk I/O)
2. Analyze database performance (slow queries, connection pool usage)
3. Check network latency between application and database
4. Review application logs for errors or warnings
5. Profile the application to identify performance bottlenecks

#### Database Connection Errors

**Symptoms**: The application reports database connection errors

**Possible Causes**:
- Database is down or not accessible
- Database connection pool is exhausted
- Network issues between application and database
- Database authentication issues

**Troubleshooting Steps**:
1. Check that the database is running and accessible
2. Verify database connection parameters (host, port, username, password)
3. Check network connectivity between application and database
4. Review database connection pool settings
5. Check database logs for connection-related errors

#### Memory Issues

**Symptoms**: The application uses excessive memory or crashes due to out-of-memory errors

**Possible Causes**:
- Memory leaks
- Insufficient memory allocated
- Memory-intensive operations
- Large datasets being loaded into memory

**Troubleshooting Steps**:
1. Monitor memory usage over time to identify leaks
2. Increase memory allocation if necessary
3. Review code for memory-intensive operations
4. Optimize queries to load less data into memory
5. Implement pagination for large datasets

### Diagnostic Tools

#### Application Logs

Application logs are the primary source of information for troubleshooting:

```bash
# For Docker Compose
docker-compose logs backend

# For Kubernetes
kubectl logs <backend-pod-name> -n cmdb-lite

# For Traditional VM
sudo journalctl -u cmdb-lite-backend -f
```

#### Database Metrics

Use PostgreSQL's built-in monitoring views to diagnose database issues:

```sql
-- Check active connections
SELECT count(*) FROM pg_stat_activity;

-- Check long-running queries
SELECT pid, now() - pg_stat_activity.query_start AS duration, query, state
FROM pg_stat_activity
WHERE (now() - pg_stat_activity.query_start) > interval '5 minutes';

-- Check table sizes
SELECT schemaname, tablename, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- Check index usage
SELECT schemaname, tablename, indexname, idx_tup_read, idx_tup_fetch
FROM pg_stat_user_indexes
ORDER BY idx_tup_read + idx_tup_fetch DESC;
```

#### System Metrics

Use system tools to monitor resource usage:

```bash
# Check CPU usage
top
htop

# Check memory usage
free -h
vmstat

# Check disk usage
df -h
du -sh

# Check network connections
netstat -tuln
ss -tuln
```

#### Profiling Tools

Use profiling tools to identify performance bottlenecks:

```bash
# For Go backend
go tool pprof http://localhost:8080/debug/pprof/profile
go tool pprof http://localhost:8080/debug/pprof/heap
go tool pprof http://localhost:8080/debug/pprof/goroutine
```

### Troubleshooting Procedures

#### Application Performance Issues

1. **Identify the Problem**:
   - Check monitoring dashboards for anomalies
   - Review application logs for errors or warnings
   - Check system resource usage

2. **Isolate the Component**:
   - Determine if the issue is with the frontend, backend, or database
   - Check if the issue affects all users or a subset
   - Check if the issue is related to specific operations

3. **Analyze the Data**:
   - For frontend issues, check browser developer tools
   - For backend issues, check application logs and metrics
   - For database issues, check database logs and metrics

4. **Implement a Solution**:
   - For frontend issues, optimize JavaScript, CSS, or assets
   - For backend issues, optimize code, queries, or configuration
   - For database issues, optimize queries, indexes, or configuration

5. **Verify the Solution**:
   - Monitor the application after implementing the solution
   - Check that the issue is resolved
   - Check for any unintended side effects

#### Database Performance Issues

1. **Identify Slow Queries**:
   - Use PostgreSQL's pg_stat_statements to identify slow queries
   - Check application logs for slow query warnings
   - Review database metrics for query performance

2. **Analyze Query Execution Plans**:
   - Use EXPLAIN ANALYZE to understand how queries are executed
   - Look for full table scans, sequential scans, or inefficient joins
   - Identify missing indexes

3. **Optimize Queries**:
   - Rewrite queries to be more efficient
   - Add appropriate indexes
   - Consider denormalization for frequently accessed data

4. **Optimize Database Configuration**:
   - Adjust memory settings (shared_buffers, work_mem)
   - Adjust connection pool settings
   - Consider partitioning large tables

5. **Monitor Performance**:
   - Monitor query performance after optimization
   - Check that the issue is resolved
   - Check for any unintended side effects

## Performance Tuning

Performance tuning is essential for ensuring that CMDB Lite can handle the expected load and provide a good user experience.

### Application Performance

#### Backend Performance

1. **Database Connection Pooling**:
   - Optimize database connection pool settings
   - Monitor connection pool usage
   - Adjust pool size based on load

   Example configuration:

   ```go
   // Database configuration
   db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
   if err != nil {
       log.Fatal("Failed to connect to database:", err)
   }
   
   sqlDB, err := db.DB()
   if err != nil {
       log.Fatal("Failed to get database instance:", err)
   }
   
   // Configure connection pool
   sqlDB.SetMaxIdleConns(5)           // Maximum number of idle connections
   sqlDB.SetMaxOpenConns(25)          // Maximum number of open connections
   sqlDB.SetConnMaxLifetime(time.Hour) // Maximum time a connection can be reused
   ```

2. **Caching**:
   - Implement caching for frequently accessed data
   - Use Redis or Memcached for distributed caching
   - Set appropriate cache expiration times

   Example implementation:

   ```go
   import (
       "github.com/go-redis/redis/v8"
   )
   
   var redisClient *redis.Client
   
   func init() {
       redisClient = redis.NewClient(&redis.Options{
           Addr:     "localhost:6379",
           Password: "", // no password set
           DB:       0,  // use default DB
       })
   }
   
   func GetCI(id string) (*models.ConfigurationItem, error) {
       // Try to get from cache first
       cached, err := redisClient.Get(ctx, "ci:"+id).Result()
       if err == nil {
           var ci models.ConfigurationItem
           if err := json.Unmarshal([]byte(cached), &ci); err == nil {
               return &ci, nil
           }
       }
       
       // If not in cache, get from database
       var ci models.ConfigurationItem
       if err := db.First(&ci, "id = ?", id).Error; err != nil {
           return nil, err
       }
       
       // Cache the result
       if data, err := json.Marshal(ci); err == nil {
           redisClient.Set(ctx, "ci:"+id, data, 5*time.Minute)
       }
       
       return &ci, nil
   }
   ```

3. **Concurrency**:
   - Optimize concurrency settings for the Go backend
   - Use goroutines and channels effectively
   - Avoid blocking operations

   Example implementation:

   ```go
   func ProcessCIs(cis []models.ConfigurationItem) []models.ProcessedCI {
       var wg sync.WaitGroup
       results := make([]models.ProcessedCI, len(cis))
       
       for i, ci := range cis {
           wg.Add(1)
           go func(idx int, c models.ConfigurationItem) {
               defer wg.Done()
               results[idx] = ProcessCI(c)
           }(i, ci)
       }
       
       wg.Wait()
       return results
   }
   ```

4. **Memory Management**:
   - Monitor memory usage
   - Identify and fix memory leaks
   - Optimize memory allocation

   Example implementation:

   ```go
   func GetLargeDataset() ([]models.ConfigurationItem, error) {
       var cis []models.ConfigurationItem
       
       // Use pagination to avoid loading all data at once
       offset := 0
       limit := 100
       
       for {
           var batch []models.ConfigurationItem
           if err := db.Offset(offset).Limit(limit).Find(&batch).Error; err != nil {
               return nil, err
           }
           
           if len(batch) == 0 {
               break
           }
           
           cis = append(cis, batch...)
           offset += limit
       }
       
       return cis, nil
   }
   ```

#### Frontend Performance

1. **Code Splitting**:
   - Implement code splitting to reduce initial bundle size
   - Use dynamic imports for large components
   - Lazy load routes

   Example implementation:

   ```javascript
   // Router configuration with lazy loading
   const routes = [
     {
       path: '/dashboard',
       name: 'Dashboard',
       component: () => import('@/views/DashboardView.vue')
     },
     {
       path: '/cis',
       name: 'CIList',
       component: () => import('@/views/CIListView.vue')
     }
   ]
   
   // Dynamic import for large components
   const LargeComponent = defineAsyncComponent(() =>
     import('@/components/LargeComponent.vue')
   )
   ```

2. **Caching**:
   - Implement browser caching for static assets
   - Use service workers for offline functionality
   - Cache API responses where appropriate

   Example implementation:

   ```javascript
   // Service worker for caching
   const CACHE_NAME = 'cmdb-lite-v1';
   const urlsToCache = [
     '/',
     '/static/js/main.js',
     '/static/css/main.css',
     '/api/v1/cis'
   ];
   
   self.addEventListener('install', event => {
     event.waitUntil(
       caches.open(CACHE_NAME)
         .then(cache => cache.addAll(urlsToCache))
     );
   });
   
   self.addEventListener('fetch', event => {
     event.respondWith(
       caches.match(event.request)
         .then(response => {
           if (response) {
             return response;
           }
           return fetch(event.request);
         })
     );
   });
   ```

3. **Optimization**:
   - Optimize images and other assets
   - Minimize CSS and JavaScript
   - Use compression for responses

   Example implementation:

   ```javascript
   // Image optimization
   function optimizeImage(imageUrl, maxWidth, maxHeight) {
     return new Promise((resolve, reject) => {
       const img = new Image();
       img.crossOrigin = 'Anonymous';
       img.src = imageUrl;
       
       img.onload = () => {
         const canvas = document.createElement('canvas');
         const ctx = canvas.getContext('2d');
         
         let width = img.width;
         let height = img.height;
         
         if (width > maxWidth) {
           height *= maxWidth / width;
           width = maxWidth;
         }
         
         if (height > maxHeight) {
           width *= maxHeight / height;
           height = maxHeight;
         }
         
         canvas.width = width;
         canvas.height = height;
         
         ctx.drawImage(img, 0, 0, width, height);
         
         canvas.toBlob(resolve, 'image/jpeg', 0.8);
       };
       
       img.onerror = reject;
     });
   }
   ```

### Database Performance

1. **Indexing**:
   - Ensure proper indexing for frequently queried columns
   - Use composite indexes for queries with multiple conditions
   - Regularly analyze index usage

   Example implementation:

   ```sql
   -- Create indexes for frequently queried columns
   CREATE INDEX idx_configuration_items_name ON configuration_items(name);
   CREATE INDEX idx_configuration_items_type ON configuration_items(type);
   CREATE INDEX idx_configuration_items_created_at ON configuration_items(created_at);
   
   -- Create composite index for queries with multiple conditions
   CREATE INDEX idx_configuration_items_type_created_at ON configuration_items(type, created_at);
   
   -- Analyze index usage
   SELECT schemaname, tablename, indexname, idx_tup_read, idx_tup_fetch
   FROM pg_stat_user_indexes
   ORDER BY idx_tup_read + idx_tup_fetch DESC;
   ```

2. **Query Optimization**:
   - Optimize slow queries
   - Use EXPLAIN ANALYZE to understand query execution plans
   - Avoid SELECT * and only select necessary columns

   Example implementation:

   ```sql
   -- Analyze query execution plan
   EXPLAIN ANALYZE SELECT id, name, type FROM configuration_items WHERE type = 'Server' AND created_at > '2023-01-01';
   
   -- Optimize query by adding appropriate index
   CREATE INDEX idx_configuration_items_type_created_at ON configuration_items(type, created_at);
   
   -- Only select necessary columns
   SELECT id, name, type FROM configuration_items WHERE type = 'Server' AND created_at > '2023-01-01';
   ```

3. **Database Configuration**:
   - Optimize PostgreSQL configuration for your workload
   - Adjust memory settings, connection settings, and checkpoint settings

   Example configuration:

   ```ini
   # Memory settings
   shared_buffers = 256MB
   effective_cache_size = 1GB
   work_mem = 4MB
   maintenance_work_mem = 64MB
   
   # Connection settings
   max_connections = 100
   
   # Query optimization
   random_page_cost = 1.1
   effective_io_concurrency = 200
   
   # Checkpoint settings
   checkpoint_completion_target = 0.9
   wal_buffers = 16MB
   
   # Logging
   log_min_duration_statement = 1000
   log_checkpoints = on
   log_connections = on
   log_disconnections = on
   ```

4. **Partitioning**:
   - Consider table partitioning for large tables
   - Use range or list partitioning based on access patterns

   Example implementation:

   ```sql
   -- Create a partitioned table
   CREATE TABLE configuration_items (
       id UUID PRIMARY KEY,
       name VARCHAR(255) NOT NULL,
       type VARCHAR(100) NOT NULL,
       attributes JSONB,
       tags TEXT[],
       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
       updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
   ) PARTITION BY RANGE (created_at);
   
   -- Create partitions
   CREATE TABLE configuration_items_2023_q1 PARTITION OF configuration_items
       FOR VALUES FROM ('2023-01-01') TO ('2023-04-01');
   
   CREATE TABLE configuration_items_2023_q2 PARTITION OF configuration_items
       FOR VALUES FROM ('2023-04-01') TO ('2023-07-01');
   
   CREATE TABLE configuration_items_2023_q3 PARTITION OF configuration_items
       FOR VALUES FROM ('2023-07-01') TO ('2023-10-01');
   
   CREATE TABLE configuration_items_2023_q4 PARTITION OF configuration_items
       FOR VALUES FROM ('2023-10-01') TO ('2024-01-01');
   ```

### Infrastructure Performance

1. **Resource Allocation**:
   - Ensure adequate CPU, memory, and disk resources
   - Monitor resource usage and adjust as needed
   - Consider vertical scaling for resource-intensive components

2. **Storage Performance**:
   - Use SSD storage for better I/O performance
   - Monitor disk I/O and latency
   - Consider using separate storage for database and logs

3. **Network Performance**:
   - Ensure adequate network bandwidth and low latency
   - Monitor network traffic and latency
   - Consider using a content delivery network (CDN) for static assets

4. **Load Balancing**:
   - Use load balancing to distribute traffic across multiple instances
   - Implement health checks for load balancer
   - Consider using session affinity if needed

## Health Checks

Health checks are essential for monitoring the health of CMDB Lite and its components.

### Application Health Checks

#### Backend Health Check

The backend provides a health check endpoint at `/health` that returns the overall health of the application and its components:

```go
func HealthCheckHandler(c *gin.Context) {
    // Check database connection
    sqlDB, err := database.DB.DB()
    if err != nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": "Database connection error"})
        return
    }
    
    err = sqlDB.Ping()
    if err != nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": "Database ping error"})
        return
    }
    
    // Check other dependencies
    // ...
    
    c.JSON(http.StatusOK, gin.H{
        "status": "healthy",
        "timestamp": time.Now().UTC().Format(time.RFC3339),
        "version": "1.0.0",
        "components": map[string]interface{}{
            "database": map[string]interface{}{
                "status": "healthy",
                "details": map[string]interface{}{
                    "connection_time_ms": 5,
                    "version":            "13.4",
                },
            },
        },
    })
}
```

#### Frontend Health Check

The frontend can implement a health check that verifies the backend API is accessible:

```javascript
// Health check service
const healthCheckService = {
  async checkHealth() {
    try {
      const response = await fetch('/api/v1/health');
      if (response.ok) {
        const data = await response.json();
        return {
          status: 'healthy',
          backend: data.status,
          timestamp: new Date().toISOString(),
        };
      } else {
        return {
          status: 'unhealthy',
          backend: 'unreachable',
          timestamp: new Date().toISOString(),
        };
      }
    } catch (error) {
      return {
        status: 'unhealthy',
        backend: 'error',
        error: error.message,
        timestamp: new Date().toISOString(),
      };
    }
  },
};

// Periodic health check
setInterval(async () => {
  const health = await healthCheckService.checkHealth();
  console.log('Health check:', health);
  
  // Send to monitoring service
  if (health.status !== 'healthy') {
    sendAlert('Application health check failed', health);
  }
}, 60000); // Check every minute
```

### Database Health Checks

#### PostgreSQL Health Check

PostgreSQL provides several ways to check the health of the database:

1. **Connection Check**:
   ```sql
   SELECT 1;
   ```

2. **Database Statistics**:
   ```sql
   SELECT * FROM pg_stat_database;
   ```

3. **Table Statistics**:
   ```sql
   SELECT * FROM pg_stat_user_tables;
   ```

4. **Index Statistics**:
   ```sql
   SELECT * FROM pg_stat_user_indexes;
   ```

Example health check script:

```bash
#!/bin/bash

# Configuration
DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="cmdb_lite"
DB_USER="cmdb_user"
DB_PASSWORD="secure_password"

# Check database connection
if ! PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT 1;" > /dev/null 2>&1; then
    echo "Database connection failed"
    exit 1
fi

# Check database size
DB_SIZE=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT pg_size_pretty(pg_database_size('$DB_NAME'));")
echo "Database size: $DB_SIZE"

# Check number of connections
CONNECTIONS=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT count(*) FROM pg_stat_activity;")
echo "Active connections: $CONNECTIONS"

# Check for long-running queries
LONG_RUNNING_QUERIES=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT count(*) FROM pg_stat_activity WHERE (now() - pg_stat_activity.query_start) > interval '5 minutes';")
echo "Long-running queries: $LONG_RUNNING_QUERIES"

echo "Database health check completed successfully"
```

### Infrastructure Health Checks

#### System Health Check

System health checks can be implemented using various tools:

1. **System Resource Check**:
   ```bash
   # Check CPU usage
   CPU_USAGE=$(top -bn1 | grep "Cpu(s)" | sed "s/.*, *\([0-9.]*\)%* id.*/\1/" | awk '{print 100 - $1}')
   echo "CPU usage: $CPU_USAGE%"
   
   # Check memory usage
   MEMORY_USAGE=$(free | grep Mem | awk '{print ($3/$2) * 100.0}')
   echo "Memory usage: $MEMORY_USAGE%"
   
   # Check disk usage
   DISK_USAGE=$(df -h / | awk 'NR==2 {print $5}' | sed 's/%//')
   echo "Disk usage: $DISK_USAGE%"
   
   # Check if any resource usage is above threshold
   if (( $(echo "$CPU_USAGE > 80" | bc -l) )); then
       echo "WARNING: High CPU usage"
   fi
   
   if (( $(echo "$MEMORY_USAGE > 80" | bc -l) )); then
       echo "WARNING: High memory usage"
   fi
   
   if [ $DISK_USAGE -gt 80 ]; then
       echo "WARNING: High disk usage"
   fi
   ```

2. **Network Health Check**:
   ```bash
   # Check network latency
   LATENCY=$(ping -c 1 google.com | grep -oP 'time=\K\d+\.?\d*')
   echo "Network latency: ${LATENCY}ms"
   
   # Check if latency is above threshold
   if (( $(echo "$LATENCY > 100" | bc -l) )); then
       echo "WARNING: High network latency"
   fi
   ```

#### Service Health Check

Service health checks can be implemented using systemd or other service managers:

1. **Systemd Service Check**:
   ```bash
   # Check if a service is running
   if systemctl is-active --quiet cmdb-lite-backend; then
       echo "cmdb-lite-backend is running"
   else
       echo "WARNING: cmdb-lite-backend is not running"
       systemctl start cmdb-lite-backend
   fi
   
   if systemctl is-active --quiet cmdb-lite-frontend; then
       echo "cmdb-lite-frontend is running"
   else
       echo "WARNING: cmdb-lite-frontend is not running"
       systemctl start cmdb-lite-frontend
   fi
   
   if systemctl is-active --quiet postgresql; then
       echo "postgresql is running"
   else
       echo "WARNING: postgresql is not running"
       systemctl start postgresql
   fi
   ```

2. **Process Check**:
   ```bash
   # Check if a process is running
   if pgrep -f "cmdb-lite-backend" > /dev/null; then
       echo "cmdb-lite-backend process is running"
   else
       echo "WARNING: cmdb-lite-backend process is not running"
   fi
   
   if pgrep -f "cmdb-lite-frontend" > /dev/null; then
       echo "cmdb-lite-frontend process is running"
   else
       echo "WARNING: cmdb-lite-frontend process is not running"
   fi
   ```

For more information on operating CMDB Lite, see the [Operations Guide](README.md) and the [Deployment Guide](deployment.md).
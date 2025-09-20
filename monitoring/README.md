# Monitoring and Observability Setup

This directory contains the configuration files and documentation for setting up monitoring and observability for the CMDB Lite application.

## Components

### 1. Prometheus
Prometheus is used for collecting and storing metrics from the application.

#### Configuration
- `prometheus/prometheus.yml` - Prometheus configuration file
- `prometheus/alert_rules.yml` - Alert rules for Prometheus

#### Key Metrics
- HTTP request metrics (count, duration, status codes)
- Database query metrics (count, duration, active connections)
- Business metrics (CI count, relationship count, user count, etc.)
- System metrics (CPU, memory, disk usage)

### 2. Grafana
Grafana is used for visualizing metrics and creating dashboards.

#### Configuration
- `grafana/provisioning/dashboards/dashboard.yml` - Dashboard provisioning configuration
- `grafana/provisioning/datasources/datasource.yml` - Data source configuration

#### Dashboards
- `grafana/dashboards/cmdb-application-metrics.json` - Application metrics dashboard
- `grafana/dashboards/cmdb-business-metrics.json` - Business metrics dashboard

### 3. Alertmanager
Alertmanager is used for handling alerts sent by Prometheus.

#### Configuration
- `alertmanager/config.yml` - Alertmanager configuration

#### Notification Channels
- Email notifications
- Webhook notifications

### 4. Jaeger
Jaeger is used for distributed tracing.

#### Configuration
- `jaeger/docker-compose.yml` - Jaeger Docker Compose configuration

#### Key Features
- Distributed tracing
- Request correlation
- Performance analysis

### 5. Node Exporter
Node Exporter is used for collecting system metrics.

### 6. cAdvisor
cAdvisor is used for collecting container metrics.

## Setup Instructions

### 1. Start the monitoring stack

```bash
# Navigate to the monitoring directory
cd monitoring

# Start the monitoring stack
docker-compose up -d
```

### 2. Start Jaeger (for distributed tracing)

```bash
# Navigate to the Jaeger directory
cd monitoring/jaeger

# Start Jaeger
docker-compose up -d
```

### 3. Access the monitoring tools

- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (username: admin, password: admin)
- **Alertmanager**: http://localhost:9093
- **Jaeger**: http://localhost:16686
- **Node Exporter**: http://localhost:9100
- **cAdvisor**: http://localhost:8080

## Configuration

### Environment Variables

The following environment variables can be used to configure the monitoring:

- `TRACING_ENABLED` - Enable/disable tracing (default: false)
- `TRACING_SERVICE` - Service name for tracing (default: cmdb-lite)
- `TRACING_ENV` - Environment for tracing (default: development)
- `TRACING_JAEGER_URL` - Jaeger collector URL (default: "")
- `TRACING_ZIPKIN_URL` - Zipkin collector URL (default: "")
- `TRACING_SAMPLING_RATE` - Sampling rate for tracing (default: 1.0)
- `METRICS_ENABLED` - Enable/disable metrics (default: true)
- `METRICS_PATH` - Path for metrics endpoint (default: /metrics)

### Application Configuration

To enable monitoring in the application, set the following environment variables:

```bash
# Enable tracing
export TRACING_ENABLED=true
export TRACING_JAEGER_URL=http://jaeger:14268/api/traces

# Enable metrics
export METRICS_ENABLED=true
```

## Alerting

### Alert Rules

The following alert rules are configured:

1. **System Alerts**
   - High CPU usage (>80%)
   - Very high CPU usage (>90%)
   - High memory usage (>80%)
   - Very high memory usage (>90%)
   - High disk usage (>80%)
   - Very high disk usage (>90%)

2. **Application Alerts**
   - High error rate (>10%)
   - Very high error rate (>20%)
   - High latency (>1s)
   - Very high latency (>2s)
   - Service down

3. **Database Alerts**
   - High database query latency (>500ms)
   - Very high database query latency (>1s)
   - High database connection usage (>80%)
   - Very high database connection usage (>90%)

4. **Business Alerts**
   - High authentication failures (>10 per minute)
   - Very high authentication failures (>20 per minute)

### Notification Channels

Alerts can be sent to the following channels:

1. **Email**
   - Configure SMTP settings in Alertmanager configuration

2. **Webhook**
   - Configure webhook URL in Alertmanager configuration

## Customization

### Adding New Metrics

To add new metrics to the application:

1. Define the metric in the `internal/metrics/metrics.go` file
2. Update the metric value in the appropriate handler
3. Add the metric to the Grafana dashboard

### Adding New Alerts

To add new alerts:

1. Define the alert rule in the `prometheus/alert_rules.yml` file
2. Configure the notification channel in the `alertmanager/config.yml` file

### Adding New Dashboards

To add new dashboards:

1. Create a new dashboard JSON file in the `grafana/dashboards/` directory
2. Update the `grafana/provisioning/dashboards/dashboard.yml` file to include the new dashboard

## Troubleshooting

### Common Issues

1. **Prometheus not scraping metrics**
   - Check if the application is running
   - Check if the metrics endpoint is accessible
   - Check the Prometheus configuration file

2. **Grafana not showing data**
   - Check if the data source is configured correctly
   - Check if Prometheus is collecting metrics
   - Check the dashboard queries

3. **Alerts not firing**
   - Check if the alert rules are configured correctly
   - Check if the Alertmanager is running
   - Check the notification channel configuration

4. **Tracing not working**
   - Check if tracing is enabled in the application
   - Check if the Jaeger collector URL is correct
   - Check if the sampling rate is appropriate

### Logs

Check the logs of the monitoring components for troubleshooting:

```bash
# Check Prometheus logs
docker logs prometheus

# Check Grafana logs
docker logs grafana

# Check Alertmanager logs
docker logs alertmanager

# Check Jaeger logs
docker logs jaeger
```

## Scaling

For production environments, consider the following:

1. **Persistent Storage**
   - Configure persistent storage for Prometheus and Grafana data

2. **High Availability**
   - Set up multiple instances of Prometheus and Grafana
   - Use a load balancer for distributing traffic

3. **Retention**
   - Configure appropriate retention periods for metrics and logs

4. **Security**
   - Enable authentication and authorization for monitoring tools
   - Use HTTPS for all monitoring endpoints
   - Restrict access to monitoring tools

## Integration with Other Tools

The monitoring stack can be integrated with other tools:

1. **Log Management**
   - Integrate with ELK stack (Elasticsearch, Logstash, Kibana)
   - Integrate with Loki and Grafana

2. **Incident Management**
   - Integrate with PagerDuty, Opsgenie, or similar tools

3. **ChatOps**
   - Integrate with Slack, Microsoft Teams, or similar tools

4. **CI/CD**
   - Integrate with Jenkins, GitLab CI, or similar tools
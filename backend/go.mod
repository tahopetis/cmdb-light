module github.com/cmdb-lite/backend

go 1.21

require (
	github.com/google/uuid v1.3.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/lib/pq v1.10.9
	github.com/pressly/goose/v3 v3.12.0
	github.com/gorilla/mux v1.8.0
	github.com/golang-jwt/jwt/v5 v5.0.0
	golang.org/x/crypto v0.11.0
	github.com/stretchr/testify v1.8.4
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/testcontainers/testcontainers-go v0.20.1
	github.com/gavv/httpexpect/v2 v2.13.0
	github.com/ory/dockertest v3.3.5+incompatible
	github.com/pact-foundation/pact-go v2.0.0+incompatible
	
	// Monitoring and logging dependencies
	github.com/rs/zerolog v1.30.0
	github.com/prometheus/client_golang v1.16.0
	github.com/prometheus/client_model v0.4.0
	github.com/prometheus/common v0.42.0
	github.com/prometheus/procfs v0.10.1
	github.com/gorilla/handlers v1.5.1
	
	// Distributed tracing dependencies
	go.opentelemetry.io/otel v1.19.0
	go.opentelemetry.io/otel/exporters/jaeger v1.17.0
	go.opentelemetry.io/otel/exporters/zipkin v1.17.0
	go.opentelemetry.io/otel/sdk v1.19.0
	go.opentelemetry.io/otel/trace v1.19.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.44.0
	go.opentelemetry.io/contrib/instrumentation/github.com/lib/pq/otelpq v0.44.0
)
🛠️ Technical Specification Document (TSD)

Project: CMDB Lite
Version: 1.2
Date: 2025-09-20
Author: [Your Name]

1. 🎯 Purpose

CMDB Lite is a lightweight, open-source Configuration Management Database (CMDB) designed to manage configuration items (CIs), their attributes, and relationships.
This document describes the technical design of the system, covering architecture, data model, API contracts, security, testing, and deployment.

2. 🏗️ System Architecture
2.1 High-Level Architecture (Mermaid)
flowchart LR
    subgraph Client
        A[Vue 3 Frontend] -->|REST Calls| B
    end

    subgraph Backend
        B[Go REST API (chi + sqlx)]
    end

    subgraph Database
        C[(PostgreSQL)]
    end

    B -->|SQL Queries| C

    subgraph Deployment
        D[Docker Compose / Kubernetes]
    end

    A -.->|Served as SPA| D
    B -.->|Containerized| D
    C -.->|Stateful Service| D


Frontend: Vue 3 SPA, TailwindCSS, D3.js for visualization.

Backend: Go with chi router, sqlx for DB, JWT authentication.

Database: PostgreSQL (schema migration with Goose or migrate).

Deployment: Docker Compose for local dev, optional Helm chart for Kubernetes production.

3. 📂 Data Model
3.1 ERD (Mermaid)
erDiagram
    CONFIGURATION_ITEMS {
        uuid id PK
        string name
        string type
        jsonb attributes
        text[] tags
        timestamp created_at
        timestamp updated_at
    }

    RELATIONSHIPS {
        uuid id PK
        uuid source_id FK
        uuid target_id FK
        string type
        timestamp created_at
    }

    AUDIT_LOGS {
        uuid id PK
        string entity_type
        uuid entity_id
        string action
        string changed_by
        timestamp changed_at
        jsonb details
    }

    USERS {
        uuid id PK
        string username
        string password_hash
        string role
        timestamp created_at
    }

    CONFIGURATION_ITEMS ||--o{ RELATIONSHIPS : "has"
    CONFIGURATION_ITEMS ||--o{ AUDIT_LOGS : "changes recorded"
    USERS ||--o{ AUDIT_LOGS : "made changes"

4. 📡 API Specification
Endpoint	Method	Auth	Description
/auth/login	POST	❌	User login, returns JWT
/cis	GET	✅	List all configuration items (paginated)
/cis	POST	✅	Create new configuration item
/cis/{id}	GET	✅	Get CI details
/cis/{id}	PUT	✅	Update CI
/cis/{id}	DELETE	✅	Delete CI
/cis/{id}/graph	GET	✅	Retrieve related CIs (graph traversal)
/relationships	POST	✅	Create relationship
/relationships/{id}	DELETE	✅	Delete relationship
/audit-logs	GET	✅	Get audit log (filter by CI/user)
5. 🔄 Sequence Diagrams
5.1 Login Flow (Mermaid)
sequenceDiagram
    participant U as User
    participant F as Frontend (Vue)
    participant A as API (Go)
    participant D as Database (Postgres)

    U->>F: Enter username/password
    F->>A: POST /auth/login
    A->>D: SELECT * FROM users WHERE username=?
    D-->>A: User record (with password hash)
    A->>A: Validate password & generate JWT
    A-->>F: 200 OK + JWT token
    F->>F: Store token (localStorage)

5.2 CI Creation Flow (Mermaid)
sequenceDiagram
    participant U as User
    participant F as Frontend
    participant A as API
    participant D as Database

    U->>F: Fill CI Form + Submit
    F->>A: POST /cis (JWT)
    A->>D: INSERT INTO configuration_items(...)
    D-->>A: Success (CI ID)
    A->>D: INSERT INTO audit_logs(...)
    A-->>F: 201 Created + CI Object
    F->>F: Update UI with new CI

6. 🔒 Security Design

Auth: JWT Bearer tokens (exp 24h)

Password: Argon2id / bcrypt hashing

RBAC: Roles: admin, viewer

Transport: HTTPS enforced in production

Audit: Immutable logs for create/update/delete actions

7. 🧪 Testing Strategy

Unit: Go testing for handlers & repos

Integration: Testcontainers for Postgres

Frontend: Vitest + Vue Testing Library

E2E: Playwright or Cypress

8. 🛠️ Deployment

Dev: Docker Compose spins up API + DB + Adminer + Vue Dev Server

Prod: Multi-stage Docker builds, optionally Helm for K8s

Config: ENV vars (DB_URL, JWT_SECRET)

Backup: Nightly pg_dump (retention policy 7 days)

9. 📊 Monitoring & Logging

Logs: Zerolog JSON structured logs with request ID

Metrics: Prometheus /metrics endpoint

Dashboards: Grafana for latency, error rate
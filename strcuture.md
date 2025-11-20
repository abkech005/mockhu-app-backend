# Project Structure Documentation (Go + Postgres + Redis)

This document describes a clean, production‑ready folder structure for a scalable backend built using **Go**, **PostgreSQL**, and **Redis**. It follows principles from Clean Architecture and Hexagonal Architecture.

---

## 📁 High-Level Structure

```
/cmd/
  api/
  worker/
/internal/
  app/
    user/
    post/
    follow/
    feed/
    notification/
  ports/
    repo/
    cache/
    queue/
  transport/
    http/
      handlers/
      middleware/
      dto/
  infra/
    db/
    cache/
    queue/
    storage/
    email/
  util/
/pkg/
  auth/
  middleware/
  errors/
/migrations/
/configs/
/scripts/
/deploy/
/tests/
```

---

## 📦 Folders Explained

### `/cmd/`

Entry points for the application.

* `api/` — HTTP API server.
* `worker/` — background jobs, cron processors.

### `/internal/`

The main application code. Not exposed outside the module.

#### `/internal/app/` — Domain & Use‑Cases

Each feature/module lives here.
Examples:

* `user/`
* `post/`
* `follow/`
* `feed/`

Each module usually contains:

* `domain.go` — core entity definitions.
* `service.go` — business logic.
* `repo.go` — repository interfaces for the module.
* `transport.go` — DTOs and request validation.

#### `/internal/ports/` — Abstract Interfaces

Defines the interfaces that the domain depends on.

* `repo/` — DB repositories.
* `cache/` — Redis interfaces.
* `queue/` — message queue interfaces.

#### `/internal/transport/` — HTTP Layer

Includes routing, handlers and middleware.

* `handlers/`
* `middleware/`
* `dto/` — input/output schemas.

#### `/internal/infra/` — Concrete Implementations

Actual implementations of ports.

* `db/` — Postgres repos.
* `cache/` — Redis client helpers.
* `queue/` — Kafka/RabbitMQ/etc.
* `storage/` — S3, local files.
* `email/` — SMTP provider.

#### `/internal/util/`

Internal utilities: time, context, pagination.

---

## `/pkg/` — Public Reusable Components

Everything placed here is safe for other projects to import.
Examples:

* `auth/` — JWT utilities.
* `middleware/` — logging, tracing, CORS.
* `errors/` — custom error formatting.

---

## `/migrations/`

SQL migrations for Postgres.

```
0001_init.up.sql
0001_init.down.sql
```

---

## `/configs/`

Config files, environment templates.

* `.env.example`
* `app.yaml`

---

## `/scripts/`

Local development scripts.

* DB seeders
* Helpers for running tests

---

## `/deploy/`

Production/CI deployment files.

* Kubernetes manifests
* Dockerfiles
* Helm charts

---

## `/tests/`

Integration and end-to-end tests.
Use libraries like `testcontainers-go` for real Postgres + Redis tests.

---

## 🧩 Architecture Overview

### Flow of a Request:

1. **Transport Layer** receives request (HTTP).
2. Parses DTO, validates.
3. Calls application **service** from `/internal/app/<module>`.
4. Service runs domain logic.
5. Service calls repositories/cache/queues through **ports** interfaces.
6. Infra implementations execute DB/Redis operations.
7. Response returned back to transport layer.

This maintains:

* Module isolation
* Easy testing
* Scalability
* Change‑friendly structure

---

If you want, I can create a full README.md with installation steps, environment setup, development scripts, and example API flow.

# MiniKube
MiniKube is a simplified container orchestration system inspired by [**Kubernetes**](https://kubernetes.io/docs/home/). It is designed as a learning-focused distributed systems project that demonstrates how modern container orchestration works internally.

The project includes:
- A custom control plane
- Pod scheduling
- Controller reconciliation loops
- Node heartbeats
- Container execution using Docker
- Service discovery and networking
- CLI tooling similar to kubectl

---

## Contents
1. [Technological Overview](#technological-overview)
2. [Tech Stack](#tech-stack)
3. [Building Phases](#building-phases)
    - [Phase 1 - Go foundations + project skeleton](#phase-1---go-foundations--project-skeleton)
    - [Phase 2 - Control plane core](#phase-2---control-plane-core)

---

## Technological Overview
![alt text](image-1.png)

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go |
| API Server | net/http / chi |
| CLI | Cobra |
| State Store | BoltDB |
| Communication | gRPC |
| Runtime | Docker SDK |
| Dashboard | Go + React/NextJs |

---

## Building Phases

### Phase 1 - Go foundations + project skeleton
- **Goal**: Get comfortable with Go patterns you'll use everywhere before touching orchestration logic.
- **What we built**: A small CLI tool and a basic HTTP server, nothing MiniKube-specific yet, just Go muscle memory.
- **Learnings**: `cobra` (CLI framework), `net/http`, JSON marshalling, and Go project layout (`cmd/`, `internal/`, `pkg/`).
- **Deliverable**: A `minik` CLI binary that can `ping` a running server and get back a response.

### Phase 2 - Control plane core
- **Goal**: Build the brain of MiniKube — the API, state store, and scheduler.
- **What we built**: A structured REST API with chi router, a BoltDB embedded state store, pod status constants, and a background scheduler that assigns pending pods to nodes.
- **Learnings**: `chi` router and method-based routing, BoltDB buckets and transactions (`db.Update`, `db.View`), goroutines and `time.NewTicker` for background loops, Go struct methods, UUID generation, and proper separation of concerns across `internal/api`, `internal/store`, and `internal/scheduler`.
- **Deliverable**: `POST /pods` creates a pod persisted in BoltDB with status `PENDING`. Within 5 seconds the scheduler picks it up, assigns it to a node round-robin, and updates its status to `SCHEDULED`. `GET /pods` reflects the live state.
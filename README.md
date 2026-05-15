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
- **What we built**: a small CLI tool and a basic HTTP server, nothing MiniKube-specific yet, just Go muscle memory.
- **Learnings**: `cobra` (CLI framework), `net/http`, JSON marshalling, and Go project layout (`cmd/`, `internal/`, `pkg/`).
- **Deliverable**: a `minik` CLI binary that can `ping` a running server and get back a response.
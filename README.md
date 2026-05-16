# MiniKube
> A lightweight container orchestrator built in Go, inspired by Kubernetes.
 
MiniKube is a distributed systems project that implements core Kubernetes concepts from scratch — pod scheduling, controller reconciliation loops, node heartbeats, container execution via Docker, service discovery, and a full CLI and web dashboard. Built as a deep dive into how modern container orchestration actually works internally.

---

## Contents
1. [Features](#features)
2. [Architecture](#architecture)
3. [Tech Stack](#tech-stack)
4. [Requirements](#requirements)
5. [Installation](#installation)
6. [Quick Start](#quick-start)
7. [CLI Reference](#cli-reference)
8. [How It Works](#how-it-works)
9. [API Reference](#api-reference)
10. [Project Structure](#project-structure)
11. [Building Phases](#building-phases)
    - [Phase 1 - Go foundations + project skeleton](#phase-1---go-foundations--project-skeleton)
    - [Phase 2 - Control plane core](#phase-2---control-plane-core)
    - [Phase 3 - Worker node and Docker container lifecycle](#phase-3---worker-node-and-docker-container-lifecycle)
    - [Phase 4 - Service Discovery and Load Balancing](#phase-4---service-discovery-and-load-balancing)
    - [Phase 5 - Multi-node support and HTTP-based worker communication](#phase-5---multi-node-support-and-http-based-worker-communication)
    - [Phase 6 - CLI commands](#phase-6---cli-commands)
    - [Phase 7 - Dashboard UI and command](#phase-7---dashboard-ui-and-command)

---
 
## Features
 
- **Pod scheduling** — round-robin scheduler assigns pods to healthy worker nodes automatically
- **Reconciliation loop** — controller continuously diffs desired vs actual state and self-heals
- **Real containers** — Docker SDK under the hood, actual containers are started and stopped
- **Multi-node support** — run multiple worker nodes as separate processes, each handling their own pods
- **Node heartbeats** — workers ping the control plane every 5 seconds; missing heartbeats mark nodes as `NOT_READY`
- **Service discovery** — named services route traffic across pods with round-robin load balancing
- **Full CLI** — `minik` CLI with `get`, `delete`, `apply`, `cluster`, and `dashboard` commands
- **Web dashboard** — live Next.js dashboard showing cluster state, auto-refreshing every 5 seconds
- **Single command setup** — `minik cluster start` spins up the entire cluster in the background
---
 
## Architecture
 
```
┌─────────────────────────────────────────────┐
│                Control Plane                │
│                                             │
│   ┌──────────┐  ┌───────────┐  ┌────────┐   │
│   │API Server│  │ Scheduler │  │  Store │   │
│   │ (chi)    │  │(round-rbn)│  │(BoltDB)│   │
│   └──────────┘  └───────────┘  └────────┘   │
└─────────────────────────────────────────────┘
          │               │
    HTTP  │               │ HTTP
          ▼               ▼
┌──────────────┐   ┌──────────────┐
│   Worker 1   │   │   Worker 2   │
│              │   │              │
│ Docker SDK   │   │ Docker SDK   │
│ [container]  │   │ [container]  │
└──────────────┘   └──────────────┘
 
┌─────────────────────────────────────────────┐
│              minik CLI + Dashboard          │
│   minik get pods / apply / cluster start    │
│   localhost:3000 (Next.js dashboard)        │
└─────────────────────────────────────────────┘
```
 
The control plane runs as a single server process. Workers are separate processes that register with the control plane, receive pod assignments, and execute containers via the Docker SDK. The CLI and dashboard both communicate with the control plane over HTTP.
 
---
 
## Tech Stack
 
| Layer | Technology |
|---|---|
| Language | Go 1.24 |
| API Server | `net/http` + `chi` router |
| CLI | Cobra |
| State Store | BoltDB (embedded key-value) |
| Container Runtime | Docker SDK for Go |
| Dashboard Frontend | Next.js + Tailwind CSS |
| UUID Generation | `github.com/google/uuid` |
| YAML Parsing | `gopkg.in/yaml.v3` |
 
---
 
## Requirements
 
- **Docker** — containers are run via the local Docker daemon
- **Node.js + npm** — required for the web dashboard (`minik dashboard`)
---
 
## Installation
 
```bash
curl -fsSL https://github.com/SHIVAM-KUMAR-59/minikube/raw/main/install.sh -o /tmp/install.sh
chmod +x /tmp/install.sh
/tmp/install.sh
```
 
This downloads the `minik`, `minik-server`, and `minik-worker` binaries for your platform and places them in `/usr/local/bin`.
 
Supported platforms:
- macOS arm64 (Apple Silicon)
- macOS amd64 (Intel)
- Linux amd64
---
 
## Quick Start
 
```bash
# 1. Start the cluster with 2 worker nodes
minik cluster start --workers 2
 
# 2. Create a pod from a YAML spec
minik apply -f pod.yaml
 
# 3. Check pod status
minik get pods
 
# 4. Open the web dashboard
minik dashboard
 
# 5. Stop everything
minik cluster stop
```
 
**Example `pod.yaml`:**
 
```yaml
name: my-nginx
image: nginx
```
 
---
 
## CLI Reference
 
### Cluster
 
| Command | Description |
|---|---|
| `minik cluster start --workers N` | Start the server and N worker nodes as background processes |
| `minik cluster stop` | Stop all cluster processes |
 
### Resources
 
| Command | Description |
|---|---|
| `minik get pods` | List all pods with status and node assignment |
| `minik get nodes` | List all registered nodes with heartbeat time |
| `minik get services` | List all services |
| `minik apply -f <file>` | Create a pod from a YAML spec |
| `minik delete pod <id>` | Delete a pod by ID |
| `minik delete node <id>` | Delete a node by ID |
| `minik delete service <id>` | Delete a service by ID |
 
### Dashboard
 
| Command | Description |
|---|---|
| `minik dashboard` | Start the web dashboard and open it in the browser |
| `minik ping` | Check if the server is running |
 
---
 
## How It Works
 
**Pod lifecycle:**
 
1. User runs `minik apply -f pod.yaml` — CLI sends `POST /pods` to the control plane
2. Pod is saved to BoltDB with status `PENDING`
3. Scheduler goroutine ticks every 5 seconds, finds pending pods, picks a ready node via round-robin, and marks the pod `SCHEDULED`
4. Worker goroutine on the assigned node ticks every 5 seconds, finds scheduled pods for its node, pulls the Docker image, creates and starts the container, and marks the pod `RUNNING`
**Node health:**
 
Workers send a heartbeat to `POST /nodes/{id}/heartbeat` every 5 seconds. The scheduler only assigns pods to nodes with status `READY`. If a node stops sending heartbeats, it can be detected and marked `NOT_READY`.
 
**Service discovery:**
 
A service is a named endpoint that maps to a list of pod IDs. `GET /services/{name}/next` returns the next pod in round-robin order, enabling basic load balancing.
 
---
 
## API Reference
 
### Pods
 
| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/pods` | Create a pod |
| `GET` | `/pods` | List all pods |
| `DELETE` | `/pods/{id}` | Delete a pod |
| `PUT` | `/pods/{id}/status` | Update pod status |
 
### Nodes
 
| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/nodes/register` | Register a worker node |
| `POST` | `/nodes/{id}/heartbeat` | Send a heartbeat |
| `GET` | `/nodes` | List all nodes |
| `DELETE` | `/nodes/{id}` | Delete a node |
 
### Services
 
| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/services` | Create a service |
| `GET` | `/services` | List all services |
| `GET` | `/services/{name}/next` | Get next pod (load balanced) |
| `DELETE` | `/services/{id}` | Delete a service |
 
---
 
## Project Structure
 
```
minikube/
├── cmd/
│   ├── minik/                    ← CLI binary entry point
│   │   ├── main.go
│   │   └── cmd/
│   │       ├── root.go
│   │       ├── ping.go
│   │       ├── apply.go
│   │       ├── dashboard.go
│   │       ├── get.go
│   │       ├── delete.go
│   │       ├── cluster.go
│   │       ├── get/              ← get subcommands
│   │       │   ├── pods.go
│   │       │   ├── nodes.go
│   │       │   └── services.go
│   │       ├── delete/           ← delete subcommands
│   │       │   ├── pods.go
│   │       │   ├── nodes.go
│   │       │   └── services.go
│   │       └── cluster/          ← cluster subcommands
│   │           ├── start.go
│   │           └── stop.go
│   ├── server/                   ← API server binary
│   │   └── main.go
│   └── worker/                   ← Worker node binary
│       └── main.go
├── internal/
│   ├── api/                      ← HTTP handlers
│   │   ├── handler.go
│   │   ├── pod_handler.go
│   │   ├── node_handler.go
│   │   ├── service_handler.go
│   │   └── ping_handler.go
│   ├── store/                    ← BoltDB persistence
│   │   ├── db.go
│   │   ├── pod.go
│   │   ├── node.go
│   │   ├── service.go
│   │   └── status.go
│   ├── scheduler/                ← Pod scheduling loop
│   │   └── scheduler.go
│   ├── worker/                   ← Container execution
│   │   └── worker.go
│   └── loadbalancer/             ← Round-robin load balancer
│       └── loadbalancer.go
├── dashboard/                    ← Next.js + Tailwind UI
├── docs/                         ← Learning notes
├── Makefile
├── install.sh
├── go.mod
└── README.md
```
 
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

### Phase 3 - Worker node and Docker container lifecycle
- **Goal**: Complete the pod lifecycle by actually running containers on worker nodes using the Docker SDK.
- **What we built**: A worker node that runs as a background goroutine, reconciles scheduled pods, pulls Docker images, creates and starts real containers, and updates pod status to `RUNNING` in the store.
- **Learnings**: Docker SDK (`ImagePull`, `ContainerCreate`, `ContainerStart`), `context.Background()` and why context is needed for long-running operations, aliasing imports to avoid naming conflicts, and chaining goroutine-based reconciliation loops.
- **Deliverable**: `POST /pods` with an image like `nginx` results in a real Docker container running on the machine within 10 seconds. `docker ps` shows the container and `GET /pods` shows status `RUNNING`.

### Phase 4 - Service Discovery and Load Balancing
- **Goal**: Allow pods to be grouped under named services and have traffic distributed across them.
- **What we built**: A `Service` data structure, service store methods on the existing `Store`, service API endpoints (`POST /services`, `GET /services`), a round-robin load balancer, and a `GET /services/{name}/next` endpoint that returns the next pod for a given service.
- **Learnings**: Chi URL parameters, separating concerns across handler files, round-robin load balancing with a per-service counter map, and why a single BoltDB connection must be shared across all store operations.
- **Deliverable**: Create a service pointing to running pods and hit `/services/{name}/next` repeatedly — each call returns the next pod in round-robin order.

### Phase 5 - Multi-node support and HTTP-based worker communication
- **Goal**: Separate the worker into its own binary so multiple independent worker nodes can run on different processes, making the architecture properly distributed.
- **What we built**: A standalone `cmd/worker/main.go` binary that accepts `--node-id` and `--server-url` flags, refactored the worker to communicate with the control plane purely over HTTP (no direct DB access), and added a `PUT /pods/{id}/status` endpoint so workers can update pod state remotely.
- **Learnings**: Why a worker shouldn't have direct database access in a distributed system, using `http.NewRequest` for PUT requests, the `flag` package for CLI flags, and how proper process separation makes a system feel real.
- **Deliverable**: Run `cmd/server/main.go` and two instances of `cmd/worker/main.go` with different node IDs — pods get distributed across both workers via round-robin scheduling and run as real Docker containers on their assigned node.

### Phase 6 - CLI commands
- **Goal**: Make MiniKube usable from the terminal without curl commands, and allow anyone to start the entire cluster with a single command.
- **What we built**: `minik get pods/nodes/services` for listing resources, `minik delete pod/node/service <id>` for deletion, `minik apply -f pod.yaml` for creating pods from YAML specs, and `minik cluster start/stop` for managing the entire cluster as background processes.
- **Learnings**: YAML parsing with `gopkg.in/yaml.v3`, running detached background processes with `os/exec` and `syscall.SysProcAttr{Setsid: true}`, saving and reading PID files to track running processes, and structuring nested CLI commands with Cobra using parent/child command packages.
- **Deliverable**: A user can now run `minik cluster start --workers 2` to spin up the entire cluster, use `minik apply -f pod.yaml` to create pods, monitor with `minik get pods`, and shut everything down with `minik cluster stop` — no curl commands or manual terminal management needed.

### Phase 7 - Dashboard UI and command
- **Goal**: Give users a visual interface to interact with the cluster without needing terminal commands or curl.
- **What we built**: A full Next.js + Tailwind dashboard with four pages — Overview (cluster stats, recent pods, nodes, services), Pods (table with create and delete), Nodes (table with heartbeat timestamps), and Services (cards with create and delete). Added a `minik dashboard` CLI command that starts the frontend server and opens the browser automatically, with cross-platform support for Mac, Windows, and Linux.
- **Learnings**: Converting Next.js server components to client components for polling, `useEffect` intervals for live data refresh every 5 seconds, relative timestamps with `date-fns`, cross-platform browser opening with `runtime.GOOS`, and keeping a child process alive from a CLI command.
- **Deliverable**: Run `minik dashboard` and the browser opens automatically at `localhost:3000` showing live cluster state that auto-refreshes every 5 seconds. Users can create and delete pods and services directly from the UI.
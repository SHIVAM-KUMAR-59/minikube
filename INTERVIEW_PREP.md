# MiniKube — Interview Preparation Guide

A comprehensive set of questions, answers, key terminology, and tradeoffs for discussing the MiniKube lightweight container orchestrator in technical interviews.

## Table of Contents

- [System Design Questions](#system-design-questions)
- [Go & Backend Questions](#go--backend-questions)
- [Distributed Systems Questions](#distributed-systems-questions)
- [Scheduling & Orchestration Questions](#scheduling--orchestration-questions)
- [Container & Docker Questions](#container--docker-questions)
- [Database & State Questions](#database--state-questions)
- [CLI & Tooling Questions](#cli--tooling-questions)
- [Frontend & Dashboard Questions](#frontend--dashboard-questions)
- [Key Terminology](#key-terminology)
- [Key Tradeoffs Made](#key-tradeoffs-made)
- [Behavioural Questions](#behavioural-questions)

---

## System Design Questions

**Q: Walk me through the high-level architecture of MiniKube.**

The system has three main components. The control plane is a Go HTTP server that owns all cluster state — it exposes a REST API, runs a scheduler goroutine, and runs a controller goroutine. Workers are separate Go binaries that register with the control plane, receive pod assignments, manage Docker containers, and expose their own HTTP server for log streaming. The CLI is a third binary that communicates with the control plane over HTTP and provides a kubectl-style interface. All state is persisted in BoltDB, an embedded key-value store that runs inside the server process — no separate database needed.

---

**Q: How does MiniKube handle state? Why did you choose BoltDB over something like PostgreSQL or Redis?**

BoltDB is an embedded key-value store — it's a single file on disk, no separate process to manage. This was a deliberate choice to make MiniKube self-contained. A user installing MiniKube shouldn't need to run a database server. BoltDB gives you ACID transactions, persistence, and reasonable performance for the scale MiniKube operates at. The tradeoff is that it can't scale horizontally — only one process can open the database at a time, which is fine since the control plane is a single process. For a production orchestrator you'd use etcd, which is also a key-value store but distributed.

---

**Q: How would you scale MiniKube if you needed to support thousands of nodes?**

Several things would need to change. First, replace BoltDB with a distributed store like etcd — the real Kubernetes uses etcd for exactly this reason. Second, the scheduler would need to become smarter — round-robin breaks down at scale, you'd want resource-aware scheduling considering CPU and memory. Third, the control plane would need to be horizontally scaled with leader election so multiple API server instances can run simultaneously. Fourth, the heartbeat interval and controller check frequency would need tuning — at 10,000 nodes, checking every node every 10 seconds with a single goroutine would create bottlenecks. You'd partition the node health checks across multiple controller instances.

---

**Q: How does MiniKube differ from Kubernetes architecturally?**

The concepts map closely — control plane, scheduler, controller, worker nodes (kubelet in Kubernetes), state store. The key differences are scale and complexity. Kubernetes has namespaces, RBAC, multi-tenancy, custom resource definitions, admission controllers, and a pluggable runtime interface. MiniKube has none of these. Kubernetes workers (kubelets) talk to the API server over a more complex protocol and manage pod lifecycles including init containers, sidecars, and health probes. MiniKube workers are simpler — they poll for scheduled pods and run Docker containers. Kubernetes uses etcd as its distributed state store; MiniKube uses BoltDB which is single-node only. The reconciliation loop concept is identical — both systems continuously compare desired state against actual state and take action to converge them.

---

**Q: What happens when the control plane goes down?**

Running pods continue to run — Docker containers are independent of the control plane. Workers keep sending heartbeats but they fail silently since the server is unreachable. The scheduler and controller stop running since they live inside the server process. New pods can't be created, existing pods can't be rescheduled if a worker dies. When the control plane comes back up, workers re-register, heartbeats resume, and the controller reconciles state from BoltDB which was persisted before the crash. This is similar to how Kubernetes handles API server downtime — workloads keep running, but no new scheduling decisions can be made.

---

## Go & Backend Questions

**Q: Why did you choose Go for this project?**

Go is the language Kubernetes, Docker, and most cloud-native infrastructure is written in. It has first-class concurrency with goroutines and channels, compiles to a single static binary with no runtime dependencies, has a strong standard library for HTTP and networking, and is fast enough for systems work. For a project about orchestration, using the same language as the tools it's inspired by made sense both technically and as a learning exercise.

---

**Q: How did you use goroutines in MiniKube?**

Every background loop runs as a goroutine. The scheduler runs a `time.NewTicker` goroutine that fires every 5 seconds and checks for pending pods. The controller runs another goroutine that fires every 10 seconds and checks node heartbeats. Each worker runs two goroutines — one for pod reconciliation and one for heartbeats. The worker's HTTP server runs in a goroutine so it doesn't block the main process. All of these run concurrently without explicit thread management — Go's runtime schedules them on available OS threads.

---

**Q: How do you handle concurrent access to shared state in Go?**

BoltDB handles concurrent access at the database level — it uses a single writer, multiple reader model with file-level locking. This means only one write transaction can happen at a time, which is fine for MiniKube's scale. The scheduler and controller both read and write pods and nodes, but since they go through BoltDB transactions, consistency is guaranteed. If I were to add in-memory caching, I'd use `sync.RWMutex` to protect the cache from concurrent reads and writes.

---

**Q: Explain how `context.Background()` is used in MiniKube.**

The Docker SDK requires a `context.Context` for every operation — image pulls, container creates, container starts, container inspects. `context.Background()` creates a context with no cancellation, no deadline, and no timeout. It's the right choice for operations that should run to completion regardless of external signals. In a production system you'd want to add timeouts — for example, `context.WithTimeout(context.Background(), 30*time.Second)` for image pulls so a slow registry doesn't block the worker indefinitely.

---

**Q: How does the worker communicate with the control plane?**

Purely over HTTP. The worker calls `POST /nodes/register` on startup, `POST /nodes/{id}/heartbeat` every 5 seconds, `GET /pods` to find scheduled pods, and `PUT /pods/{id}/status` to update pod status after starting a container. There's no direct database access from the worker — the control plane is the single source of truth. This is the correct distributed systems pattern — the database is an implementation detail of the control plane, not something workers should know about.

---

## Distributed Systems Questions

**Q: How does MiniKube detect node failures?**

Each worker sends a heartbeat by calling `POST /nodes/{id}/heartbeat` every 5 seconds. A controller goroutine runs inside the server every 10 seconds and checks `time.Since(node.LastHeartbeat)` for every registered node. If a node's last heartbeat is more than 15 seconds ago, the controller marks it `NOT_READY` and finds all pods assigned to that node with status `RUNNING` or `SCHEDULED` and resets them to `PENDING`. The scheduler then picks them up on the next tick and assigns them to a healthy node.

---

**Q: What are the failure modes in your heartbeat system?**

There are a few. First, false positives — a node could be healthy but temporarily unable to reach the server due to a network partition. The controller would mark it `NOT_READY` and reschedule its pods, but the containers are still running. This creates a split-brain scenario where the same logical pod is running on two nodes. In production you'd handle this with fencing — explicitly stopping the containers on the old node before starting them on the new one. Second, the 15-second threshold is arbitrary — too short causes false positives under load, too long means slow failure detection. Third, if the controller crashes, heartbeat checking stops entirely.

---

**Q: What is a reconciliation loop and why is it important?**

A reconciliation loop is a background process that continuously compares the desired state of the system (what's stored in the database) against the actual state (what's actually running) and takes action to make them converge. It's the core pattern behind Kubernetes and MiniKube. The scheduler reconciles `PENDING` pods by assigning them to nodes. The controller reconciles node health by marking dead nodes and resetting their pods. The worker reconciles `SCHEDULED` pods by running containers and `RUNNING` pods by restarting crashed containers. The key insight is that the system is eventually consistent — you don't need atomic distributed transactions, you just need to keep running the reconciliation loop until everything converges.

---

**Q: How does MiniKube handle the case where a pod is rescheduled but the old container is still running?**

This is a known limitation. When a node is marked `NOT_READY`, the controller resets pods to `PENDING` and they get rescheduled to a new node. The old containers on the dead node keep running — Docker doesn't know about the orchestrator's decision. If the old node comes back online, it would see those pods as `RUNNING` in Docker but `PENDING` in the database and try to start new containers, hitting a name conflict. The correct fix is to fence the old node — explicitly stop its containers before rescheduling. I handled the name conflict partially by appending the pod ID to container names, but full fencing would require the controller to call the worker's HTTP server to stop containers before resetting pods.

---

**Q: How does service discovery work in MiniKube?**

A service is a named resource that maps to a list of pod IDs. When a client calls `GET /services/{name}/next`, the load balancer looks up the service, gets its pod list, uses a per-service counter with modulo to pick the next pod in round-robin order, and returns that pod's details. Pods find each other by service name rather than IP address. This is analogous to Kubernetes Services — an abstraction layer over pod IPs that survives pod restarts. The limitation is that MiniKube's services are purely informational — they return a pod to route to, but don't actually proxy traffic. A real implementation would use iptables or a sidecar proxy.

---

## Scheduling & Orchestration Questions

**Q: How does the MiniKube scheduler work?**

The scheduler runs a goroutine with a 5-second ticker. On each tick it calls `GetAllPods()` from the store and filters for pods with status `PENDING`. For each pending pod it calls `GetAllNodes()` and filters for nodes with status `READY`. If there are ready nodes, it picks one using a counter modulo the number of ready nodes — round-robin — updates the pod's `NodeID` and status to `SCHEDULED`, and saves it back to the store. The worker on the assigned node picks it up on its next reconciliation tick.

---

**Q: What are the limitations of round-robin scheduling?**

Round-robin is simple and fair in terms of pod count, but it's blind to resource utilization. Node 1 might be running 3 CPU-intensive containers while node 2 is idle — round-robin would still assign the next pod to node 1. A better scheduler would consider available CPU and memory on each node before assigning. Kubernetes's default scheduler does this — it scores nodes based on resource requests, affinity rules, taints, and tolerations. For MiniKube's scope round-robin is appropriate, but it's a clear extension point.

---

**Q: How does pod restart policy work?**

The worker's reconciliation loop handles this. After processing scheduled pods, it loops through all pods with status `RUNNING` assigned to its node. For each one, it calls `ContainerInspect` on the Docker client using the container name. If the inspect call fails (container doesn't exist) or the container state is not `running`, the worker calls `removeContainer` to clean up the old container, then calls `RunPod` to start a fresh one. This gives MiniKube automatic restart behavior similar to Kubernetes's `restartPolicy: Always`.

---

**Q: How do replicas work in MiniKube?**

When a pod is created with `replicas: N`, the `CreatePod` handler loops N times and creates N separate pod records in the database, each with a unique ID and a name like `nginx-1`, `nginx-2`, `nginx-3`. Each pod is an independent entity that gets scheduled and run independently. The scheduler distributes them across nodes via round-robin. This is similar to how Kubernetes ReplicaSets work, except MiniKube has no concept of a ReplicaSet controller watching the desired replica count — if you delete `nginx-2`, MiniKube doesn't automatically create a replacement.

---

## Container & Docker Questions

**Q: How does MiniKube interact with Docker?**

Through the official Docker SDK for Go — `github.com/docker/docker/client`. The worker creates a Docker client using `client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())` which reads the Docker socket from environment variables. For each pod, the worker calls `ImagePull` to pull the image, `ContainerCreate` to create the container with the image and name, and `ContainerStart` to start it. For log streaming it calls `ContainerLogs` which returns an `io.ReadCloser`. For health checking it calls `ContainerInspect` and checks `State.Running`.

---

**Q: Why do you use `stdcopy.StdCopy` instead of `io.Copy` for Docker logs?**

Docker multiplexes stdout and stderr into a single stream with an 8-byte header on each log line indicating which stream it came from and the payload length. If you use `io.Copy` directly, those header bytes appear as garbage characters at the start of each line. `stdcopy.StdCopy` understands the multiplexing format and strips the headers, writing stdout and stderr to separate writers. In MiniKube both writers are the HTTP response, so all logs appear cleanly in order.

---

**Q: How does log streaming work end to end?**

The flow has three hops. First, the CLI calls `GET /pods/{name}/logs` on the control plane. The control plane looks up the pod by name, finds its node, and makes a `GET /logs/{containerName}` request to that worker's HTTP server. The worker calls `ContainerLogs` on the Docker client and uses `stdcopy.StdCopy` to write the logs into the HTTP response. The control plane streams that response back to the CLI using `io.Copy`. The CLI scans the response line by line and colorizes output based on log level keywords. This proxy pattern means the control plane doesn't need a Docker client — all Docker interaction is encapsulated in the worker.

---

**Q: Why do container names include part of the pod ID?**

To avoid name conflicts when pods are rescheduled. Docker requires container names to be unique — if a pod is rescheduled from node1 to node2 after node1 fails, the worker on node2 would try to create a container with the pod's name. If the original container on node1 still exists (node1 came back online), Docker would reject the creation with a conflict error. By appending the first 8 characters of the pod UUID to the container name (`nginx-abc12345`), each container has a unique name even across rescheduling events.

---

## Database & State Questions

**Q: Why BoltDB instead of etcd?**

etcd is a distributed key-value store designed for high availability — it requires a cluster of nodes and handles leader election and consensus via Raft. It's the right choice for production Kubernetes. BoltDB is an embedded single-file database — simpler to use, no cluster to manage, zero configuration. For a single-machine learning project that needs persistence without external dependencies, BoltDB is the right tradeoff. The API is similar — both are key-value stores organized into buckets/namespaces. Migrating from BoltDB to etcd would be a matter of swapping the store layer.

---

**Q: How is data organized in BoltDB?**

BoltDB organizes data into buckets, which are like tables or namespaces. MiniKube has three buckets — `pods`, `nodes`, and `services`. Within each bucket, records are stored as key-value pairs where the key is the resource ID and the value is the JSON-serialized struct. All reads use `db.View` transactions and all writes use `db.Update` transactions. BoltDB guarantees that `db.Update` transactions are atomic — either all writes succeed or none do.

---

**Q: How do you handle the case where two goroutines try to write to the database simultaneously?**

BoltDB serializes write transactions — only one `db.Update` can run at a time. If a second goroutine tries to open a write transaction while one is already in progress, it blocks until the first one completes. This is fine for MiniKube's scale — the scheduler and controller run on 5-10 second intervals so write contention is minimal. In a high-throughput system you'd want a write-ahead log or a database that supports concurrent writes.

---

## CLI & Tooling Questions

**Q: How is the CLI structured and why did you use Cobra?**

The CLI uses Cobra, the same library `kubectl` and most Go CLI tools use. Cobra organizes commands into a tree — `minik` is the root, `cluster` is a parent command with `start` and `stop` children, `get` is a parent with `pods`, `nodes`, `services` children. Each command is a separate file in its own package, which keeps the code organized as the command set grows. Cobra handles flag parsing, help text generation, and command routing automatically.

---

**Q: How does `minik cluster start` work under the hood?**

`cluster start` uses `os/exec` to launch `minik-server` and N instances of `minik-worker` as detached background processes. Each process is started with `cmd.Start()` rather than `cmd.Run()` so they run independently without waiting for them to finish. `syscall.SysProcAttr{Setsid: true}` detaches each process from the CLI's process group so they survive after the CLI exits. The PID of each process is saved to `~/.minik/cluster.pid` so `minik cluster stop` can kill them later by reading the PID file and calling `syscall.Kill`.

---

**Q: How does `minik apply -f minik.yaml` work?**

The command reads and parses the YAML file. If it contains a `pods` or `cluster` key it's treated as a full cluster config — it checks if the cluster is already running by pinging `localhost:8080/ping`, starts it if not using `os.Executable()` to find the binary path, waits for N nodes to register by polling `GET /nodes` with a spinner, then creates all pods via `POST /pods` and all services via `POST /services`. If the YAML only has `name` and `image` keys it's treated as a simple pod spec for backward compatibility.

---

**Q: Why use `os.Executable()` instead of hardcoding the binary path?**

`os.Executable()` returns the path of the currently running binary. When a user installs MiniKube to `/usr/local/bin/minik`, `os.Executable()` returns `/usr/local/bin/minik`, so `minik-server` and `minik-worker` are found at `/usr/local/bin/minik-server` and `/usr/local/bin/minik-worker`. When running from the project root with the local build, it returns `./minik` and finds `./minik-server` and `./minik-worker` in the same directory. This makes the binary relocatable — it works correctly wherever it's installed without hardcoded paths.

---

## Frontend & Dashboard Questions

**Q: How does the dashboard get its data?**

The Next.js dashboard is a client-side application that polls the MiniKube REST API every 5 seconds using `useEffect` with `setInterval`. Each page fetches its relevant data — the pods page calls `GET /pods`, the nodes page calls `GET /nodes`, the services page calls `GET /services`. The overview page calls all three in parallel using `Promise.all`. CORS is configured on the Go server to allow requests from `localhost:3000`.

**Q: Why polling instead of WebSockets?**

Polling is simpler to implement and reason about. WebSockets would give lower latency updates but add complexity — connection management, reconnection logic, server-side event broadcasting. For a dashboard where 5-second staleness is acceptable, polling is the right tradeoff. The dashboard already felt live enough with 5-second polling.

---

**Q: How is the dashboard distributed to users?**

The dashboard is a Next.js app that runs as a dev server. When a user runs `minik dashboard`, it finds the dashboard directory relative to the installed binary, runs `npm run dev` as a child process, waits 3 seconds for the server to start, then opens the browser. The dashboard directory is copied to `/usr/local/bin/dashboard` during `make install` with `npm install` run post-copy. The install script does the same when installing from GitHub releases.

---

## Key Terminology

**Control Plane** — the server component that owns cluster state, runs the scheduler and controller, and exposes the REST API. Analogous to the Kubernetes API server + scheduler + controller manager.

**Worker Node** — a separate process that registers with the control plane, receives pod assignments, manages Docker containers, and sends heartbeats. Analogous to the Kubernetes kubelet.

**Pod** — the basic unit of scheduling in MiniKube. Represents a single container with a name, image, status, and node assignment.

**Reconciliation Loop** — a background goroutine that continuously compares desired state against actual state and takes action to converge them. The fundamental pattern behind Kubernetes and MiniKube.

**Heartbeat** — a periodic signal sent by workers to the control plane to indicate they are alive. Missing heartbeats trigger the controller to mark a node `NOT_READY`.

**BoltDB** — an embedded key-value store used as MiniKube's state store. Stores all pods, nodes, and services as JSON in named buckets. Analogous to etcd in Kubernetes.

**Scheduler** — a goroutine that assigns `PENDING` pods to `READY` nodes using round-robin. Runs every 5 seconds inside the control plane.

**Controller** — a goroutine that watches node health and reschedules pods from dead nodes. Runs every 10 seconds inside the control plane. Analogous to the Kubernetes node controller.

**Service** — a named resource that maps a name and port to a list of pod IDs, enabling service discovery and round-robin load balancing.

**Replica** — a copy of a pod. `replicas: N` creates N independent pod records that are scheduled and run independently across worker nodes.

**Round-robin** — a scheduling strategy that distributes pods across nodes in sequence. Simple and fair in terms of count but blind to resource utilization.

**Fencing** — the act of explicitly stopping containers on a dead node before rescheduling its pods elsewhere. Prevents split-brain scenarios where the same logical pod runs on two nodes simultaneously. Not yet implemented in MiniKube.

---

## Key Tradeoffs Made

**BoltDB vs etcd** — chose BoltDB for simplicity and zero external dependencies. The tradeoff is single-node only — can't scale the control plane horizontally. etcd would be the correct choice for a production system.

**HTTP vs gRPC for worker communication** — chose HTTP for simplicity and debuggability. gRPC would give better performance, strongly-typed contracts, and built-in streaming. A natural Phase 10 upgrade.

**Round-robin vs resource-aware scheduling** — chose round-robin for simplicity. Resource-aware scheduling would require workers to report CPU and memory usage and the scheduler to consider it. This is the most impactful missing feature for real-world usefulness.

**Polling vs WebSockets for dashboard** — chose polling for simplicity. WebSockets would give real-time updates with lower latency but add connection management complexity.

**No fencing on node failure** — when a node is marked `NOT_READY`, pods are rescheduled but old containers aren't stopped. This risks split-brain. The correct fix requires the controller to call the worker's HTTP server to stop containers before resetting pods.

**Single binary install vs Docker Compose** — chose a single CLI binary with embedded server/worker management over Docker Compose. This removes the Docker Compose dependency and gives a better UX, but means the install is more complex and the dashboard requires Node.js.

**Node.js dependency for dashboard** — chose Next.js dev server over embedding static files into the Go binary. The tradeoff is requiring Node.js on the user's machine. The alternative — `go:embed` with a static Next.js export — was attempted but abandoned due to routing complexity with chi.

---

## Behavioural Questions

**Q: Tell me about a technical challenge you faced building MiniKube.**

The hardest problem was the log streaming architecture. Initially I gave the server a Docker client so it could fetch logs directly, but that violated the architectural boundary — the server shouldn't know about Docker, that's the worker's job. I refactored to give each worker its own HTTP server on a separate port. The control plane now proxies log requests to the correct worker based on which node the pod is assigned to. This required adding an `Address` field to the `Node` struct so the control plane knows where to find each worker's HTTP server. The result was architecturally cleaner and closer to how Kubernetes actually handles log requests — the API server proxies to the kubelet.

---

**Q: What would you do differently if you started MiniKube over?**

I'd design the worker communication protocol more carefully upfront. Starting with HTTP was fine, but I'd think earlier about where gRPC would fit — particularly for streaming log data and for the heartbeat mechanism. I'd also add resource tracking from the beginning — workers reporting available CPU and memory, and the scheduler considering it. This would make the scheduling decisions more realistic. Finally I'd write tests alongside the code rather than after — the reconciliation loops and state transitions are complex enough that unit tests would have caught several bugs earlier.

---

**Q: How did building MiniKube improve your understanding of Kubernetes?**

Before MiniKube, Kubernetes felt like magic — pods appear, containers run, failures self-heal. Building MiniKube demystified every piece. The reconciliation loop is just a goroutine with a ticker. The scheduler is just a function that reads pending pods and assigns nodes. Node health detection is just checking a timestamp. Service discovery is just a round-robin counter over a list of pod IDs. Understanding these fundamentals makes reading Kubernetes source code and documentation much easier — I now know what problem each component is solving and why the design decisions were made.

---

**Q: How did you approach testing and validating MiniKube?**

Most testing was integration testing — running the full system and verifying behavior end to end. I'd start the server and two workers, create pods, verify they reached `RUNNING` status, then kill a worker and verify pods were rescheduled. For the restart policy I'd manually stop a Docker container and verify it restarted within 10 seconds. The codebase would benefit from more unit tests — particularly for the scheduler's round-robin logic, the controller's heartbeat threshold calculation, and the store layer's CRUD operations. This is an acknowledged gap.

---

**Q: What would Phase 10 of MiniKube look like?**

The most impactful additions would be gRPC for worker-to-server communication, TLS for encrypted communication between components, and resource-aware scheduling where workers report CPU and memory usage and the scheduler considers it when placing pods. Beyond that, rolling updates — updating pod images one replica at a time with zero downtime — would make MiniKube genuinely useful for deployments. Persistent volumes, environment variables in pod specs, and a proper health check mechanism (liveness and readiness probes) would round out the feature set.

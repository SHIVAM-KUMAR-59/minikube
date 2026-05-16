### Docker SDK

### What is Docker SDK?
- The Docker SDK for Go allows you to interact with the Docker daemon programmatically using Go.
- Official documentation: [here](https://pkg.go.dev/github.com/docker/docker/client)
- Installation command: `go get github.com/docker/docker@v28.5.2+incompatible`

### Key Concepts
- **Docker Client** — the entry point for all Docker operations. Created once and reused across the application.
- **Container lifecycle** — creating and starting a container are always two separate steps.
- **Context** — every Docker SDK operation requires a `context.Context`. Use `context.Background()` for operations with no timeout or cancellation.

---

### Creating a Docker Client
```go
dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
```
- `client.FromEnv` — configures the client from environment variables (uses your local Docker daemon by default)
- `client.WithAPIVersionNegotiation()` — automatically negotiates the API version with the daemon

---

### Container Lifecycle

#### Step 1 — Pull the image
```go
reader, err := dockerClient.ImagePull(ctx, "nginx", image.PullOptions{})
io.Copy(io.Discard, reader) // must read to completion or pull won't finish
reader.Close()
```

#### Step 2 — Create the container
```go
container, err := dockerClient.ContainerCreate(ctx, &container.Config{
    Image: "nginx",
}, nil, nil, nil, "my-container-name")
```

#### Step 3 — Start the container
```go
err = dockerClient.ContainerStart(ctx, container.ID, container.StartOptions{})
```

#### Step 4 — Stop the container
```go
err = dockerClient.ContainerStop(ctx, container.ID, container.StopOptions{})
```

---

### Important Points
- Always read and close the `ImagePull` response body — it's an `io.ReadCloser` representing a live stream. If you don't read it to completion, the pull won't finish.
- `ContainerCreate` returns a container ID — always use this ID (not the name) for subsequent operations like `ContainerStart` and `ContainerStop`.
- Import aliasing is sometimes needed to avoid conflicts: `dockerContainer "github.com/docker/docker/api/types/container"`
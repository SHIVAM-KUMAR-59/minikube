# Bolt DB

### What is Bolt DB?
- It is a pure Go key/value store.
- Its goal is to provide a simple, fast, and reliable database for projects that don't require a full database server such as Postgres or MySQL.
- Since Bolt is meant to be used as such a low-level piece of functionality, simplicity is key. The API will be small and only focus on getting values and setting values. That's it.
- Official documentation [here](https://github.com/etcd-io/bbolt)
- Installation command: `go get go.etcd.io/bbolt@latest`

### Why Bolt DB instead of SQL DB?
- It is **embedded**, i.e. the database runs inside your Go process. there is no separate database server to install, start, or connect to. BoltDB is just a single `.db` file on disk. Your Go program opens that file directly and reads/writes to it.
Compare that to Postgres where you have:
    - A separate Postgres server process running
    - Your app connecting to it over a network/socket 
    - A connection string, credentials, etc.

- Bolt is good for read intensive workloads. Sequential write performance is also fast but random writes can be slow.

- This is perfect for MiniKube because we don't want to make someone install a database just to run our tool. Everything is self-contained in one binary + one file.

- This is also exactly how etcd works conceptually — Kubernetes' real state store is also an embedded key-value store, just a more sophisticated one.

### Key Concepts
- **Transactions**: Bolt allows only one read-write transaction at a time but allows as many read-only transactions as you want at a time. Each transaction has a consistent view of the data as it existed when the transaction started.

- **Buckets**: Buckets are collections of key/value pairs within the database. All keys in a bucket must be unique. You can create a bucket using the `Tx.CreateBucket()` or `Tx.CreateBucketIfNotExists()` function.

### Key Commands
- Storing a key-value pair: `Bucket.Put(key, value)`
- Fetching a key-value pair: `Bucket.Get(key)`
    - The `Get()` function does not return an error because its operation is guaranteed to work (unless there is some kind of system failure). If the key exists then it will return its byte slice value. If it doesn't exist then it will return `nil`
- Deleting a key: `Bucket.Delete(key)`
# DistributedTwitter
Make do Twitter clone with a Raft based distributed key-value data store.

Project requirements are stated in `prompt.md`

### Steps to run:

1. Run `go mod download` to download all dependencies
2. Compile proto files: `./scripts/protoc.sh`
2. Start etcd instances if using Raft as the storage implementation (Hint: look into `scripts/initEtcd.sh` file)
3. Check the server config (`cmd/server/config.yaml`)
4. Check the client config (`cmd/web/config.yaml`)
5. Start the server: `go run cmd/server/twitter.go`
6. Start the client: `go run cmd/web/web.go`
7. Navigate to `localhost:3000` (or depending on your config.yaml in client)


Current storage implementations:

1. Memory: In-memory storage for users and posts, non-persistent.
2. etcd: Distributed key-value based persistent storage

TODO Implementations:

1. Zookeeper
2. PostgreSQL
3. Cassandra
4. Consul

TODO:

1. Better unit testing
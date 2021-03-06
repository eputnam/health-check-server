# health-check-server
We're building a backend for team health checks in golang!

# Create a config file
Copy [config.sample.yaml] to `config.yaml` and adjust options accordingly.
```shell
cp config.sample.yaml config.yaml
```

# Run the server
Run this command in the project root
```shell
go run ./main.go
```

# Docker

## The server
Build the Docker image (optionally set TAG, defaults to latest)
```shell
TAG=my-tag make server-build
```

Run the Docker container (optionally set TAG, defaults to latest)
```shell
TAG=my-tag make server-run
```

## PostgreSQL
You can run a local PostgreSQL server with these make targets:
```shell
make db-run
make db-stop
```

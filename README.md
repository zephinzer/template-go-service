# Template Go Service

This is a template repository that provides boilerplate code for a service written in Go.

**Supported data stores**
- Kafka
- MongoDB
- MySQL
- NATS
- PostgreSQL
- Redis

**Supported deployment methods**
- Kubenretes via Helm

# Boilerplate information

## Directory structure

| Directory path | Description |
| --- | --- |
| `/.data` | Persistent data storage |
| `/bin` | Contains built binaries |
| `/cmd` | Contains command entrypoints |
| `/deploy` | Contains deployment manifests |
| `/internal` | Contains app-related packages |
| `-> ./api` | Contains HTTP-related packages and defines mainly routing |
| `-> ./constants` | Contains build configurations and app-wide constants |
| `-> ./database` | Contains database connection management code |
| `-> ./docs` | Contains Swagger documentation |
| `-> ./example` | Contains an example endpoint that connects to the HTTP API |
| `-> ./server` | Contains the server instance code |
| `/tests` | Contains test artifacts |
| `/vendor` | Contains dependenices |

## CLI information

To test the CLI you can run `go run ./cmd/app`, help information should be available.

The following commands are included in the boilerplate which you can remove if not applicable to your project:

- `app debug kafka` - tests the connection to Kafka (run `make kafka-jks` and `make start-kafka` before this)
- `app debug mongo` - tests the connection to MongoDB (run `make start-mongo` before this)
- `app debug mysql` - tests the connection to MySQL (run `make start-mysql` before this)
- `app debug nats` - tests the connection to NATS (run `make start-nats` before this)
- `app debug postgres` - tests the connection to PostgreSQL (run `make start-postgres` before this)
- `app debug redis` - tests the connection to Redis (run `make start-redis` before this)
- `app start job` - sample command to start a one-time job
- `app start server` - sample command to start the server
- `app start worker` - sample command to start a worker service


## Notable packages and software used

### Go

1. [`cobra`](https://github.com/spf13/cobra) for structuring the app's CLI invocations and flags configuration
1. [`viper`](https://github.com/spf13/viper) for managing configuration on a service level
1. [`gofiber`](https://github.com/gofiber/fiber) for the HTTP server and HTTP API scaffolding
2. [`swaggo`](https://github.com/swaggo/swag) to generate Swagger documentation

### Tooling

1. [Docker](https://www.docker.com/) for build/release packaging
2. [Docker Compose](https://docs.docker.com/compose/) for bringing up development services
3. [KinD](https://github.com/kubernetes-sigs/kind/) for bringing up a local Kubernetes cluster
4. [`nk`](https://github.com/nats-io/nkeys) for generating NATS nkeys
5. [Make](https://www.gnu.org/software/make/) for development operations recipes

### Supported data stores

#### Queues/Message brokers

1. Kafka
5. NATS for queueing/message brokering

#### Relational databases
7. MySQL for RDBMS usage
8. PostgreSQL for RDBMS usage

#### NoSQL/document databases/Key stores
9.  MongoDB for NoSQL usage
10. Redis for caching

## SDLC Usage

### Development operations

This repository uses Makefile for storing development operations. Included operations are:

- `make binary` builds the binary and stores it in `./bin/`
- `make deploy-kind` deploys the application onto a local Kubernetes cluster
- `make deploy-k8s` deploys the application onto the currently selected Kubernetes cluster
- `make deps` pulls in the dependencies into `./vendor/`
- `make docs-swaggo` generates the Swagger documentation using `swaggo`
- `make image` builds the Docker image
- `make install-swaggo` installs `swaggo` at the latest version
- `make kafka-jks` generates the keys and certificates required for Kafka security
- `make kind-load` loads the image into the local Kubernetes cluster
- `make nats-nkey` creates a new nkey for use with NATS
- `make start` starts the HTTP server using defaults
- `make start-kafka` starts a local Kafka instance via Docker Compose
- `make start-kind` starts a local Kubernetes cluster with KinD
- `make start-mongo` starts a local MongoDB instance
- `make start-mysql` starts a local MySQL instance
- `make start-nats` starts a local NATS instance
- `make start-postgres` starts a local PostgreSQL instance
- `make start-redis` starts a local Redis instance
- `make test` runs tests on the application and produces coverage artifacts at `./tests/`

To configure the Makefile recipes, create a new `Makefile.properties` in the root of this repository. Configurable values can be found in the `Makefile` above the line `-include Makefile.properties`.

Some common configurations are documented below. Insert the code snippets into `Makefile.properties`.

See the [`./Makefile`](Makefile) for full details

### Database connections

Sample connection code with test pings can be found in [`./cmd/app/commands/debug`](./cmd/app/commands/debug).

All database packages at [`./internal/database`](./internal/database) expose the global interface:

- `Db(name ...string)`: Use this to retrieve a database connection. If `name` is specified, the named connection will be returned.
- `Init(opts database.ConnectionOpts, name ...string)`: Use this to initialise a database connection. If `name` is specified, a named connection will be created and added to a map of connections which you can retrieve using `Db(connectionName)`.
- `Close(name ...string)`: Use this to close connections you no longer need. If `name` is specified, a named connection will be closed and removed from the connection map. You need to call this before you `Init()` another connection with the same name.

Named connections can be used to implement logical tenant separations at the connection level. For example for B2B businesses where each tenant has a separate database instance, a new connection for `tenantA` can be initialised using `Init(opts, "tenantA")` and retrieved using `Db("tenantA")`.

### Building on top of the boilerplate

#### Changing the application name

Insert the following into `Makefile.properties`:

```Makefile
APP_NAME := "app2"
```

You will also need to:
1. Rename the module in `go.mod` to `app2` and run `go mod vendor`
2. Rename the chart in `./deploy/charts/app/Chart.yaml` to `app2`
3. Rename the directory `./deploy/charts/app` to `./deploy/charts/app2`
4. Rename the directory in `./cmd/app` to `./cmd/app2`

#### Changing the Docker image URL

Insert the following into `Makefile.properties`:

```Makefile
IMAGE_REGISTRY := "docker.registry.domain.com"
IMAGE_PATH := "your-org/your-app"
IMAGE_TAG := "main-123456"
```

This example will cause images to be tagged and pushed to `docker.registry.domain.com/your-org/your-app:main-123456`

#### Adding a new functional domain

Domains should be in their own internal package in `./internal/`. A package should contain both the controllers and the HTTP interfaces exposed with a reference performed in `./internal/api/`

An example with both controllers and HTTP interfaces can be found in `./internal/example` for an `example` domain.

### Pre-production checklist

This repository being a boilerplate contains some redunancies like database code which might not be used depending on your service. Before going to production:

- [ ] Helm chart is name appropriately
- [ ] Service is named appropriately
- [ ] Remove the `./cmd/app/commands/debug/*` directories that aren't used. Eg if MongoDB is not used, remove `./cmd/app/commands/debug/mongo` and remove the reference to it in `./cmd/app/commands/debug/debug.go`. This way, code for handling that database type will not be included in your binary.
- [ ] Disable Swagger in production (unless you want it in production)
- [ ] Generate fresh keys on a privileged system and delete any keys that come with the boilerplate

# Licensing

See [`./LICENSE`](./LICENSE). Basically, do whatever you want with this.

# PokemonGO API

A Go REST API to make a working API that fetch public API https://pokeapi.co/ to give RESTFUL response to React Vite frontend

## Prerequisites

- Go 1.21+, tested on `go1.25.5 X:nodwarf5 linux/amd64`
- PostgreSQL 15, tested on `docker compose postgres 15-alpine`
- Docker, tested on `29.0.4`, thus using `go1.25.4` by 18th December 2025
- PokeAPI, accessible on 18th December 2025

The API was tested on Fedora 43 but following requirement should be compatible with Windows and MacOS in their own ways.

## Installation

### 1. Clone Repository

`git clone https://github.com/deanAirre/pokeGOAPI`
`cd pokeGOAPI`

### 2. Install Go Dependencies

`go mod download`

### 3. Install Docker or Podman (For consistencies the readme will use Docker)

`sudo dnf install docker docker-compose`
or
`sudo apt install docker.io docker-compose`

### 4. Start Docker service to run PostgreSQL

```
# Start docker service
sudo systemctl start docker
sudo systemctl enable docker


# Add user to docker group
sudo usermod -a6 docker $USER
newgrp docker
```

_For Windows / Mac:_

- Install Docker Desktop from https://www.docker.com/products/docker-desktop/

---

## Quick Start

1. Start PostgreSQL Database
   `docker compose up -d`
2. Verify databse is running
   `docker ps`
3. Run the API Server
   `go run main.go`

Expected response:

```docker ps
CONTAINER ID   IMAGE                COMMAND                  CREATED         STATUS                            PORTS                                         NAMES
975066d55033   postgres:15-alpine   "docker-entrypoint.s…"   4 seconds ago   Up 3 seconds (health: starting)   0.0.0.0:5432->5432/tcp, [::]:5432->5432/tcp   pokemongo-postgres
```

````This is default password for test, change it later
2025/12/18 14:57:15 Config loaded
2025/12/18 14:57:15 Database connected
2025/12/18 14:57:15 Migrations completed
2025/12/18 14:57:15  Server starting on http://localhost:8080 ```

````

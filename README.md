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

_Test the API_

```curl http://localhost:8080/health

```

_Expected Response_

```
{"status":"healthy", "service":"pokemon-api"}

```

## Configuration

_Default Configuration_

## **Variable** **Default Value** **Description**

`DB_HOST` | localhost | Database host

`DB_POST` | 5432 | Database port

`DB_PASSWORD` | postgres Database | password

`DB_NAME` | pokemon_db | Database name

`SERVER_PORT` | 8080 | API server port

---

\*The default password are meant only for first installation, for later production it is recommended to change the password for better security.

--

## API Endpoints

---

## _Method_ _Endpoint_ _Description_

GET /health Health check endpoint

## Various commands

### View database logs

`docker logs pokemongo-postgres`

### View live logs

`docker logs -f pokemongo-postgres`

### Restart Database

`docker compose restart postgres`

### Stop Database

`docker compose down`

### Reset Database (Delete all data)

`docker compose down -v`
`docker compose up -d //turn it back on, clean slate`

### Access via psql

`docker exec -it pokemongo-postgres psql -U postgres -d pokemon_db`

## Docker Management (Linux)

### Start Docker

`sudo systemctl start docker`

### Stop Docker

`sudo systemctl stop docker`

### Restart Docker

`sudo systemctl restart docker`

### Check Docker status

`sudo systemctl status docker`

### Enable Docker on boot

`sudo systemctl enable docker`

### Disable Docker on boot

`sudo systemctl disable docker`

Docker uses plenty of resources and personally recommended to be turned off in times it is not required to run. If sudden spike in RAM or CPU uses every startup, better check your `docker compose`

---

## Go Development

### Run Application

`go run main.go`

### Build Binary

`go build -o bin/pokemongo main/main.go ./bin/pokemongo`

### Run Tests

`go test ./..`

## Troubleshooting

- "docker command not ofund"
  Install docker following the installation steps above

- "permission denied while trying to connect to Docker
  Add your user to the docker group (linux solution)
  `sudo usermod -aG docker $USER`
  `newgrp docker`

- "port 5432 is already in use"
  Another PostgreSQL instance probably occupying the port, check your port with `lsof` or `ss` or `netstat`

` sudo ss -tulpn | grep :5432`

If it's like `docker` or `postgresql` then it is another instance of either still running from past operations.

Turn it off with
`sudo systemctl stop postgresql`
or
`sudo systemctl stop docker`
or might as well clean the docker-compose and start it again if needed
`docker compose down`

If necessary, or the port is necessarily occupied, change the port in `docker-compose.yml`

- "container is restarting continously"
  Check logs first with
  `docker logs pokemongo-postgres`
  `docker logs -f pokemongo-postgres`

Common causes:

- Permission issues with mounted files (initiating new .db or something)
- Port conflicts
- Corrupted Docker volume

Solution:
`docker compose down -v`
`docker compose up -d`

"Failed to load config:DB_PASSWORD required"
Probable postgresql config issues, may want to check
`psql -U postgres`
or if it's already in Docker
`docker exec -it pokemongo-postgres psql -U postgres`

then change the password
`ALTER USER postgres PASSWORD 'new_password'`;

verify that it worked with ALTER ROLE or try logging in again

- "Database connection refuse"
  Ensure Docker and PostgreSQL container are running
  `sudo systemctl status docker`
  `docker ps`
  `docker logs pokemongo-postgres`

If it's on windows, nevermind the sudo part.
If it's not running, try restarting it.
If it's running but somehow still have that problem, well..
I'd recommend ask someone to debug it. In my experience it's postgresql not correctly installed, docker not correctly installed, version outdated, dependency not found, etc. Goodluck

## Contact

For questions or support, please open issue on GitHub

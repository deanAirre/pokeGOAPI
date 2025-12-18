# PokemonGO API

A Go REST API to make a working API that fetch public API https://pokeapi.co/ to give RESTFUL response to React Vite frontend

This API uses Go and PostgreSQL stack, to use this API, readers must at least preinstall things listed below suitable with user's OS is it Windows or Linux:

1. `postgresql-server` and `postgresql-contrib`, tested using `dbeaver 25.2` with:

- `postgresql-private-libs-18.1-1.fc43.x86_64`
- `postgresql-18.1-1.fc43.x86_64`
- `postgresql-contrib-18.1-1.fc43.x86_64`
- `postgresql-server-18.1-1.fc43.x86_64`

2. Golang with `go1.25.5` version or equal, tested using

- `go 1.25.5 X:nodwarf5 linux/amd64`

---

To run docker compose and get the .envs and setup postgreSQL database immediately
`docker compose up -d`

To run the API, use:
`go run main.go`

## Little Notes for Linux Users

For linux users:

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

or add these envs manually

```POSTGRES_USER: postgres
POSTGRES_PASSWORD: postgres
POSTGRES_DB: pokemon_db
```

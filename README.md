# OKDP Server

## Using Makefile

make run

## Using docker compose

### Prerequisites

Manually add the following entry in /etc/hosts

```shell
127.0.0.1       keycloak
```

### Start

```shell
docker-compose up
```

Open swagger UI at: http://localhost:8092/

### Rebuild

```shell
docker-compose build --no-cache 
```


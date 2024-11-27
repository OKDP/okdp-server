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

Authenticate with Swagger : oauth2 (OAuth2, authorizationCode with PKCE)

dev1/user

adm1/user

view1/user

Test the following endpoint:
```
http://localhost:8092/#/deployments/ListDeployments
spaceId: 1
compositionId:1
```

Then:
```
spaceId: 2
compositionId:1
```

### Rebuild

```shell
docker-compose build --no-cache 
```


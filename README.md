# OKDP Server

## Using Makefile

make run

## Using docker compose

### Prerequisites

Manually add the following entry in /etc/hosts

```shell
127.0.0.1       keycloak
```

Forward kad webserver to localhost:

```shell
kubectl port-forward svc/kad-webserver 6553:6553
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

### Rebuild

```shell
docker-compose build --no-cache 
```

```shell
curl -H "Authorization: Bearer qn7ccrJQhBJ9bdJ4sPa3LAXR8mrjsHen" -X PUT 'https://kad.ingress.kubo4.mbp/api/git/v1/mycluster/deployments/minio3.yaml' -F kadfile=@minio3.yaml -F commit-message='A commit Message' -F committer-name='Serge' -F committer-email=''

```


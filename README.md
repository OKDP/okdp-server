# OKDP Server

## Using Makefile
```shell
make help (or make)
make run
```

## Using docker compose

### Prerequisites

Manually add the following entry in /etc/hosts

```shell
127.0.0.1       keycloak
```

### Start

Start docker compose using your robot account token to access private registries:
```shell
docker-compose rm -f
OCI_USERNAME=okdp+okdp_quay_robot OCI_PASSWORD=****** docker-compose up --build
```

Open swagger UI at: http://localhost:8092/

Authenticate with Swagger : oauth2 (OAuth2, authorizationCode with PKCE)

dev1/user

adm1/user

view1/user

# Helm

```
docker build -t quay.io/okdp/okdp-server:0.1.0-snapshot  .
docker push quay.io/okdp/okdp-server:0.1.0-snapshot

helm package ./helm/okdp-server
helm push okdp-server-0.1.0-snapshot.tgz oci://quay.io/okdp/charts

helm pull oci://quay.io/okdp/charts/swagger-ui --version 0.1.0 --destination helm/okdp-server/charts

helm upgrade --install okdp-server \
     --namespace okdp-server \
     --create-namespace helm/okdp-server \
     --values helm/okdp-server/values.keycloak.yaml
```

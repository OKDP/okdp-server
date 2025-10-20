[![ci](https://github.com/okdp/okdp-server/actions/workflows/ci.yml/badge.svg)](https://github.com/okdp/okdp-server/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/OKDP/okdp-server)](https://github.com/OKDP/okdp-server/releases/latest)
[![Release](https://img.shields.io/github/v/release/OKDP/okdp-ui)](https://github.com/OKDP/okdp-ui/releases/latest)
[![KuboCD](https://img.shields.io/badge/kubocd-v0.2.1-green.svg)](https://github.com/kubocd/kubocd)
[![License Apache2](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)
<a href="https://okdp.io">
<img src="https://okdp.io/logos/okdp-notext.svg" height="20px" style="margin: 0 2px;" />
</a>


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
docker-compose up --build
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
helm push okdp-server-0.2.0-snapshot.tgz oci://quay.io/okdp/charts

helm pull oci://quay.io/okdp/charts/swagger-ui --version 0.2.0 --destination helm/okdp-server/charts

helm upgrade --install okdp-server \
     --namespace okdp-server \
     --create-namespace helm/okdp-server \
     --values helm/okdp-server/values.keycloak.yaml
```

Swagger UI: https://okdp-server.okdp.sandbox/swagger/


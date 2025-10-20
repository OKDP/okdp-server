[![ci](https://github.com/okdp/okdp-server/actions/workflows/ci.yml/badge.svg)](https://github.com/okdp/okdp-server/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/okdp/okdp-server)](https://github.com/OKDP/okdp-server/releases/latest)
[![okdp-ui](https://img.shields.io/badge/okdp--ui-v0.4.1-%2341928F?style=flat)](https://github.com/OKDP/okdp-ui/releases/latest)
[![KuboCD](https://img.shields.io/badge/kubocd-v0.2.1-green.svg)](https://github.com/kubocd/kubocd)
[![License Apache2](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)
<a href="https://okdp.io">
<img src="https://okdp.io/logos/okdp-notext.svg" height="20px" style="margin: 0 2px;" />
</a>

OKDP Server is the backend API for the OKDP Platform , providing a unified REST API for managing deployments, clusters, projects, GitOps repositories, and package catalogs.
It powers the [OKDP Control Plane UI](https://github.com/OKDP/okdp-ui), exposing a consistent API to interact with Kubernetes and Git-based deployment environments.

# Test the Rest API

The easiest way to start testing the **OKDP Rest API** is by using the [okdp sandbox](https://github.com/OKDP/okdp-sandbox):

The sandbox provides a local, preconfigured OKDP environment — including the necessary front end services, dependencies, and sample data — so you can quickly validate and interact with the Rest API without setting up the entire platform manually.


# Developing/Testing locally

This project is configured with a **Dev Container** for a consistent development environment.  

## Using Makefile

```shell
make help (or make)
make run
```

## Using docker compose

1. Manually add the following entry in /etc/hosts

```shell
127.0.0.1       keycloak
```

2. Start docker compose using your robot account token to access private registries:
```shell
docker-compose rm -f
docker-compose up --build
```

3. Open swagger UI at: http://localhost:8092/

Authenticate with Swagger using OAuth2/authorizationCode with PKCE

| Username | Password |
|-----------|---------|
| dev1      | user    |
| adm1      | user    |
| view1     | user    |


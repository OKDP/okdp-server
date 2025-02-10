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

Forward kad webserver to localhost:

```shell
kubectl port-forward svc/kad-webserver 6553:6553
```

### Start

```shell
docker-compose rm -f
docker-compose up --build
```

Open swagger UI at: http://localhost:8092/

Authenticate with Swagger : oauth2 (OAuth2, authorizationCode with PKCE)

dev1/user

adm1/user

view1/user

```shell

kadcli git remove projectA-1.yaml -m "Test remove" -n "idir" --insecureSkipVerify
kadcli git add -m "test add" -n"idir"  projectA-1.yaml --insecureSkipVerify 

curl -H "Authorization: Bearer HEDG296X4XjnjETBJ1HGEUEqQbn3pNaD" -X PUT 'https://kad.ingress.kind.local/api/git/v1/mycluster/deployments/minio3.yaml' -F kadfile=@.tmp/z.tmp/curl/minio3.yaml -F commit-message='A commit Message' -F committer-name='Serge' -F committer-email='serge@example.com' -k


kadcli kad componentReleases apply --insecureSkipVerify minio1
kadcli kad componentReleases apply --insecureSkipVerify _all_ 

```


http://localhost:8092/#/componentreleases/CreateOrUpdateComponentRelease
```
{
  "comment": "Create minio deployment example",
  "gitRepoFolder": "deployments",
  "componentReleases": [
    {
      "name": "minio3",
      "component": {
        "name": "minio",
        "version": "1.0.0",
        "protected": true,
        "config": {
          "install": {
            "createNamespace": true
          }
        },
        "parameters": {
          "ingressName": "minio3",
          "ldap": "openldap"
        },
        "parameterFiles": [
          {
            "document": "minio-flavor-small"
          },
          {
            "document": "data1-minio-parameters",
            "unwrap": "minio"
          }
        ]
      },
      "namespace": "minio3",
      "dependsOn": [
        "ldapLocalServer"
      ]
    }
  ]
}
```


# Helm

```
docker build -t quay.io/okdp/okdp-server:0.1.0-snapshot  .
docker push quay.io/okdp/okdp-server:0.1.0-snapshot 
helm upgrade --install okdp-server \
     --namespace okdp-server \
     --create-namespace helm/okdp-server \
     --values helm/okdp-server/values.keycloak.yaml
```

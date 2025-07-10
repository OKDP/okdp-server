# okdp-server

![Version: 0.3.0](https://img.shields.io/badge/Version-0.3.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.3.0](https://img.shields.io/badge/AppVersion-0.3.0-informational?style=flat-square)

A Helm chart for okdp-server

**Homepage:** <https://okdp.io>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| idirze | <idir.izitounene@kubotal.io> | <https://github.com/idirze> |

## Source Code

* <https://github.com/OKDP/okdp-ui>

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| oci://quay.io/okdp/charts | swagger-ui | 0.1.0 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | Affinity for pod scheduling. |
| autoscaling.enabled | bool | `false` |  |
| autoscaling.maxReplicas | int | `100` |  |
| autoscaling.minReplicas | int | `1` |  |
| autoscaling.targetCPUUtilizationPercentage | int | `80` |  |
| configuration.catalog | list | `[{"description":"My Storage packages","id":"storage","name":"Storage catalog","packages":[{"name":"redis"},{"name":"minio"},{"name":"cnpg"}],"repoUrl":"quay.io/kubocd/packages"},{"description":"My Auth packages","id":"auth","name":"Auth catalog","packages":[{"name":"openldap"}],"repoUrl":"quay.io/kubocd/packages"},{"description":"My Infra packages","id":"infra","name":"Infra catalog","packages":[{"name":"cert-manager"},{"name":"ingress-nginx"},{"name":"metallb"}],"repoUrl":"quay.io/kubocd/packages"},{"description":"My Stack packages","id":"stack","name":"Stack catalog","packages":[{"name":"podinfo"}],"repoUrl":"quay.io/kubocd/packages"}]` | List of catalogs available to this chart |
| configuration.catalog[0] | object | `{"description":"My Storage packages","id":"storage","name":"Storage catalog","packages":[{"name":"redis"},{"name":"minio"},{"name":"cnpg"}],"repoUrl":"quay.io/kubocd/packages"}` | Unique identifier for the catalog |
| configuration.catalog[0].description | string | `"My Storage packages"` | Description of the catalog's purpose |
| configuration.catalog[0].name | string | `"Storage catalog"` | Human-readable name of the catalog |
| configuration.catalog[0].packages | list | `[{"name":"redis"},{"name":"minio"},{"name":"cnpg"}]` | List of packages under this catalog |
| configuration.catalog[0].packages[0] | object | `{"name":"redis"}` | Name of the package |
| configuration.catalog[0].repoUrl | string | `"quay.io/kubocd/packages"` | OCI registry URL to pull packages from |
| configuration.clusters | list | `[{"auth":{"inCluster":true},"env":"dev","id":"kubo2","name":"My k8s cluster 1"}]` | List of Kubernetes clusters this chart will interact with |
| configuration.clusters[0] | object | `{"auth":{"inCluster":true},"env":"dev","id":"kubo2","name":"My k8s cluster 1"}` | Unique identifier for the cluster |
| configuration.clusters[0].auth.inCluster | bool | `true` | Use in-cluster authentication |
| configuration.clusters[0].env | string | `"dev"` | Environment tag (e.g., dev, staging, prod) |
| configuration.clusters[0].name | string | `"My k8s cluster 1"` | Human-readable name for the cluster |
| configuration.logging.format | string | `"console"` | Specify the logging format. One of `console` or `json`. |
| configuration.logging.level | string | `"debug"` | Specify the logging level. One of `debug`, `info`, `warn`, `error`, `fatal` or `panic`. |
| configuration.security.authN.bearer.groupsAttributePath | string | `"realm_access.groups"` | Specify the groups attribute path from json access token. |
| configuration.security.authN.bearer.issuerUri | string | `""` | Specify the issuer uri. |
| configuration.security.authN.bearer.jwksURL | string | `""` | Specify the jwks URL. |
| configuration.security.authN.bearer.rolesAttributePath | string | `"realm_access.roles"` | Specify the roles attribute path from json access token. |
| configuration.security.authN.bearer.skipIssuerCheck | bool | `false` | Wether to skip issuer check. |
| configuration.security.authN.bearer.skipSignatureCheck | bool | `false` | Wether to skip issuer signature check. |
| configuration.security.authN.provider | list | `["bearer"]` | Specify the oidc privider. One of `openid` or `bearer`. |
| configuration.security.authZ.inline | object | `{"model":"[request_definition]\nr = sub, obj, act\n\n[policy_definition]\np = sub, obj, act\n\n[role_definition]\ng = _, _\n\n[policy_effect]\ne = some(where (p.eft == allow))\n\n[matchers]\nm = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && (r.act == p.act || p.act == \"*\")\n","policy":"p, role:viewers, /api/v1/users/myprofile, *\np, role:viewers, /api/v1/catalogs, *\np, role:viewers, /api/v1/catalogs/*, *\n\np, role:viewers, /api/v1/clusters, *\np, role:viewers, /api/v1/clusters/*/gitrepos, *\np, role:viewers, /api/v1/clusters/*/gitrepos/*, *\n\ng, role:admins, role:developers\ng, role:developers, role:viewers\n"}` | More info: https://casbin.org/docs/how-it-works/ file:   modelPath: ".local/authz-model.conf"   policyPath: ".local/authz-policy.csv" |
| configuration.security.authZ.inline.model | string | `"[request_definition]\nr = sub, obj, act\n\n[policy_definition]\np = sub, obj, act\n\n[role_definition]\ng = _, _\n\n[policy_effect]\ne = some(where (p.eft == allow))\n\n[matchers]\nm = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && (r.act == p.act || p.act == \"*\")\n"` | More info: https://casbin.org/docs/how-it-works/ |
| configuration.security.authZ.provider | string | `"inline"` | Specify the authZ storage provider. One of `inline` or `file`. |
| configuration.security.cors.allowCredentials | bool | `true` | Determine whether cookies and authentication credentials should be included in cross-origin requests. |
| configuration.security.cors.allowedHeaders | list | `["Origin","Accept","Authorization","Content-Length","Content-Type"]` | List the headers that clients are allowed to include in requests. |
| configuration.security.cors.allowedMethods | list | `["GET","POST","PUT","DELETE","PATCH","OPTIONS","HEAD"]` | Define the HTTP methods permitted for CORS requests. |
| configuration.security.cors.allowedOrigins | list | `["*"]` | Specify the allowed origins for cross-origin requests. "*" allows all origins. |
| configuration.security.cors.exposedHeaders | list | `["Content-Length"]` | Specify which response headers should be exposed to the client. |
| configuration.security.cors.maxAge | int | `3600` | Define how long (in seconds) the results of a preflight request can be cached by the client. |
| configuration.security.headers.X-Content-Type-Options | string | `"nosniff"` | Prevent browsers from MIME-sniffing a response away from the declared content type. |
| configuration.security.headers.X-Frame-Options | string | `"DENY"` | Prevent the page from being embedded in an iframe, mitigating clickjacking attacks. |
| configuration.server.listenAddress | string | `"0.0.0.0"` | Specify the Server listen address. |
| configuration.server.mode | string | `"debug"` | Specify the Server Mode. One of `debug`, `release` or `test`. |
| configuration.server.port | int | `8090` | Specify the Server listen port. |
| configuration.swagger.securitySchemes.oauth2.flows.authorizationCode.authorizationUrl | string | `nil` |  |
| configuration.swagger.securitySchemes.oauth2.flows.authorizationCode.scopes.email | string | `"User Email"` |  |
| configuration.swagger.securitySchemes.oauth2.flows.authorizationCode.scopes.openid | string | `"OpenId Authentication"` |  |
| configuration.swagger.securitySchemes.oauth2.flows.authorizationCode.scopes.profile | string | `"User Profile"` |  |
| configuration.swagger.securitySchemes.oauth2.flows.authorizationCode.scopes.roles | string | `"User Roles"` |  |
| configuration.swagger.securitySchemes.oauth2.flows.authorizationCode.tokenUrl | string | `nil` |  |
| configuration.swagger.securitySchemes.oauth2.type | string | `"oauth2"` |  |
| configuration.swagger.security[0].oauth2[0] | string | `"openid"` |  |
| configuration.swagger.security[0].oauth2[1] | string | `"email"` |  |
| configuration.swagger.security[0].oauth2[2] | string | `"profile"` |  |
| configuration.swagger.security[0].oauth2[3] | string | `"roles"` |  |
| fullnameOverride | string | `""` | Overrides the release name. |
| image.pullPolicy | string | `"Always"` | Image pull policy. |
| image.repository | string | `"quay.io/okdp/okdp-server"` | Docker image registry. |
| image.tag | string | `"0.3.0"` | Image tag. |
| imagePullSecrets | list | `[]` | Secrets to be used for pulling images from private Docker registries. |
| ingress.annotations | object | `{}` |  |
| ingress.className | string | `""` | Specify the ingress class (Kubernetes >= 1.18). |
| ingress.enabled | bool | `false` |  |
| ingress.hosts[0].host | string | `"chart-example.local"` |  |
| ingress.hosts[0].paths[0].path | string | `"/"` |  |
| ingress.hosts[0].paths[0].pathType | string | `"ImplementationSpecific"` |  |
| ingress.tls | list | `[]` |  |
| livenessProbe | object | `{"httpGet":{"path":"/healthz","port":"http"},"initialDelaySeconds":60,"periodSeconds":30,"timeoutSeconds":10}` | Liveness probe for the okdp-server container. |
| nameOverride | string | `""` | Override for the `okdp-server.fullname` template, maintains the release name. |
| nodeSelector | object | `{}` | Node selector for pod scheduling. |
| podAnnotations | object | `{}` | Additional annotations for the okdp-server pod. |
| podLabels | object | `{}` | Additional labels for the okdp-server pod. |
| podSecurityContext | object | `{}` |  |
| rbac.annotations | object | `{}` | Specify annotations for the proxy. |
| rbac.create | bool | `true` | Specify whether a RBAC should be created |
| readinessProbe | object | `{"httpGet":{"path":"/readiness","port":"http"}}` | Readiness probe for the okdp-server container. |
| replicaCount | int | `1` | Desired number of okdp-server pods to run. |
| resources | object | `{}` |  |
| securityContext | object | `{}` | Security context for the container. |
| service.port | int | `8090` |  |
| service.type | string | `"ClusterIP"` |  |
| serviceAccount.annotations | object | `{}` | Annotations to add to the service account |
| serviceAccount.automount | bool | `true` | Automatically mount a ServiceAccount's API credentials? |
| serviceAccount.create | bool | `true` | Specify whether a service account should be created |
| serviceAccount.name | string | `""` | If not set and create is true, a name is generated using the fullname template |
| swagger-ui.configuration.extraEnv | list | `[{"name":"BASE_URL","value":"/"},{"name":"URLS","value":null},{"name":"OAUTH_CLIENT_ID","value":null},{"name":"OAUTH_SCOPES","value":"openid profile email roles"},{"name":"OAUTH_USE_PKCE","value":true},{"name":"CORS_ENABLED","value":true}]` | Specify swagger configuration with environment variables (https://swagger.io/docs/open-source-tools/swagger-ui/usage/oauth2/) |
| swagger-ui.configuration.extraEnv[0] | object | `{"name":"BASE_URL","value":"/"}` | Specify base url. |
| swagger-ui.configuration.extraEnv[1] | object | `{"name":"URLS","value":null}` | Specify the list of api servers swagger docs to serve. |
| swagger-ui.configuration.extraEnv[2] | object | `{"name":"OAUTH_CLIENT_ID","value":null}` | Specify the Oauth2 client Id. The client ID must be be public for production. |
| swagger-ui.configuration.extraEnv[3] | object | `{"name":"OAUTH_SCOPES","value":"openid profile email roles"}` | Specify the Oauth2 scopes. |
| swagger-ui.enabled | bool | `true` |  |
| tolerations | list | `[]` | Tolerations for pod scheduling. |
| volumeMounts | list | `[]` | Additional volumeMounts on the output Deployment definition. |
| volumes | list | `[]` | Additional volumes on the output Deployment definition. |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.13.1](https://github.com/norwoodj/helm-docs/releases/v1.13.1)

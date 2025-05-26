#!/bin/bash

CONFIDENTIAL_CLIENT='confidential-oidc-client'
PUBLIC_CLIENT='public-oidc-client'
WEB_ORIGINS='["*"]'
REDIRECT_URIS='[
  "http://localhost:8090/oauth2/callback",
  "http://localhost:8092/oauth2-redirect.html",
  "http://localhost:4200/index.html",
  "http://localhost:4200/silent-refresh.html",
  "http://localhost:4200/home",
  "http://okdp-ui.okdp.sandbox/index.html",
  "https://okdp-ui.okdp.sandbox/index.html",
  "http://okdp-server.okdp.sandbox/swagger/oauth2-redirect.html",
  "https://okdp-server.okdp.sandbox/swagger/oauth2-redirect.html"
]'

get_client_id() {
  local client_name=$1
  /opt/keycloak/bin/kcadm.sh get clients -r master --fields id,clientId \
    | grep -B1 "\"clientId\" : \"${client_name}\"" \
    | grep '"id"' \
    | sed -E 's/.*"id" : "([^"]+)".*/\1/'
}

echo "Creating users, roles and clients ..."

# Connect to kecloak
/opt/keycloak/bin/kcadm.sh config credentials --server http://keycloak:$KC_HOSTNAME_PORT \
    --realm master --user $KC_BOOTSTRAP_ADMIN_USERNAME --password $KC_BOOTSTRAP_ADMIN_PASSWORD

# Create Users
/opt/keycloak/bin/kcadm.sh create users -r master -s username=dev1 -s firstName=dev1 -s lastName=dev1 -s enabled=true \
    -s email=dev1.developers@example.org -s emailVerified=true
/opt/keycloak/bin/kcadm.sh set-password -r master --username dev1 --new-password user

/opt/keycloak/bin/kcadm.sh create users -r master -s username=adm1 -s firstName=adm1 -s lastName=adm1 -s enabled=true \
    -s email=adm1.admins@example.org -s emailVerified=true
/opt/keycloak/bin/kcadm.sh set-password -r master --username adm1 --new-password user

/opt/keycloak/bin/kcadm.sh create users -r master -s username=view1 -s firstName=view1 -s lastName=view1 -s enabled=true \
    -s email=view1.viewers@example.org -s emailVerified=true
/opt/keycloak/bin/kcadm.sh set-password -r master --username view1 --new-password user

# Create Groups
/opt/keycloak/bin/kcadm.sh create groups -r master -s name=developers
/opt/keycloak/bin/kcadm.sh create groups -r master -s name=viewers
/opt/keycloak/bin/kcadm.sh create groups -r master -s name=admins

# Create roles
/opt/keycloak/bin/kcadm.sh create roles -r master -s name=developers
/opt/keycloak/bin/kcadm.sh create roles -r master -s name=viewers
/opt/keycloak/bin/kcadm.sh create roles -r master -s name=admins

# Assign Roles to users
/opt/keycloak/bin/kcadm.sh add-roles -r master --uusername dev1 --rolename developers
/opt/keycloak/bin/kcadm.sh add-roles -r master --uusername view1 --rolename viewers
/opt/keycloak/bin/kcadm.sh add-roles -r master --uusername adm1 --rolename admins

# Create OAuth2 clients
/opt/keycloak/bin/kcadm.sh create clients -r master -s clientId=$PUBLIC_CLIENT -s name=$PUBLIC_CLIENT -s publicClient=true \
                           -s "redirectUris=${REDIRECT_URIS}" \
                           -s "webOrigins=${WEB_ORIGINS}"
/opt/keycloak/bin/kcadm.sh create clients -r master -s clientId=$CONFIDENTIAL_CLIENT -s name=$CONFIDENTIAL_CLIENT  -s 'secret=secret1' \
                           -s "redirectUris=${REDIRECT_URIS}" \
                           -s "webOrigins=${WEB_ORIGINS}"

CONF_CLIENT_ID=$(get_client_id "$CONFIDENTIAL_CLIENT")
/opt/keycloak/bin/kcadm.sh update clients/$CONF_CLIENT_ID -r master \
  -s "redirectUris=${REDIRECT_URIS}" \
  -s "webOrigins=${WEB_ORIGINS}"

PUB_CLIENT_ID=$(get_client_id "$PUBLIC_CLIENT")
/opt/keycloak/bin/kcadm.sh update clients/$PUB_CLIENT_ID -r master \
  -s "redirectUris=${REDIRECT_URIS}" \
  -s "webOrigins=${WEB_ORIGINS}"

# Update access token lifetime
echo "Update access token lifetime to 8H"
/opt/keycloak/bin/kcadm.sh update realms/master -s accessTokenLifespan=28800
echo "Users, roles and clients created successfuly"
exit 0


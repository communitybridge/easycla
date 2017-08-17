#!/usr/bin/env bash

echo "export const CINCO_API_URL: string = \"${CINCO_SERVER_URL}\";" > /srv/app/src/app/src/services/constants.ts
echo "Wrote /srv/app/src/app/src/services/constants.ts"

echo "{" > /srv/app/src/app/src/assets/keycloak.json
echo "  \"realm\": \"LinuxFoundation\"," >> /srv/app/src/app/src/assets/keycloak.json
echo "  \"auth-server-url\": \"${KEYCLOAK_SERVER_URL}\"," >> /srv/app/src/app/src/assets/keycloak.json
echo "  \"ssl-required\": \"external\"," >> /srv/app/src/app/src/assets/keycloak.json
echo "  \"resource\": \"pmc\"," >> /srv/app/src/app/src/assets/keycloak.json
echo "  \"public-client\": true" >> /srv/app/src/app/src/assets/keycloak.json
echo "}" >> /srv/app/src/app/src/assets/keycloak.json
echo "Wrote /srv/app/src/app/src/assets/keycloak.json"
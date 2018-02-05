#!/usr/bin/env bash

echo "" > /srv/app/src/ionic/services/constants.ts
echo "export const CINCO_API_URL: string = \"${CINCO_SERVER_URL}\";" >> /srv/app/src/ionic/services/constants.ts
echo "export const CLA_API_URL: string = \"${CLA_SERVER_URL}\";" >> /srv/app/src/ionic/services/constants.ts
echo "Wrote /srv/app/src/ionic/services/constants.ts"

echo "{" > /srv/app/src/ionic/assets/keycloak.json
echo "  \"realm\": \"LinuxFoundation\"," >> /srv/app/src/ionic/assets/keycloak.json
echo "  \"auth-server-url\": \"${KEYCLOAK_SERVER_URL}/auth\"," >> /srv/app/src/ionic/assets/keycloak.json
echo "  \"ssl-required\": \"external\"," >> /srv/app/src/ionic/assets/keycloak.json
echo "  \"resource\": \"cla-management-console\"," >> /srv/app/src/ionic/assets/keycloak.json
echo "  \"public-client\": true" >> /srv/app/src/ionic/assets/keycloak.json
echo "}" >> /srv/app/src/ionic/assets/keycloak.json
echo "Wrote /srv/app/src/ionic/assets/keycloak.json"

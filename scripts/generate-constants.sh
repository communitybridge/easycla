#!/usr/bin/env bash

cat > /srv/app/src/app/src/services/constants.ts <<'endmsg'
export const CINCO_API_URL: string = "${CINCO_SERVER_URL}";
endmsg
echo "Wrote /srv/app/src/app/src/services/constants.ts"

cat > /srv/app/src/app/src/assets/keycloak.json <<'endmsg'
{
  "realm": "LinuxFoundation",
  "auth-server-url": "${KEYCLOAK_SERVER_URL}",
  "ssl-required": "external",
  "resource": "pmc",
  "public-client": true
}
endmsg
echo "Wrote /srv/app/src/app/src/assets/keycloak.json"
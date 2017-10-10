#!/usr/bin/env bash

echo "" > /srv/app/src/ionic/services/constants.ts
echo "export const CINCO_API_URL: string = \"${CINCO_SERVER_URL}\";" >> /srv/app/src/ionic/services/constants.ts
echo "export const CLA_API_URL: string = \"${CLA_SERVER_URL}\";" >> /srv/app/src/ionic/services/constants.ts
echo "Wrote /srv/app/src/ionic/services/constants.ts"

#!/usr/bin/env bash

echo "export const CINCO_API_URL: string = \"${CINCO_SERVER_URL}\";" > /srv/app/src/ionic/services/constants.ts
echo "Wrote /srv/app/src/ionic/services/constants.ts"
#!/usr/bin/env bash

/srv/app/scripts/generate-constants.sh

cd /srv/app/src

npm install; npm run build

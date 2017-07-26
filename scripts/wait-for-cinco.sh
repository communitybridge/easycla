#!/usr/bin/env bash

echo "Waiting for CINCO to become available..."

until $(curl --output /dev/null --silent --head --fail ${CINCO_SERVER_URL}about/version); do
    printf '.'
    sleep 5
done
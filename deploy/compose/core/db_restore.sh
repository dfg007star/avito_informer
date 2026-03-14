#!/usr/bin/env bash
set -e

cp ./.env ./app/deploy/compose/core

cd ./app/deploy/compose/core
docker-compose up -d postgres-avito-informer

echo "prepare to restore database!"
sleep 10
docker exec -i postgres-avito-informer pg_restore -U avito-informer-service-user -d avito-informer-service -v <../../../../backup.dump

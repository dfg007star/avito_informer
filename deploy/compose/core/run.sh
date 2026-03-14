#!/usr/bin/env bash
set -e

cp ./.env ./app/deploy/compose/core

echo "task deps:update"
cd ./app
if [ ! -f "go.work" ]; then
        go work init
        for module in ./collector ./http ./notification ./platform; do
            go work use $module
        done
        go work sync
fi
task deps:update
task rebuild

#!/bin/bash
set -e
bunx openapi-typescript \
  http://0.0.0.0:9788/openapi.json \
  -o ../client/src/lib/api/schema.d.ts

#!/bin/bash
set -e
bunx openapi-typescript \
  ../openapi.json \
  -o ../client/src/lib/api/schema.d.ts

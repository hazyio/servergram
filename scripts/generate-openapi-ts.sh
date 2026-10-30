#!/bin/bash
set -e
bunx openapi-typescript \
  ../openapi.json \
  -o ../client/src/lib/types/request-schema.d.ts

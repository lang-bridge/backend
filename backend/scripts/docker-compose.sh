#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

#GOLANG_IMAGE=${GOLANG_IMAGE:-"golang:1.22.2-bullseye"}
#APP_BASE_IMAGE=${APP_BASE_IMAGE:-"alpine:3.18.6"}
APP_NAME="langbridge"
COMPOSE_DIR="deployments/local/"
COMPOSE_ARGS=( -p "$APP_NAME" )

# Postgres
COMPOSE_POSTGRES_ENABLE=${COMPOSE_POSTGRES_ENABLE:-""}
if [[ -n "$COMPOSE_POSTGRES_ENABLE" ]]; then
  COMPOSE_ARGS+=( -f "${COMPOSE_DIR}postgres.docker-compose.yml" )
  COMPOSE_ARGS+=( -f "${COMPOSE_DIR}migrator.docker-compose.yml" )

  export POSTGRES_IMAGE=${POSTGRES_IMAGE:-"postgres:16.3-bullseye"}
  export POSTGRES_LOG_STATEMENT=${POSTGRES_LOG_STATEMENT:-"none"}
fi

# ORY
COMPOSE_ORY_ENABLE=${COMPOSE_ORY_ENABLE:-""}
if [[ -n "$COMPOSE_ORY_ENABLE" ]]; then
  COMPOSE_ARGS+=( -f "${COMPOSE_DIR}ory/kratos.docker-compose.yml" )
  export ORY_KRATOS_IMAGE=${ORY_KRATOS_IMAGE:-"oryd/kratos:v1.2.0"}

  COMPOSE_ARGS+=( -f "${COMPOSE_DIR}ory/oathkeeper.docker-compose.yml" )
  export ORY_OATHKEEPER_IMAGE=${ORY_OATHKEEPER_IMAGE:-"oryd/oathkeeper:v0.40.7"}
fi

# API
COMPOSE_API_ENABLE=${COMPOSE_API_ENABLE:-""}
if [[ -n "$COMPOSE_API_ENABLE" ]]; then
  COMPOSE_ARGS+=( -f "${COMPOSE_DIR}api.docker-compose.yml" )
fi

echo "docker-compose ${COMPOSE_ARGS[@]} up"

docker-compose "${COMPOSE_ARGS[@]}" down --remove-orphans
docker-compose "${COMPOSE_ARGS[@]}" up --remove-orphans --build
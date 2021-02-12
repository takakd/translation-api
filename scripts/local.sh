#!/usr/bin/env bash

#
# Helper scripts for developing
#

SCRIPT_DIR=$(cd $(dirname "$0"); pwd)
ARGV=("$@")
ARGC=$#

function usage() {
cat <<_EOT_
Usage:
  $0 Command

Example.
  $0 build

Command:
  run           Run envoy and gRPC server on local.
  run:go        Run gRPC server on local.
  run:envoy     Run envoy on local.
  down:envoy    Stop envoy on local.
  grpc          Generate gRPC codes.
_EOT_
exit 1
}

run() {
    proxy_envoy_run

    # Call if it's entered Ctrl+C
    trap proxy_envoy_down SIGINT

    run_go

    proxy_envoy_down
}

run_go() {
    echo Run go cmd.
    cd "${SCRIPT_DIR}/../cmd/api" || exit
    APP_ENV=local ENV_FILE="${SCRIPT_DIR}/../deployments/local/.env" go run ./api.go
}

proxy_envoy_run() {
    docker-compose -f ${SCRIPT_DIR}/../deployments/local/docker-compose.yml up -d
    # wait for starting
    sleep 5
}

proxy_envoy_down() {
    docker-compose -f ${SCRIPT_DIR}/../deployments/local/docker-compose.yml down
}

grpc() {
    # Ref.
    #   https://grpc.io/docs/languages/go/quickstart/#regenerate-grpc-code
    #   https://developers.google.com/protocol-buffers/docs/reference/go-generated#invocation

    # Generate Go codes
    #   translator
    mkdir -p ${SCRIPT_DIR}/../internal/app/grpc/translator
    protoc \
        --proto_path=${SCRIPT_DIR}/../api/grpc/translator \
        --go-grpc_out=${SCRIPT_DIR}/../internal/app/grpc/translator \
        --go-grpc_opt=paths=source_relative \
        --go_out=${SCRIPT_DIR}/../internal/app/grpc/translator \
        --go_opt=paths=source_relative \
        ${SCRIPT_DIR}/../api/grpc/translator/translator.proto
    "${SCRIPT_DIR}/mock.pl" "${SCRIPT_DIR}/../internal/app/grpc/translator/translator_grpc.pb.go"
}


if [[ $# -lt 1 ]]; then
    usage
fi

if [[ $1 = "run" ]]; then
    run
elif [[ $1 = "run" ]]; then
    proxy_envoy_run
    run_go
elif [[ $1 = "run:go" ]]; then
    run_go
elif [[ $1 = "run:envoy" ]]; then
    proxy_envoy_run
elif [[ $1 = "down:envoy" ]]; then
    proxy_envoy_down
elif [[ $1 = "grpc" ]]; then
    grpc
else
    usage
fi

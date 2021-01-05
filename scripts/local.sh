#!/usr/bin/env bash

#
# Helper scripts for Makefile
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
  build         Build app binary.
  fmt           Format sources.
  run           Run on local.
  test          Run test on local.
  install       Install dependency modules
  docker:run    Run docker on local.
  docker:stop   Stop docker on local.
_EOT_
exit 1
}

build() {
    cd "${SCRIPT_DIR}/../cmd/api" || exit
    go build -p 2 -v -x api.go
}

fmt() {
    go fmt ./...
    # Ref: https://gist.github.com/bgentry/fd1ffef7dbde01857f66#gistcomment-1618537
    goimports -w $(find . -type f -name "*.go" -not -path "./vendor/*")
    golint ./cmd/... ./internal/...
    go vet ./cmd/... ./internal/...
}

run() {
    proxy_envoy_run

    # Call if it's entered Ctrl+C
    trap proxy_envoy_down SIGINT

    echo Run go cmd.
    cd "${SCRIPT_DIR}/../cmd/api" || exit
    DOT_ENV=local ./api

    proxy_envoy_down
}
proxy_envoy_run() {
    docker-compose -f ${SCRIPT_DIR}/../deployments/envoy/docker-compose.yml up -d
    # wait for starting
    sleep 5
}
proxy_envoy_down() {
    docker-compose -f ${SCRIPT_DIR}/../deployments/envoy/docker-compose.yml down
}

cmd_test() {
    docker_run

    cd ${SCRIPT_DIR}/..

    # @see https://stackoverflow.com/questions/16353016/how-to-go-test-all-tests-in-my-project/35852900#35852900
    # NG
    #go test -v -cover "${ARGS}" ./...
    # OK
    sh -c "go $(echo ${ARGV[@]})"

    docker_cleanup
}

grpc() {
    protoc --go-grpc_out=${SCRIPT_DIR}/../internal/app/grpc/translator --go-grpc_opt=paths=source_relative --proto_path=${SCRIPT_DIR}/../internal/app/grpc/translator --go_out=${SCRIPT_DIR}/../internal/app/grpc/translator --go_opt=paths=source_relative ${SCRIPT_DIR}/../internal/app/grpc/translator/translator.proto
}

install() {
    go env -w GO111MODULE=on
    go mod vendor -v
}

if [[ $# -lt 1 ]]; then
    usage
fi

if [[ $1 = "build" ]]; then
    build
elif [[ $1 = "fmt" ]]; then
    fmt
elif [[ $1 = "run" ]]; then
    run
elif [[ $1 = "test" ]]; then
    cmd_test
elif [[ $1 = "install" ]]; then
    install
elif [[ $1 = "docker:run" ]]; then
    docker_run
elif [[ $1 = "docker:stop" ]]; then
    docker_cleanup
elif [[ $1 = "grpc" ]]; then
    grpc
else
    usage
fi

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
  $0 command

Example.
  $0 build

command:
  build         Build app binary.
  fmt           Format sources.
  test          Run test on local.
_EOT_
exit 0
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

cmd_test() {
    cd ${SCRIPT_DIR}/..

    # @see https://stackoverflow.com/questions/16353016/how-to-go-test-all-tests-in-my-project/35852900#35852900
    # NG
    #go test -v -cover "${ARGS}" ./...
    # OK
    sh -c "go $(echo ${ARGV[@]})"
}


if [[ $# -lt 1 ]]; then
    usage
fi

if [[ $1 = "build" ]]; then
    build
elif [[ $1 = "fmt" ]]; then
    fmt
elif [[ $1 = "test" ]]; then
    cmd_test
else
    usage
fi

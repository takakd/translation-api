#!/usr/bin/env bash

#
# Build docker image script
#

SCRIPT_DIR=$(cd $(dirname "$0"); pwd)
ARGV=("$@")
ARGC=$#

function usage() {
cat <<_EOT_
Usage:
  $0 [image tag]

Example.
  $0 com.example.mydockerimage:latest
_EOT_
exit 1
}


# Check parameters
if [[ $# -lt 1 || "$1" == "" ]]; then
    usage
fi


# Start to build image

DOCKER_DIR=${SCRIPT_DIR}/../deployments/docker-image

#   API
DOCKER_API_SRC_DIR=${DOCKER_DIR}/api/src

#   Clean up
rm -rf "$DOCKER_API_SRC_DIR"

#   Copy codes to Docker working directory
for file in $(git ls-files | grep -E "(\\.go|Makefile|scripts|go.mod)"); do
    REL_DIR=$(dirname "$file")
    DIR="${DOCKER_API_SRC_DIR}/${REL_DIR}"
    if [[ ! -e "$DIR" ]]; then
        mkdir -p "$DIR"
    fi
    cp -fr "$file" "$DIR"
done

#   Build
cd "${DOCKER_DIR}/api"
docker build -t $1 .

#!/bin/bash
#!/usr/bin/env bash
set -eo pipefail

########################################################################################################################
#
# Vars that should not be edited
#
########################################################################################################################

RED=$(printf "\33[1;31m")
GREEN=$(printf "\33[1;32m")
YELLOW=$(printf "\33[33m")
RESET=$(printf "\33[0m")
BLUE=$(printf "\33[1;34m")

MIN_DOCKER_VERSION="19.03"
MIN_COMPOSE_VERSION="1.24"
COMPILED_BASE_NAME="ty_go_poc"
APP_NAME="

8888888 8888888888 \`8.\`8888.      ,8'  ,o888888o.        ,o888888o.     8 888888888o       ,o888888o.         ,o888888o.
      8 8888        \`8.\`8888.    ,8'  8888     \`88.   . 8888     \`88.   8 8888    \`88.  . 8888     \`88.      8888     \`88.
      8 8888         \`8.\`8888.  ,8',8 8888       \`8. ,8 8888       \`8b  8 8888     \`88 ,8 8888       \`8b  ,8 8888       \`8.
      8 8888          \`8.\`8888.,8' 88 8888           88 8888        \`8b 8 8888     ,88 88 8888        \`8b 88 8888
      8 8888           \`8.\`88888'  88 8888           88 8888         88 8 8888.   ,88' 88 8888         88 88 8888
      8 8888            \`8. 8888   88 8888           88 8888         88 8 888888888P'  88 8888         88 88 8888
      8 8888             \`8 8888   88 8888   8888888 88 8888        ,8P 8 8888         88 8888        ,8P 88 8888
      8 8888              8 8888   \`8 8888       .8' \`8 8888       ,8P  8 8888         \`8 8888       ,8P  \`8 8888       .8'
      8 8888              8 8888      8888     ,88'   \` 8888     ,88'   8 8888          \` 8888     ,88'      8888     ,88'
      8 8888              8 8888       \`8888888P'        \`8888888P'     8 8888             \`8888888P'         \`8888888P'
"

RETURN_CODE=0

echo "$APP_NAME"
echo -n "${BLUE}Checking dependencie...${RESET} "

# BASE Function
# Help
help() {
  echo -e "$0 <api_version> <cmd> <args>
$COMPILED_BASE_NAME is an helper script to cleanly work within TY services.

Command:
    help        display this message
    gen_protoc  compile .proto files to .pd.go
    build       build servers for dev
    serve       run server

Use <cmd> help to have the detailed help for command
  "
}

# Version check functions ------------------------------------------------------
version_comp() {
  if [[ "$1" == "$2" ]]; then
    return 0
  fi
  local IFS=.
  local i ver1=($1) ver2=($2)
  # fill empty fields in ver1 with zeros
  for ((i = ${#ver1[@]}; i < ${#ver2[@]}; i++)); do
    ver1[i]=0
  done
  for ((i = 0; i < ${#ver1[@]}; i++)); do
    if [[ -z ${ver2[i]} ]]; then
      # fill empty fields in ver2 with zeros
      ver2[i]=0
    fi
    if ((10#${ver1[i]} > 10#${ver2[i]})); then
      return 0
    fi
    if ((10#${ver1[i]} < 10#${ver2[i]})); then
      return 2
    fi
  done
  return 0
}

test_version_comp() {
  version_comp "$1" "$2"
  case $? in
  0) op='>=' ;;
  *) op='<' ;;
  esac
  if [[ $op == '<' ]]; then
    echo -e "${RED}FAIL: Your version is older than require.,  '$1', '$2' $RESET"
    return 1
  else
    echo -e "${GREEN}Pass: '$1 $op $2'.$RESET"
    return 0
  fi
}
# ------------------------------------------------------------------------------
# Spinner
spin() {
  local -r pid="${1}"
  local delay=0.5
  spinner="/|\\-/|\\-"
  while ps a | awk '{print $1}' | grep -q "${pid}"; do
    for i in $(seq 0 7); do
      echo -n "${spinner:$i:1}"
      echo -en "\010"
      sleep $delay
    done
  done
}

# Check deps
check_dep() {
  # Test Docker installation -------------------------------------------------------
  echo
  if VERSION=$(docker version --format '{{.Server.Version}}'); then

    if test_version_comp "$VERSION" $MIN_DOCKER_VERSION; then
      echo -e "${GREEN}Docker well installed $RESET"
    else
      echo -e "${RED}Please update Docker from https://docs.docker.com/install/ $RESET"
      RETURN_CODE=1
    fi
  else
    echo -e "${RED}Please Install docker from https://docs.docker.com/install/ or make it run without sudo. Read the full doc :) $RESET"
    RETURN_CODE=1
  fi
  echo
  VERSION=0
  # ------------------------------------------------------------------------------
  # Test Docker compose installation ---------------------------------------------
  VERSION=$(docker-compose version --short)
  if [ $? -eq 0 ]; then
    echo
    test_version_comp $VERSION $MIN_COMPOSE_VERSION
    if [ $? -eq 0 ]; then
      echo -e "${GREEN}Compose well installed$RESET"
    else
      echo -e "${RED}Please update Docker compose from https://docs.docker.com/compose/install/ $RESET"
      RETURN_CODE=1
    fi
  else
    echo -e "${RED}Please Install docker compose from https://docs.docker.com/compose/install/ $RESET"
    RETURN_CODE=1
  fi
  echo
  VERSION=0

  # Ensure protoc
  if protoc --version >/dev/null 1>/dev/null 2>/dev/null; then
    echo -e "${GREEN}Protocol buffer well installed $RESET"
  else
    echo -e "${RED}Please install protocol buffer from https://github.com/protocolbuffers/protobuf/releases $RESET"
    RETURN_CODE=1
  fi
  # Ensure golang

  if go version >/dev/null 1>/dev/null 2>/dev/null; then
    echo -e "${GREEN}Go well installed $RESET"
  else
    echo -e "${RED}Please install golang from https://golang.org/dl/ $RESET"
    RETURN_CODE=1
  fi
  # Ensure go protoc generator latest version
  if command -v protoc-gen-go >/dev/null 1>/dev/null 2>/dev/null; then
    echo -e "${GREEN}Protoc go generator well installed $RESET"
  else
    GOFLAGS="" go get github.com/golang/protobuf/protoc-gen-go >/dev/null 1>/dev/null
    if command -v protoc-gen-go; then
      echo -e "${GREEN}Protoc go generator well installed $RESET"
    else
      echo -e "${RED}Could not install protoc generator $RESET"
      RETURN_CODE=1
    fi
  fi
}

# Main functions --------------------
gen_protoc() {
  if [ "$1" = "help" ]; then
    echo "gen_protoc <version> [SERVICES]

gen_protoc command build go .pb.go files from .proto definition.

parameter:
  version    api version matching path in api (v1/v2/demo/etc ...)
  SERVICES   opt: services to build. If none, build all

services: all folders in api/VERSION without the trailing plurar s
    "
    return 0
  fi
  local api_version=${1}
  shift
  local build_spec=()
  local base_path="api/${api_version}"
  for to_build in "$@"; do
    build_spec+=("$to_build")
  done
  if [ ${#build_spec[@]} -eq 0 ]; then
    # shellcheck disable=SC2045
    for file in $(ls "$base_path"/**/*.proto); do
      service="${file##*/}"
      service="${service%.proto}"
      file=${file##${base_path}/}

      echo -e "${YELLOW}Building: $service${RESET}"
      if protoc --proto_path="$base_path" --proto_path=third_party --go_out=plugins=grpc:"$base_path" "$file"; then
        echo -e "${GREEN}$service built${RESET}"
      else
        echo -e "${RED}Could not build $service${RESET}"
      fi
    done
  else
    for service in "${build_spec[@]}"; do
      ser_path="${base_path}/${service}s"
      file="${ser_path}/${service}-service.proto"

      echo -e "${YELLOW}Building: $service${RESET}"
      if protoc --proto_path="$base_path" --proto_path=third_party --go_out=plugins=grpc:"$base_path" "$file"; then
        echo -e "${GREEN}$service built${RESET}"
      else
        echo -e "${RED}Could not build $service${RESET}"
      fi
    done
  fi
}

build() {
  if [ "$1" = "help" ]; then
    echo "build [SERVICES]

build command build servers.

parameter:
  SERVICES   opt: servers to build. If none, build all

services:
    grpc     grpc server only
    rest     rest api server only
    all      res+grpc (default)
    "
    return 0
  fi
  local build_spec=()
  local build_all=0
  for to_build in "$@"; do
    [ $build_all -eq 0 ] && build_all=2
    [ $to_build = "all" ] && build_all=1
    build_spec+=("$to_build")
  done
  if [ $build_all -eq 0 ] || [ $build_all -eq 1 ]; then
    rm -R build || true
    echo "Building dev"
    GOOS=darwin GOARCH=amd64 go build -o "build/osx64/${COMPILED_BASE_NAME}-dev" dev.go
    GOOS=linux GOARCH=amd64 go build -o "build/linux64/${COMPILED_BASE_NAME}-dev" dev.go
    echo "Building grpc"
    GOOS=darwin GOARCH=amd64 go build -o "build/osx64/${COMPILED_BASE_NAME}-grpc" grpc.go
    GOOS=linux GOARCH=amd64 go build -o "build/linux64/${COMPILED_BASE_NAME}-grpc" grpc.go
    echo "Building rest"
    GOOS=darwin GOARCH=amd64 go build -o "build/osx64/${COMPILED_BASE_NAME}-rest" rest.go
    GOOS=linux GOARCH=amd64 go build -o "build/linux64/${COMPILED_BASE_NAME}-rest" rest.go
  else
    for service in "${build_spec[@]}"; do
      echo "Building grpc"
      if [ "$service" = "grpc" ]; then
        rm "build/osx64/${COMPILED_BASE_NAME}-grpc" || true
        rm "build/linux64/${COMPILED_BASE_NAME}-grpc" || true
        GOOS=darwin GOARCH=amd64 go build -o "build/osx64/${COMPILED_BASE_NAME}-grpc" grpc.go
        GOOS=linux GOARCH=amd64 go build -o "build/linux64/${COMPILED_BASE_NAME}-grpc" grpc.go
      fi
      if [ "$service" = "rest" ]; then
        echo "Building rest"
        rm "build/osx64/${COMPILED_BASE_NAME}-rest" || true
        rm "build/linux64/${COMPILED_BASE_NAME}-rest" || true
        GOOS=darwin GOARCH=amd64 go build -o "build/osx64/${COMPILED_BASE_NAME}-rest" grpc.go
        GOOS=linux GOARCH=amd64 go build -o "build/linux64/${COMPILED_BASE_NAME}-rest" grpc.go
      fi
    done
  fi
}

serve() {
  if [ "$1" = "help" ]; then
    echo "serve [SERVER]

serve command run server.

parameter:
  SERVER   opt: server to run. If none, run dev

servers:
    grpc     grpc server only
    rest     rest api server only
    dev      res+grpc (default)
    "
    return 0
  fi
  local server=${1:-dev}
  go run "${server}.go"
}

check_dep &
echo
spin $!
if [ $RETURN_CODE -ne 0 ]; then exit $RETURN_CODE; fi
echo
echo
"${@:-help}"

#!/usr/bin/env bash

set -e

main(){
    datadir="$(mktemp -d)"
    trap "rm -rf $datadir" EXIT
    ./build/bin/core-geth --datadir="$datadir" init "$1"
	./build/bin/core-geth --datadir="$datadir" import "$2"
}

main $*
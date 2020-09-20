#!/bin/bash

function backtrace() {
    local func="${FUNCNAME[1]}"
    local line="${BASH_LINENO[0]}"
    local src="${BASH_SOURCE[0]}"
    echo "  called from file $src, func $func(), line $line"
}

set -euo errtrace
trap backtrace ERR

go vet -v ./...
#go test -p 1 -v ./...
go build -v cmd/stocker/stocker.go

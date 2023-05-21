#!/bin/bash
export GOCACHE=${BOX_ROOT}/go/cache
/usr/local/go/bin/go build -o main main.go
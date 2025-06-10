#!/bin/bash
set -e
set -x

cd /Users/jrush/Development/trueblocks-dalledress
go build -o enhanced-demo cmd/enhanced-demo/main.go
./enhanced-demo

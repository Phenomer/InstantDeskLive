#!/bin/sh

set -e
set -x

GOARCH=amd64 GOOS=windows go build -o desklive.exe


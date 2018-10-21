#!/bin/bash

set -e

GOOS=darwin GOARCH=amd64 go build -o bin/lpass-auditor-osx
GOOS=windows GOARCH=amd64 go build -o bin/lpass-auditor.exe
GOOS=linux GOARCH=amd64 go build -o bin/lpass-auditor

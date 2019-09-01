#!/usr/bin/env bash

PLUGIN_NAME=$1
PLUGIN_VERSION=$2

# Linux
GOOS=linux GOARCH=amd64 go build -o ${PLUGIN_NAME}_${PLUGIN_VERSION}.linux64
GOOS=linux GOARCH=386 go build -o ${PLUGIN_NAME}_${PLUGIN_VERSION}.linux32

# Windows
GOOS=windows GOARCH=amd64 go build -o ${PLUGIN_NAME}_${PLUGIN_VERSION}.win64
GOOS=windows GOARCH=386 go build -o ${PLUGIN_NAME}_${PLUGIN_VERSION}.win32

# Mac OS X
GOOS=darwin GOARCH=amd64 go build -o ${PLUGIN_NAME}_${PLUGIN_VERSION}.osx
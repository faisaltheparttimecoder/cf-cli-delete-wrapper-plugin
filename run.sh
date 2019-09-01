#!/usr/bin/env bash

PLUGIN_VERSION=v0.1.0
PLUGIN_NAME=cf-delete-wrapper

# Create the build
/bin/sh build.sh ${PLUGIN_NAME} ${PLUGIN_VERSION}

# Install the plugin
cf install-plugin ${PLUGIN_NAME}_${PLUGIN_VERSION}.osx -f
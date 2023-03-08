#!/bin/bash

set -e

# Build the Go binaries
./build.sh

# Deploy the service
sls deploy

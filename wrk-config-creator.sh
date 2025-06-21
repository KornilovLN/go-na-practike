#!/bin/bash
# Script to create a directory structure for work-configs project
# Usage: ./wrk-config-creator.sh
# Start from the root project directory

mkdir -p cmd/wrk-configs
cd cmd/wrk-configs
mkdir -p {configs/{examples,schemas},pkg/{parsers,generators,utils,types},examples/{01-basic-json,02-basic-yaml,03-basic-ini,04-dynamic-json,05-universal-reader,06-config-manager,07-json-to-struct},cmd/{config-converter,config-validator},internal}

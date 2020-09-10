#!/bin/bash

set -e

#SCRIPT_PATH="$( cd "$(dirname "$0")" ; pwd -P )"
go build -o $PROJECTHOME/bin/5gmonitor $SOURCE_ROOT/5gsurveillance/pkg

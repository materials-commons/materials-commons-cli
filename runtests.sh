#!/bin/sh

P=github.com/materials-commons/materials
export MATERIALS_WEBDIR=""
export MATERIALS_ADDRESS=""
export MATERIALS_PORT=""
export MATERIALS_UPDATE_CHECK_INTERVAL=""
export MCDOWNLOADURL=""
export MCAPIURL=""
export MCURL=""

go test -v $P $P/site

#!/bin/sh

export MATERIALS_WEBDIR=""
export MATERIALS_ADDRESS=""
export MATERIALS_PORT=""
export MATERIALS_UPDATE_CHECK_INTERVAL=""
export MCDOWNLOADURL=""
export MCAPIURL=""
export MCURL=""

#mcfs/mcfs --db-connect='localhost:30815'&
#MCFSPID=$!
go test -v ./...
#kill -HUP $MCFSPID

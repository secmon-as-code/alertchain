#!/bin/bash

docker run -d -p 9010:9010 -p 9020:9020 gcr.io/cloud-spanner-emulator/emulator
gcloud config configurations activate $SPANNER_EMULATOR_CONFIG
gcloud spanner instances create $SPANNER_INSTANCE_ID --config=$SPANNER_EMULATOR_CONFIG --nodes=1 --description=testing
wrench create --directory ./db/spanner

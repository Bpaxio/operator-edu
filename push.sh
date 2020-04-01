#!/bin/sh
# OPERATOR_VERSION=$1
# echo $OPERATOR_VERSION
OPERATOR_VERSION=`cat ./version/version.go | grep -oE "\d+\.\d+\.\d+"`
echo $OPERATOR_VERSION
operator-sdk build bpaxio/operator_myapp:$OPERATOR_VERSION
docker login
docker tag bpaxio/operator_myapp:$OPERATOR_VERSION bpaxio/operator_myapp:latest
docker push bpaxio/operator_myapp:$OPERATOR_VERSION
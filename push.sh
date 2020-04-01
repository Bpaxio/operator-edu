#!/bin/sh
PROJECT_PATH=$1
echo $PROJECT_PATH
OPERATOR_VERSION=`cat ./${PROJECT_PATH}/version/version.go | grep -oE "\d+\.\d+\.\d+"`
echo $OPERATOR_VERSION
operator-sdk build bpaxio/operator_${PROJECT_PATH}:$OPERATOR_VERSION
docker login
docker tag bpaxio/operator_${PROJECT_PATH}:$OPERATOR_VERSION bpaxio/operator_${PROJECT_PATH}:latest
docker push bpaxio/operator_${PROJECT_PATH}:$OPERATOR_VERSION
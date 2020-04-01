#!/bin/sh
PROJECT_PATH=$1
DIR=${PWD}
cd ./$PROJECT_PATH
echo "cd to ./$PROJECT_PATH"
OPERATOR_VERSION=`cat ./version/version.go | grep -oE "\d+\.\d+\.\d+"`
echo "build operator: bpaxio/operator_${PROJECT_PATH}:$OPERATOR_VERSION"
operator-sdk build bpaxio/operator_${PROJECT_PATH}:$OPERATOR_VERSION
docker login
docker tag bpaxio/operator_${PROJECT_PATH}:$OPERATOR_VERSION bpaxio/operator_${PROJECT_PATH}:latest
docker push bpaxio/operator_${PROJECT_PATH}:$OPERATOR_VERSION
cd $DIR
echo "returned to $DIR"
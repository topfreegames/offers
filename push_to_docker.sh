#!/bin/bash

VERSION=$(cat ./metadata/version.go | grep "var Version" | awk ' { print $4 } ' | sed s/\"//g)
COMMIT=$(git rev-parse --short HEAD)

docker login -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD"
if [ $? -ne 0 ]; then
    exit 1
fi

docker build -t offers .
if [ $? -ne 0 ]; then
    exit 1
fi

docker tag offers:latest tfgco/offers:$TRAVIS_BUILD_NUMBER-v$VERSION-$COMMIT
docker tag offers:latest tfgco/offers:v$VERSION
docker tag offers:latest tfgco/offers:latest

docker push tfgco/offers:$TRAVIS_BUILD_NUMBER-v$VERSION-$COMMIT
if [ $? -ne 0 ]; then
    exit 1
fi
docker push tfgco/offers:v$VERSION
if [ $? -ne 0 ]; then
    exit 1
fi
docker push tfgco/offers:latest
if [ $? -ne 0 ]; then
    exit 1
fi

if [ "$DOCKERHUB_LATEST" != "$TRAVIS_BUILD_NUMBER-v$VERSION-$COMMIT" ]; then
    echo "Last version is not in docker hub!"
    echo "docker hub: $DOCKERHUB_LATEST, expected: $TRAVIS_BUILD_NUMBER-v$VERSION-$COMMIT"
    exit 1
fi

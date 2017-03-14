#!/bin/bash

VERSION=$(cat ./api/version.go | grep "var VERSION" | awk ' { print $4 } ' | sed s/\"//g)
COMMIT=$(git rev-parse --short HEAD)

docker login -e="$DOCKER_EMAIL" -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD"

docker build -t offers .
docker tag offers:latest tfgco/offers:$TRAVIS_BUILD_NUMBER-v$VERSION-$COMMIT
docker tag offers:latest tfgco/offers:v$VERSION
docker tag offers:latest tfgco/offers:latest
docker push tfgco/offers:$TRAVIS_BUILD_NUMBER-v$VERSION-$COMMIT
docker push tfgco/offers:v$VERSION
docker push tfgco/offers:latest

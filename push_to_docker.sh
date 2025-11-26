#!/bin/bash

set -e  # Exit on error

VERSION=$(cat ./metadata/version.go | grep "var Version" | awk ' { print $4 } ' | sed s/\"//g)
COMMIT=$(git rev-parse --short HEAD)
IMAGE_REPOSITORY=tfgco/offers
PLATFORMS=${PLATFORMS:-linux/amd64,linux/arm64}

echo "Building multi-arch binaries..."
make build-multiarch

echo "Logging into Docker Hub..."
docker login -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD"

echo "Setting up Docker buildx for multi-architecture builds..."
docker buildx create --use --name multiarch-builder 2>/dev/null || docker buildx use multiarch-builder

echo "Building and pushing multi-arch Docker images..."
echo "Platforms: ${PLATFORMS}"
echo "Tags: v${VERSION}, latest"

docker buildx build \
    --platform ${PLATFORMS} \
    --tag ${IMAGE_REPOSITORY}:v${VERSION} \
    --tag ${IMAGE_REPOSITORY}:latest \
    --push \
    .

echo "Successfully pushed multi-arch images!"
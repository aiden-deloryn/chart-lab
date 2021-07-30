#!/bin/bash

if [ "$1" == "" ]
then
    echo "ERROR: You must provide a version number, e.g: '0.0.0'"
    exit 1
else 
    echo $1
fi

if grep -Fq "${1}" config/deployment.yaml
then
    echo "Version tag verified: config/deployment.yaml"
else
    echo "ERROR: Version mismatch. Update the image tag version in config/deployment.yaml"
    exit 1
fi

# Build and push docker image
docker build -t aidendeloryn/chartlab:latest -t aidendeloryn/chartlab:${1} .
docker push aidendeloryn/chartlab:latest
docker push aidendeloryn/chartlab:${1}

# Apply and push version tag
git tag -a v${1} -m "Release version ${1}"
git push origin v${1}
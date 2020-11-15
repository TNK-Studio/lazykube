#!/bin/bash

# allow specifying different destination directory
DIR="${DIR:-"/usr/bin/"}"

# prepare the download URL
GITHUB_LATEST_VERSION=$(curl -L -s -H 'Accept: application/json' https://github.com/TNK-Studio/lazykube/releases/latest | sed -e 's/.*"tag_name":"\([^"]*\)".*/\1/')
echo "GITHUB_LATEST_VERSION ${GITHUB_LATEST_VERSION}"
GITHUB_FILE="lazykube_linux_amd64"
GITHUB_URL="https://github.com/TNK-Studio/lazykube/releases/download/${GITHUB_LATEST_VERSION}/${GITHUB_FILE}.tar.gz"
echo "GITHUB_URL ${GITHUB_URL}"

# install/update the local binary
curl -L -o lazykube.tar.gz $GITHUB_URL --progress-bar
tar -xzvf lazykube.tar.gz ${GITHUB_FILE}/lazykube
sudo mv -f ${GITHUB_FILE}/lazykube "$DIR"
echo "lazykube install to '${DIR}'"
rm -rf lazykube.tar.gz ${GITHUB_FILE}

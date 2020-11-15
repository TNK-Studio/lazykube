#!/bin/bash

# allow specifying different destination directory
DIR="${DIR:-"/usr/local/bin"}"

# prepare the download URL
GITHUB_LATEST_VERSION=$(curl -L -s -H 'Accept: application/json' https://github.com/TNK-Studio/lazykube/releases/latest | sed -e 's/.*"tag_name":"\([^"]*\)".*/\1/')
GITHUB_FILE="lazykube_linux_amd64.tar.gz"
GITHUB_URL="https://github.com/TNK-Studio/lazykube/releases/download/${GITHUB_LATEST_VERSION}/${GITHUB_FILE}"

# install/update the local binary
curl -L -o lazykube.tar.gz $GITHUB_URL --progress-bar
tar -xzvf lazykube.tar.gz lazykube
sudo mv -f lazykube "$DIR"
rm lazykube.tar.gz

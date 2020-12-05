#!/bin/bash

# allow specifying different destination directory
DIR="${DIR:-"/usr/local/bin"}"
PROXY="https://cool-moon-43e4.elfgzp.workers.dev/"

ARCH=$(uname -m)
case $ARCH in
    i386|i686|x86) ARCH=386 ;;
    armv6*) ARCH=arm ;;
    armv7*) ARCH=arm ;;
    aarch64*) ARCH=arm ;;
    x86_64) ARCH=amd64 ;;
esac

# prepare the download URL
GITHUB_LATEST_VERSION=$(curl -L -s -H 'Accept: application/json' ${PROXY}https://github.com/TNK-Studio/lazykube/releases/latest | sed -e 's/.*"tag_name":"\([^"]*\)".*/\1/')
echo "GITHUB_LATEST_VERSION ${GITHUB_LATEST_VERSION}"
GITHUB_FILE="lazykube_linux_${ARCH}"
GITHUB_URL="https://github.com/TNK-Studio/lazykube/releases/download/${GITHUB_LATEST_VERSION}/${GITHUB_FILE}.tar.gz"
echo "GITHUB_URL ${GITHUB_URL}"

# install/update the local binary
curl -L -o lazykube.tar.gz $PROXY$GITHUB_URL --progress-bar
tar -xzvf lazykube.tar.gz ${GITHUB_FILE}/lazykube
sudo mv -f ${GITHUB_FILE}/lazykube "$DIR"
echo "lazykube install to '${DIR}'"
rm -rf lazykube.tar.gz ${GITHUB_FILE}

# Compatible with previous installation methods
if [ -f "/usr/bin/lazykube" ]; then
  rm -f /usr/bin/lazykube
  ln -s ${DIR}/lazykube /usr/bin/lazykube
  echo "Add '/usr/bin/lazykube' link to '${DIR}/lazykube'"
fi

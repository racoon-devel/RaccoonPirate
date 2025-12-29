#!/bin/bash

set -e

rm -rf .release
mkdir -p .release
make clean

arch=`uname -m`
version=`git tag --sort=-version:refname | head -n 1`
script_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

echo "Buidling for host ${arch}..."
make build
cp .build/raccoon-pirate .release/raccoon-pirate-linux-${arch}
if [ "${arch}" = "x86_64" ]; then
    cp .build/raccoon-pirate .release/raccoon-pirate-linux-amd64
fi
VERSION=${version} ARCH=${arch} nfpm pkg -p deb --target .release/
VERSION=${version} ARCH=${arch} nfpm pkg -p rpm --target .release/
make clean

arch="arm64"
echo "Building for ${arch}..."
TARGET_TRIPLET=aarch64-rpi3-linux-gnu make rpi
cp .build/rpi/build/raccoon-pirate .release/raccoon-pirate-linux-${arch}
cp .build/rpi/build/raccoon-pirate .build/raccoon-pirate
VERSION=${version} ARCH=${arch} nfpm pkg -p deb --target .release/
VERSION=${version} ARCH=${arch} nfpm pkg -p rpm --target .release/
make clean

arch="aarch64"
echo "Building package for Batocera..."
TARGET_TRIPLET=aarch64-rpi3-linux-gnu make batocera
short_version="${version:1}"
cp .build/batocera/raccoon-pirate-${short_version}-1-${arch}.pkg.tar.zst .release/raccoon-pirate-batocera-${short_version}-1-${arch}.pkg.tar.zst
make clean

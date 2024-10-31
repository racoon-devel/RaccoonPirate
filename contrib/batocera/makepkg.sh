#!/bin/bash

set -ex

echo "Make Batocera package..."

if [ -z VERSION ]; then
    VERSION=v0.0.1
fi
VERSION="${VERSION:1}"

BUILD_CMD="git clone https://github.com/libfuse/libfuse.git && cd libfuse && git checkout fuse-2.9.9 && git apply ../fuse/0001-fix-aarch64-build.patch && git apply ../fuse/0002-util-ulockmgr_server-c-conditionally-define-closefrom-fix-glibc-2-34.patch && ./makeconf.sh && CC=/home/develop/opt/x-tools/${TARGET_TRIPLET}/bin/${TARGET_TRIPLET}-gcc ./configure --host=${TARGET_TRIPLET} --disable-example --enable-lib --enable-util && make && sudo cp ./lib/.libs/*.so* ../fuse && sudo chown -R ${UID}:${UID} ../fuse/"
docker run --rm -v ./fuse:/home/develop/fuse -it raccoon_pirate_rpi_build-${TARGET_TRIPLET} /bin/bash -c "${BUILD_CMD}"
PKG_DIR=package/userdata/system/raccoon_pirate
cp ${PATH_TO_BINARY} ${PKG_DIR}/raccoon-pirate
cp fuse/lib* ${PKG_DIR}
cd package && wget https://raw.githubusercontent.com/batocera-linux/batocera.linux/ced8e7abd2ce2b0833c649ddffe0cbf82a06e086/package/batocera/utils/pacman/batocera-makepkg
sed -i "s/__VERSION__/${VERSION}-1/" .PKGINFO
chmod u+x batocera-makepkg
./batocera-makepkg



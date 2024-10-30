#!/bin/bash

echo "Make Batocera package..."

BUILD_CMD="git clone https://github.com/libfuse/libfuse.git && cd libfuse && git checkout fuse-2.9.9 && git apply ../fuse/0001-fix-aarch64-build.patch && git apply ../fuse/0002-util-ulockmgr_server-c-conditionally-define-closefrom-fix-glibc-2-34.patch && ./makeconf.sh && CC=/home/develop/opt/x-tools/${TARGET_TRIPLET}/bin/${TARGET_TRIPLET}-gcc ./configure --host=${TARGET_TRIPLET} --disable-example --enable-lib --enable-util && make && sudo cp ./lib/.libs/*.so* ../fuse && sudo chown -R ${UID}:${UID} ../fuse/"
docker run --rm -v ./fuse:/home/develop/fuse -it raccoon_pirate_rpi_build /bin/bash -c "${BUILD_CMD}"



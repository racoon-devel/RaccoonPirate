#!/bin/bash

set -ex

if [ -z "${TARGET_TRIPLET}" ]; then
    echo "TARGET_TRIPLET must be set"
    exit 1
fi

if [ -z "${BINARY_NAME}" ]; then
    echo "BINARY_NAME must be set"
    exit 1
fi

if [ -z "${SOURCE_DIR}" ]; then
    echo "SOURCE_DIR must be set"
    exit 1
fi

if [ -z "${SOURCE_MAIN}" ]; then
    echo "SOURCE_DIR must be set"
    exit 1
fi

# Build docker container
sed -i "s/__TARGET_TRIPLET__/${TARGET_TRIPLET}/" Dockerfile
docker build --tag raccoon_pirate_rpi_build-${TARGET_TRIPLET} .

GO_OPTS="GOARCH=arm"
if [ "${TARGET_TRIPLET}" == "aarch64-rpi3-linux-gnu" ]; then
    GO_OPTS="GOARCH=arm64"
fi
if [ "${TARGET_TRIPLET}" == "armv8-rpi3-linux-gnueabihf" ]; then
    GO_OPTS="GOARCH=arm"
fi
if [ "${TARGET_TRIPLET}" == "armv6-rpi-linux-gnueabihf" ]; then 
    GO_OPTS="GOARCH=arm GOARM=6"
fi

echo ${GO_OPTS}
	
# Build raccoon-pirate
RPI_TOOLCHAIN="/home/develop/opt/x-tools/${TARGET_TRIPLET}/bin/${TARGET_TRIPLET}"
BUILD_CMD="cd /home/develop/RaccoonPirate && sudo CGO_ENABLED=1 ${GO_OPTS} CC=${RPI_TOOLCHAIN}-gcc CXX=${RPI_TOOLCHAIN}-g++ go build -ldflags '${LDFLAGS}' -o /home/develop/build/${BINARY_NAME} ${SOURCE_MAIN} && sudo chown -R ${UID}:${UID} /home/develop/build"
docker run --rm -v ./build:/home/develop/build -v ${SOURCE_DIR}:/home/develop/RaccoonPirate:ro -i raccoon_pirate_rpi_build-${TARGET_TRIPLET} /bin/bash -c "${BUILD_CMD}"
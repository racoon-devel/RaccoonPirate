#!/bin/bash

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
docker build --tag raccoon_pirate_rpi_build .
	
# Build raccoon-pirate
RPI_TOOLCHAIN="/home/develop/opt/x-tools/${TARGET_TRIPLET}/bin/${TARGET_TRIPLET}"
BUILD_CMD="cd /home/develop/RaccoonPirate && CGO_ENABLED=1 GOARCH=arm64 CC=${RPI_TOOLCHAIN}-gcc CXX=${RPI_TOOLCHAIN}-g++ sudo go build -ldflags '${LDFLAGS}' -o /home/develop/build/${BINARY_NAME} ${SOURCE_MAIN} && sudo chown -R ${UID}:${UID} /home/develop/build"
echo "BUILD_CMD='${BUILD_CMD}'"
docker run --rm -v ./build:/home/develop/build -v ${SOURCE_DIR}:/home/develop/RaccoonPirate:ro -it raccoon_pirate_rpi_build /bin/bash -c "${BUILD_CMD}"
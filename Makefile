PROJECT_NAME=raccoon-pirate
BINARY_NAME=${PROJECT_NAME}
SOURCE_MAIN=app/${PROJECT_NAME}/${PROJECT_NAME}.go
LDFLAGS="-X main.Version=`git tag --sort=-version:refname | head -n 1`"

all: build

build:
	mkdir -p .build
	go build -ldflags ${LDFLAGS} -o .build/${BINARY_NAME} ${SOURCE_MAIN}

rpi:
	PWD=`pwd`
	mkdir -p .build
	rm -rf .build.rpi
	cp -r contrib/rpi .build/rpi
	cd .build/rpi && TARGET_TRIPLET=${TARGET_TRIPLET} BINARY_NAME=${BINARY_NAME} LDFLAGS=${LDFLAGS} SOURCE_DIR=$(PWD) SOURCE_MAIN=${SOURCE_MAIN} ./build.sh

batocera: rpi
	rm -rf .build/batocera
	cp -r contrib/batocera .build/batocera
	cd .build/batocera && TARGET_TRIPLET=${TARGET_TRIPLET} ./makepkg.sh

clean:
	rm -rf .build
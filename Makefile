PROJECT_NAME=raccoon-pirate
BINARY_NAME=${PROJECT_NAME}.out
SOURCE_MAIN=app/${PROJECT_NAME}/${PROJECT_NAME}.go
LDFLAGS="-X main.Version=`git tag --sort=-version:refname | head -n 1`"

all: build

build:
	go build -ldflags ${LDFLAGS} -o ${BINARY_NAME} ${SOURCE_MAIN}

run:
	go build -ldflags ${LDFLAGS} -o ${BINARY_NAME} ${SOURCE_MAIN}
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}
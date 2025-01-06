BUILD_APP ?= scrapper
EXECUTABLE_WINDOWS ?= ${BUILD_APP}_windows_amd64.exe
EXECUTABLE_LINUX ?= ${BUILD_APP}_linux_amd64
EXECUTABLE_FOLDER ?= bin

APP_REPO=github.com/pkierski/wokanda-scrapper
APP_MAIN_DIR ?= cmd/${BUILD_APP}
APP=${APP_REPO}/${APP_MAIN_DIR}

# APP_VERSION ?= $(shell git describe --tags --always --dirty)
# CURRENT_ISO_TIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
# LINKER_FLAGS = -X ${APP}internal/buildinfo.VersionStr=${APP_VERSION} -X ${APP}internal/buildinfo.CompileTimeStr=${CURRENT_ISO_TIME}
LINKER_FLAGS=
USE_CGO=0

test:
	go test ./... -race -v

build: build-windows build-linux

build-windows:
	env CGO_ENABLED=${USE_CGO} GOOS=windows GOARCH=amd64 go build -ldflags="${LINKER_FLAGS}" -o ${EXECUTABLE_FOLDER}/${EXECUTABLE_WINDOWS} ${APP}

build-linux:
	env CGO_ENABLED=${USE_CGO} GOOS=linux GOARCH=amd64 go build -ldflags="${LINKER_FLAGS}" -o ${EXECUTABLE_FOLDER}/${EXECUTABLE_LINUX} ${APP}

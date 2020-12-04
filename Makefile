GO111MODULE=on
OS=linux
ARCH=amd64

default: build

build:
	env GOOS=${OS} GOARCH=${ARCH} go build

build-windows:
	env GOOS=windows GOARCH=amd64 go build
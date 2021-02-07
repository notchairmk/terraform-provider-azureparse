GO111MODULE=on

default: build

build:
	go build

build-windows:
	env GOOS=windows GOARCH=amd64 go build

install:
	go install

develop: install
	mkdir -p ~/.terraform.d/plugins/example.com/test/azureparse/1.0.0/darwin_amd64
	ln -sf ~/go/bin/terraform-provider-azureparse ~/.terraform.d/plugins/example.com/test/azureparse/1.0.0/darwin_amd64/terraform-provider-azureparse


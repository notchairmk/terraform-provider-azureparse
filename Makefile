default: build

.PHONY: build
build:
	go build -o bin/terraform-provider-azureparse main.go

.PHONY: build-windows
build-windows:
	env GOOS=windows GOARCH=amd64 go build

.PHONY: install
install:
	go install

.PHONY: develop
develop: install
	mkdir -p ~/.terraform.d/plugins/example.com/test/azureparse/1.0.0/darwin_amd64
	ln -sf ~/go/bin/terraform-provider-azureparse ~/.terraform.d/plugins/example.com/test/azureparse/1.0.0/darwin_amd64/terraform-provider-azureparse

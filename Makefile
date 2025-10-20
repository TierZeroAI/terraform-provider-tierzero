default: build

.PHONY: build
build:
	go build -o terraform-provider-tierzero

.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/tierzero/tierzero/1.0.0/darwin_arm64
	cp terraform-provider-tierzero ~/.terraform.d/plugins/registry.terraform.io/tierzero/tierzero/1.0.0/darwin_arm64/

.PHONY: test
test:
	go test ./... -v

.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v -timeout 120m

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: clean
clean:
	rm -f terraform-provider-tierzero
	rm -rf dist/

.PHONY: docs
docs:
	tfplugindocs generate

.PHONY: tidy
tidy:
	go mod tidy

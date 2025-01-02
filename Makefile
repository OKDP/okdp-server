.PHONY: help
help:
	@echo "Available commands:"
	@echo "  format      Format the project code"
	@echo "  lint        Run the configured linters (before pushing)"
	@echo "  compile     Compile the project code"
	@echo "  test        Run the unit tests"
	@echo "  build       Build a binary"
	@echo "  run         Run the project locally"

generate: tools gogenerate
format: generate gofmt
lint: format golint
compile: lint gocompile
test: compile gotest
build: test gobuild
run: test gorun
rundev: generate gocompile gotest gobuild gorun

.PHONY: tools
tools:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2
	go install golang.org/x/tools/cmd/goimports@v0.28.0

.PHONY: gogenerate
gogenerate:
    # Generate Open API V3 model (all types)
	oapi-codegen -package=_api -generate types,skip-prune \
	             -o api/openapi/v3/_api/types.go \
				 api/openapi/v3/api.yaml
	
    # Generate Open API server interfaces
    # https://github.com/oapi-codegen/oapi-codegen/issues/1513
	for path in api/openapi/v3/paths/*; do \
	    tag=$${path##*/}; \
		echo "Generate Server Interface for $${tag} path" ; \
		mkdir -p api/openapi/v3/_api/$${tag}/  ; \
	    oapi-codegen -package=_$${tag} -generate gin-server,types \
	                 -o api/openapi/v3/_api/$${tag}/$${tag}-server.go \
				     -include-tags $${tag}  \
				     api/openapi/v3/api.yaml  ; \
    done

    # Generate swagger spec (doc)
	oapi-codegen -package=_api -generate spec \
	             -o api/openapi/v3/_api/spec.go \
				 api/openapi/v3/api.yaml

    # Generate client stubs
	# oapi-codegen -package=_api -generate client \
	#              -o api/openapi/v3/_api/client.go \
	# 			   api/openapi/v3/api.yaml

.PHONY: gofmt
gofmt:
	goimports  -w ./internal/ ./cmd/

.PHONY: golint
golint:
	golangci-lint run

.PHONY: gobuild
gobuild:
	mkdir -p .bin/
	CGO_ENABLED=0 go build -a -o .bin/okdp-server main.go

.PHONY: gocompile
gocompile:
	CGO_ENABLED=0 go build -a -o /dev/null main.go

.PHONY: gotest
gotest:
	go test ./... -v

.PHONY: gorun
gorun:
	go run *.go --config=.local/application-local.yaml


generate: install gogenerate
compile: generate gocompile
test: compile gotest
build: test gobuild
run: test gorun

.PHONY: install
install:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1

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

.PHONY: lint
lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.61.0
	golangci-lint run ./ --out-format=colored-line-number --timeout=5m

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


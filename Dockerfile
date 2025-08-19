ARG GO_VERSION=1.24.5

FROM golang:${GO_VERSION} AS go-build

ARG GIT_COMMIT="_unset_"
ARG LDFLAGS="-X localbuild=true"
ARG TARGETOS="linux"
ARG TARGETARCH

WORKDIR /workspace/okdp-server

COPY Makefile Makefile
COPY go.* ./
COPY *.go ./
COPY api/ api/
COPY internal/ internal/
COPY cmd/ cmd/

RUN make generate
RUN go mod tidy \
    && go mod download
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    LDFLAGS=${LDFLAGS##-X localbuild=true} GIT_COMMIT=$GIT_COMMIT \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -o okdp-server main.go

FROM alpine:3.21.3

RUN apk --no-cache add ca-certificates && update-ca-certificates

COPY --from=go-build /workspace/okdp-server /usr/local/bin/

USER 65534:65534

EXPOSE 8090

ENTRYPOINT ["okdp-server"]


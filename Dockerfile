# builder: golang
FROM golang:1.19-alpine as builder

WORKDIR /.build

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# see https://medium.com/@diogok/on-golang-static-binaries-cross-compiling-and-plugins-1aed33499671
RUN CGO_ENABLED=0 go build -tags netgo -ldflags '-w -extldflags "-static"' -o ./bin/healthcheck cmd/healthcheck.go

# final layer
FROM scratch as final

COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /.build/bin/healthcheck /bin/healthcheck

ARG BUILD_COMMIT=dev
ENV BUILD_COMMIT=$BUILD_COMMIT
ARG BUILD_VERSION=dev
ENV BUILD_VERSION=$BUILD_VERSION

ENTRYPOINT ["/bin/healthcheck"]

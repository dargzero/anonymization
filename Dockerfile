FROM golang:1.10-alpine as builder

# Install tools required to build the project
# We need to run `docker build --no-cache .` to update those dependencies
RUN apk add --no-cache git
RUN go get github.com/golang/dep/cmd/dep

# Build the healthcheck, since it very rarely changes, so it is good to hav it cached early
COPY ./src/healthcheck/ /go/src/healthcheck/
WORKDIR /go/src/healthcheck/
# Install library dependencies
RUN dep ensure -vendor-only
# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/healthcheck

# Gopkg.toml and Gopkg.lock lists project dependencies
# These layers are only re-built when Gopkg files are updated
COPY ./src/anonmodel/Gopkg.lock ./src/anonmodel/Gopkg.toml /go/src/anonmodel/
WORKDIR /go/src/anonmodel/
# Install library dependencies
RUN dep ensure -vendor-only

COPY ./src/anondb/Gopkg.lock ./src/anondb/Gopkg.toml /go/src/anondb/
WORKDIR /go/src/anondb/
# Install library dependencies
RUN dep ensure -vendor-only

COPY ./src/anonbll/Gopkg.lock ./src/anonbll/Gopkg.toml /go/src/anonbll/
WORKDIR /go/src/anonbll/
# Install library dependencies
RUN dep ensure -vendor-only

COPY ./src/swagger/Gopkg.lock ./src/swagger/Gopkg.toml /go/src/swagger/
WORKDIR /go/src/swagger/
# Install library dependencies
RUN dep ensure -vendor-only

COPY ./src/server/Gopkg.lock ./src/server/Gopkg.toml /go/src/server/
WORKDIR /go/src/server/
# Install library dependencies
RUN dep ensure -vendor-only

# Copy all project and build it
# This layer is rebuilt when ever a file has changed in the project directory
COPY ./src /go/src

# Build all the projects
WORKDIR /go
RUN go install -v anonmodel
RUN go install -v anondb
RUN go install -v anonbll
RUN go install -v swagger
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/server server

# Run unit tests
RUN go test anonbll -v
RUN go test swagger -v


# This results in a single layer image
FROM scratch as final
WORKDIR /bin

# Healthcheck pings the webserver periodically
COPY --from=builder /go/bin/healthcheck healthcheck
HEALTHCHECK --interval=60s --timeout=45s --start-period=5s --retries=3 CMD [ "healthcheck" ]

# Entry point is the Go webserver
COPY --from=builder /go/bin/server server
CMD [ "server" ]
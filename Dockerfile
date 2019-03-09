FROM golang:1.10-alpine as anon-server

# tools
RUN apk add --no-cache git

# dependencies
RUN go get github.com/globalsign/mgo
RUN go get github.com/gorilla/mux
RUN go get github.com/satori/go.uuid
RUN go get -u -t gonum.org/v1/gonum/...
RUN go get bitbucket.org/dargzero/k-anon

# build
COPY ./ /go/src/github.com/dargzero/anonymization
WORKDIR /go/src
RUN go install github.com/dargzero/anonymization/healthcheck
RUN go install github.com/dargzero/anonymization/server
RUN go test github.com/dargzero/anonymization/...

# start
WORKDIR /go/bin
HEALTHCHECK --interval=60s --timeout=45s --start-period=5s --retries=3 CMD [ "healthcheck" ]
CMD [ "server" ]

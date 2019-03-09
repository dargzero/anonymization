FROM golang:1.10-alpine as anon-server

# install tools
RUN apk add --no-cache git

# go get dependencies
RUN go get github.com/globalsign/mgo
RUN go get github.com/gorilla/mux
RUN go get github.com/satori/go.uuid
RUN go get -u -t gonum.org/v1/gonum/...

# go get graph based anonymizer library
RUN go get bitbucket.org/dargzero/k-anon

# copy project files
COPY ./ /go/src/github.com/dargzero/anonymization

# build server
WORKDIR /go/src
RUN go install github.com/dargzero/anonymization/server
RUN go test github.com/dargzero/anonymization/...

# start anonimization server
WORKDIR /go/bin
CMD [ "server" ]
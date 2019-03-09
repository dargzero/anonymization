FROM golang:1.12-stretch as anon-server

COPY ./ /anonymization
WORKDIR /anonymization
RUN go install github.com/dargzero/anonymization/healthcheck
RUN go install github.com/dargzero/anonymization/server
RUN go test ./...

HEALTHCHECK --interval=60s --timeout=45s --start-period=5s --retries=3 CMD [ "healthcheck" ]
CMD [ "server" ]

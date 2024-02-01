# ---------------------------------------------------------------------
#  The first stage container, for building the application
# ---------------------------------------------------------------------
FROM golang:1.21.2-alpine as builder

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=linux

RUN apk --no-cache add ca-certificates
RUN apk add --update git

RUN mkdir -p $GOPATH/src/github.com/aopoltorzhicky/astria/

COPY ./go.* $GOPATH/src/github.com/aopoltorzhicky/astria/
WORKDIR $GOPATH/src/github.com/aopoltorzhicky/astria
RUN go mod download

COPY cmd/api cmd/api
COPY internal internal
COPY pkg pkg

WORKDIR $GOPATH/src/github.com/aopoltorzhicky/astria/cmd/api/
RUN go build -a -installsuffix cgo -o /go/bin/api .

# ---------------------------------------------------------------------
#  The second stage container, for running the application
# ---------------------------------------------------------------------
FROM scratch

WORKDIR /app/api

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/api /go/bin/api
COPY ./configs/dipdup.yml ./
COPY database database

ENTRYPOINT ["/go/bin/api", "-c", "dipdup.yml"]
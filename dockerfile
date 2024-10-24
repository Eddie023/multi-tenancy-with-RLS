FROM golang:1.23 as build
ARG VERISON 
WORKDIR /go/src/server

COPY go.mod go.sum ./
RUN go mod download 

ADD . .
RUN go install ./cmd/main.go

FROM debian:bookworm
ARG VERSION 
RUN apt-get update && apt-get install -y -q --no-install-recommends 

COPY --from=build /go/bin/main /bin/apiserver
ENTRYPOINT [ "apiserver" ]

LABEL image.authors="Manish Chaulagain" image.version=${VERSION}
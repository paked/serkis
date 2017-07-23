FROM golang:1.8.3

MAINTAINER Harrison Shoebridge <harrison@theshoebridges.com>

ADD . /go/src/github.com/paked/serkis

WORKDIR /go/src/github.com/paked/serkis

RUN go install

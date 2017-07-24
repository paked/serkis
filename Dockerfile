FROM golang:1.8.3

MAINTAINER Harrison Shoebridge <harrison@theshoebridges.com>

ADD . /go/src/github.com/paked/serkis

WORKDIR /go/src/github.com/paked/serkis

RUN go get github.com/gorilla/mux
RUN go get github.com/paked/configure
RUN go get github.com/russross/blackfriday

RUN go install

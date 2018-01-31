FROM golang:1.9.2

ENV GOPATH /go:$GOPATH
ENV PATH /go/bin:$PATH

ADD . /go/src
WORKDIR /go/src

RUN go get github.com/dispatchlabs/disgover
RUN go get github.com/nic0lae/JerryMouse/Servers
ENTRYPOINT go run node2.go
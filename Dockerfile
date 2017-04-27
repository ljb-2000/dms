FROM alpine:3.4

RUN apk add --update go git

ENV GOROOT=/usr/lib/go
ENV GOPATH=/go
ENV GOBIN=/go/bin
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin:/usr/local/bin

WORKDIR /go/src/github.com/lavrs/dms
ADD . /go/src/github.com/lavrs/dms

RUN go get github.com/tools/godep \
    && godep restore \
    && go install \
    && go build

EXPOSE 4222

ENTRYPOINT ["/go/bin/dms"]
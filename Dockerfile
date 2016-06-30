FROM golang:1.6
MAINTAINER Octoblu, Inc. <docker@octoblu.com>

WORKDIR /go/src/github.com/octoblu/image-upload-to-dropbox-and-post-to-slack
COPY . /go/src/github.com/octoblu/image-upload-to-dropbox-and-post-to-slack

RUN env CGO_ENABLED=0 go build -o image-upload-to-dropbox-and-post-to-slack -a -ldflags '-s' .

CMD ["./image-upload-to-dropbox-and-post-to-slack"]

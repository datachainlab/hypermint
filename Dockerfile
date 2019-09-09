FROM golang:1.13-alpine as builder

RUN apk add --no-cache make gcc libc-dev git

ADD . /go/src/github.com/bluele/hypermint
ENV WORKDIR=/go/src/github.com/bluele/hypermint
ENV GO111MODULE=on
WORKDIR ${WORKDIR}
RUN make build

FROM alpine:3.10

ENV WORKDIR=/go/src/github.com/bluele/hypermint
COPY --from=builder ${WORKDIR}/build/hmd /
COPY --from=builder ${WORKDIR}/build/hmcli /

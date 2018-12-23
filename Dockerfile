FROM golang:1.11-alpine as builder

RUN apk add --no-cache make gcc libc-dev git
RUN go get -u github.com/golang/dep/cmd/dep

ADD . /go/src/github.com/bluele/hypermint
ENV WORKDIR=/go/src/github.com/bluele/hypermint
WORKDIR ${WORKDIR}
RUN dep ensure
RUN make build init

FROM alpine:3.8

ENV WORKDIR=/go/src/github.com/bluele/hypermint
COPY --from=builder ${WORKDIR}/build/hmd /
COPY --from=builder ${WORKDIR}/build/hmcli /
COPY --from=builder /root/.hmd /root/.hmd
COPY --from=builder /root/.hmcli /root/.hmcli

CMD [ "/hmd", "start", "--log_level", "main:error", "--home", "/root/.hmd" ]

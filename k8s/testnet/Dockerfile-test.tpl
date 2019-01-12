FROM bluele/hypermint:<VERSION>

ENV WORKDIR=/go/src/github.com/bluele/hypermint
WORKDIR ${WORKDIR}

RUN apk add bash
RUN /hmd testnet -v=<VALS_NUM> --address=1 -o=/mytestnet

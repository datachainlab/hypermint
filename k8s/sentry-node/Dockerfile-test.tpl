FROM bluele/hypermint:<VERSION>

ENV WORKDIR=/go/src/github.com/bluele/hypermint
WORKDIR ${WORKDIR}

RUN apk add bash
RUN /hmd testnet -v=<VALS_NUM> -n=<VALS_NUM> --address=<GENESIS_ADDR> -o=/mytestnet

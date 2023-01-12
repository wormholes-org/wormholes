FROM alpine:latest

EXPOSE 8545

EXPOSE 30303

RUN apk update && apk add make git go linux-headers gcc musl-dev

ENV CGO_ENABLED=1

RUN git clone https://github.com/wormholes-org/wormholes

RUN cd wormholes && make wormholes

RUN mkdir -p /app/wormholes && cp -r /wormholes/build/bin/. /app/wormholes

WORKDIR /app/wormholes

CMD ["./wormholes", "--devnet", "--datadir", ".wormholes", "--mine", "--syncmode=full", "--http", "--http.addr", "0.0.0.0", "--http.port", "8545"]



from golang:1.11.2-stretch

WORKDIR /go/src/sentry

COPY . .

RUN make install

ENTRYPOINT [ "sentry" ]
FROM golang:1.11.2-stretch
ENV GO111MODULE=on
ENV GOBIN=/usr/local/bin

WORKDIR /go/src/sentry

COPY . .

RUN make install

ENTRYPOINT [ "sentry" ]
FROM golang:alpine

COPY . /go/src/github.com/madhanshanmugam/app-proxy

WORKDIR /go/src/github.com/madhanshanmugam/app-proxy

RUN go build main/proxy.go

CMD ["https://facebook.com"]

ENTRYPOINT ["./proxy"]

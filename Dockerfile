FROM golang:1.16-alpine3.12 AS builder

ENV GO111MODULE=on

#WORKDIR /go/src/github.com/gusaul/grpcox
WORKDIR /data

COPY . .

RUN export GOPROXY=https://proxy.golang.com.cn,direct && go mod vendor -v
RUN go build -o grpcox grpcox.go


FROM alpine

COPY ./index /index
COPY --from=builder /data/grpcox ./
RUN mkdir /log
EXPOSE 6969
ENTRYPOINT ["./grpcox"]

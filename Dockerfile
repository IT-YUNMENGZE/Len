# Build the manager binary
FROM golang:1.16

ENV GO111MODULE "on"
ENV GOPROXY "https://goproxy.cn,direct"
RUN go env -w GO111MODULE=on &&\
    go env -w GOPROXY=https://goproxy.cn,direct &&\
    GO111MODULE=on GOPROXY=https://goproxy.cn,direct go list -m -json -versions golang.org/x/text@latest
RUN go env

WORKDIR /

COPY bin/len /usr/local/bin

CMD ["len"]
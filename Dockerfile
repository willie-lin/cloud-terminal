FROM golang:alpine as builder

ENV GO111MODULE=on

ENV GOPROXY=https://goproxy.io,direct

WORKDIR /app

COPY . .

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add gcc g++
RUN go env && CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -ldflags '-linkmode external -extldflags "-static"' -o cloud-terminal main.go


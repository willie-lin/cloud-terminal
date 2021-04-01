FROM golang:1.15 as build

MAINTAINER Linjiezui "linjiezui@gmail.com"

# 容器环境变量添加，会覆盖默认的变量值
ENV GO111MODULE=on
#ENV GOPROXY=https://goproxy.io,direct
ENV GOPROXY=https://goproxy.cn,direct

# 设置工作区
WORKDIR /go/release


# 把全部文件添加到/go/release目录
ADD . .

# 编译：把cmd/main.go编译成可执行的二进制文件，命名为app
#RUN GOOS=linux CGO_ENABLED=1 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o cloud-terminal main.go
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build  -o cloud-terminal main.go
#RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -ldflags '-linkmode external -extldflags "-static"' -o cloud-terminal main.go
#RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -ldflags '-linkmode external -extldflags "-static"' -o cloud-terminal main.go
# 运行：使用scratch作为基础镜像
FROM scratch as prod

# 在build阶段复制时区到
COPY --from=build /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# 在build阶段复制可执行的go二进制文件app
COPY --from=build /go/release/cloud-terminal /

#COPY --from=build /go/release/terminal /
# 在build阶段复制配置文件
#COPY --from=build /go/release/config ./config

# Set environment variables
#ENV PORT=2021
ENV HOST=0.0.0.0

# Expose default port
EXPOSE 2021 3306
# 启动服务
CMD ["./cloud-terminal"]
#ENTRYPOINT ["./app"]
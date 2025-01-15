#FROM scratch
#
## 设置维护者信息
#MAINTAINER Linjiezui "JieZui.Lin@gmail.com"
#
## 复制时区信息
#COPY /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
#
## 复制已编译的二进制文件到根目录
#COPY dist/cloud-terminal_linux_amd64_v1/cloud-terminal /
#
## 复制配置文件到 /config 目录
#COPY app/config/config.json /app/config/
#
## 复制证书文件到 /cert 目录
#COPY cert/cert.pem /cert/cert.pem
#COPY cert/key.pem /cert/key.pem
#
## 设置环境变量（可选）
#ENV HOST=0.0.0.0
#
## 公开端口
#EXPOSE 80 443 3306
#
## 启动服务
#CMD ["./cloud-terminal"]
#

FROM golang:1.23-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制所有文件到工作目录
COPY . .

# 编译 Go 程序
RUN go build -o cloud-terminal ./cmd/cloud-terminal

# 使用 Caddy 作为反向代理服务器
FROM caddy:2

# 设置工作目录
WORKDIR /app

# 复制后端编译的二进制文件到工作目录
COPY --from=builder /app/cloud-terminal /app/cloud-terminal

# 复制 Caddyfile 到 Caddy 配置目录
COPY Caddyfile /etc/caddy/Caddyfile

# 复制配置文件到 /app/config 目录
COPY app/config/config.json /app/config/config.json

# 暴露端口
EXPOSE 80
EXPOSE 443

# 启动后端应用程序和 Caddy
CMD ["sh", "-c", "./cloud-terminal & caddy run --config /etc/caddy/Caddyfile --adapter caddyfile"]


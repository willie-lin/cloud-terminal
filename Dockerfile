FROM scratch

# 设置维护者信息
MAINTAINER Linjiezui "JieZui.Lin@gmail.com"

# 复制时区信息
#COPY /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# 复制已编译的二进制文件到根目录
COPY dist/cloud-terminal_linux_amd64_v1/cloud-terminal /

# 复制配置文件到 /config 目录
#COPY app/config/config.json /app/config/

# 设置环境变量（可选）
ENV HOST=0.0.0.0

# 公开端口
EXPOSE 8080 3306

# 启动服务
CMD ["./cloud-terminal"]


## 使用 Caddy 作为基础镜像
#FROM caddy:2
#
## 设置工作目录
#WORKDIR /app
#
## 复制编译好的二进制文件到工作目录
#COPY dist/cloud-terminal_linux_amd64_v1/cloud-terminal /app/cloud-terminal
#
## 复制配置文件到 /app/config 目录
##COPY app/config/config.json /app/config/config.json
#
## 复制 Caddyfile 到 Caddy 配置目录
#COPY Caddyfile /etc/caddy/Caddyfile
#
## 复制静态资源文件到工作目录
#COPY picture /app/picture
#
## 暴露端口
#EXPOSE 80
#EXPOSE 443
#
## 启动应用程序和 Caddy
#CMD ["sh", "-c", "./cloud-terminal & caddy run --config /etc/caddy/Caddyfile --adapter caddyfile"]

# 使用官方的 Node 镜像作为基础镜像
FROM node:23-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 package.json 和 yarn.lock 文件到工作目录
COPY package.json yarn.lock ./

# 安装项目依赖
RUN yarn install

# 复制所有文件到工作目录
COPY . .

# 构建项目
RUN NODE_OPTIONS="--max_old_space_size=4096" yarn build



## 使用 Caddy 作为静态文件服务器
#FROM caddy:2
#
## 复制构建结果到 Caddy 服务器目录
#COPY --from=0 /app/build /usr/share/caddy
#
## 复制 Caddyfile 到 Caddy 配置目录
#COPY Caddyfile /etc/caddy/Caddyfile
#
## 暴露端口
#EXPOSE 80
#EXPOSE 443
#
## 启动 Caddy
#CMD ["caddy", "run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]

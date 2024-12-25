# 使用官方的 Node 镜像作为基础镜像
FROM node:14 as build

# 设置工作目录
WORKDIR /app

# 复制 package.json 和 package-lock.json
COPY package*.json ./

# 安装项目依赖
RUN npm install

# 复制项目文件
COPY . .

# 编译项目
RUN npm run build

# 使用 Nginx 作为生产环境的服务器
FROM nginx:alpine

# 复制编译后的文件到 Nginx 目录
COPY --from=build /app/build /usr/share/nginx/html

# 复制自定义 Nginx 配置文件（可选）
COPY nginx.conf /etc/nginx/nginx.conf

# 暴露端口
EXPOSE 80

# 启动 Nginx
CMD ["nginx", "-g", "daemon off;"]

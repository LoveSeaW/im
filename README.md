# fim — 即时通讯项目（Fengfeng Instant Messaging）

一个基于 go-zero 微服务架构 + Vue3 的全功能即时通讯系统。

## 目录结构

```
im/
├── backend/                       # Go 后端（go-zero 微服务）
│   ├── main.go                    # 数据库表结构初始化入口
│   ├── go.mod
│   ├── Dockerfile                 # 基础镜像 fim_server，统一编译所有服务二进制
│   ├── core/                      # gorm/redis/etcd 等通用初始化
│   ├── common/  utils/  template/ # 公共代码 / 工具 / 模板
│   ├── fim_gateway/               # 网关（统一入口 :8080）
│   ├── fim_auth/                  # 认证服务（auth_api）
│   ├── fim_user/                  # 用户服务（user_api + user_rpc）
│   ├── fim_chat/                  # 单聊服务（chat_api + chat_rpc）
│   ├── fim_group/                 # 群聊服务（group_api + group_rpc）
│   ├── fim_file/                  # 文件服务（file_api + file_rpc）
│   ├── fim_settings/              # 系统设置（settings_api + settings_rpc）
│   ├── fim_logs/                  # 日志服务（logs_api，依赖 kafka）
│   ├── develop/docker-compose.yaml  # 本地依赖：mysql/redis/etcd/kafka
│   └── deploy/                    # 生产编排：docker-compose.yaml + 每个服务的 Dockerfile/yaml + nginx
└── frontend/                      # Vue3 + Vite + TypeScript + Element Plus
    ├── package.json
    ├── vite.config.ts             # 开发模式下 /api 代理到后端网关
    ├── .env / .env.dev1           # VITE_SERVER_URL 配置
    └── src/
```

## 技术栈

- **后端**：Go 1.21+、go-zero、gorm、go-redis、gorilla/websocket、etcd、Kafka、MySQL 8、Redis 5
- **前端**：Vue 3、Vite、TypeScript、Element Plus、Pinia、Vue Router、video.js、vue-cropper
- **基础设施**：Docker、Docker Compose、Nginx

---

## 一、本地启动（开发模式）

### 1. 启动基础组件（MySQL / Redis / etcd / Kafka）

进入 `backend/develop/`，使用 docker-compose 一键拉起所有依赖：

```bash
cd backend/develop
docker compose up -d
```

启动后服务暴露端口：

| 组件   | 容器内端口 | 宿主端口 | 说明                                    |
|--------|------------|----------|-----------------------------------------|
| MySQL  | 3306       | 3307     | root / root，库名 `fim_server_db`       |
| Redis  | 6379       | 6378     | 无密码                                  |
| etcd   | 2379       | 2379     | 无认证                                  |
| Kafka  | 9092       | 9092     | 通过 zookeeper 注册                     |

> 注意：宿主机端口为 3307（MySQL）和 6378（Redis），与默认端口不同，请同步修改下方 `etc/*.yaml` 配置或使用容器网络地址。

### 2. 修改各服务配置

每个服务的配置文件位于 `backend/<service>/<api|rpc>/etc/*.yaml`，主要修改：

- `Mysql.DataSource`：本地访问填 `root:root@tcp(127.0.0.1:3307)/fim_server_db?charset=utf8mb4&parseTime=True&loc=Local`
- `Redis.Addr`：`127.0.0.1:6378`
- `Etcd.Hosts` / `Etcd`：`127.0.0.1:2379`
- 网关 `backend/fim_gateway/settings.yaml` 中 `addr` 改为本机 IP（如 `127.0.0.1:8080`），`etcd` 同上

### 3. 初始化数据库表结构

```bash
cd backend
go mod tidy
go run main.go -db
```

> `main.go` 中默认连 `127.0.0.1:3306`，本地若使用 docker-compose 暴露的 3307 请先修改连接串再运行，或直接执行 `backend/deploy/mysql/db.sql` 导入完整结构。

### 4. 启动后端微服务

在 `backend/` 目录下，按以下顺序启动各 RPC，再启动各 API，最后启动网关。每条命令请新开一个终端：

```bash
# RPC 服务（被 API 依赖，需先启动）
go run fim_user/user_rpc/userrpc.go         -f fim_user/user_rpc/etc/userrpc.yaml
go run fim_chat/chat_rpc/chatrpc.go         -f fim_chat/chat_rpc/etc/chatrpc.yaml
go run fim_group/group_rpc/grouprpc.go      -f fim_group/group_rpc/etc/grouprpc.yaml
go run fim_file/file_rpc/filerpc.go         -f fim_file/file_rpc/etc/filerpc.yaml
go run fim_settings/settings_rpc/settingsrpc.go -f fim_settings/settings_rpc/etc/settingsrpc.yaml

# API 服务
go run fim_auth/auth_api/auth.go            -f fim_auth/auth_api/etc/auth.yaml
go run fim_user/user_api/users.go           -f fim_user/user_api/etc/users.yaml
go run fim_chat/chat_api/chat.go            -f fim_chat/chat_api/etc/chat.yaml
go run fim_group/group_api/group.go         -f fim_group/group_api/etc/group.yaml
go run fim_file/file_api/file.go            -f fim_file/file_api/etc/file.yaml
go run fim_settings/settings_api/settings.go -f fim_settings/settings_api/etc/settings.yaml
go run fim_logs/logs_api/logs.go            -f fim_logs/logs_api/etc/logs.yaml

# 网关（统一入口 :8080）
go run fim_gateway/gateway.go               -f fim_gateway/settings.yaml
```

### 5. 启动前端

```bash
cd frontend
npm install
npm run dev
```

修改 `frontend/.env` 中 `VITE_SERVER_URL` 为后端网关地址，例如：

```
VITE_SERVER_URL=http://127.0.0.1:8080
```

`vite.config.ts` 已配置 `/api`、`/api/chat/ws`、`/api/group/ws` 代理到该地址。开发服务器默认监听 **80** 端口（如被占用，可在 `vite.config.ts` 修改 `server.port`），访问：

```
http://127.0.0.1
```

---

## 二、生产部署（前端编译 + 后端 Docker）

整套生产环境通过 `backend/deploy/docker-compose.yaml` 编排，包含全部依赖、所有微服务以及 nginx 静态资源容器。

### 1. 编译前端

```bash
cd frontend

# 设置生产环境后端地址（按实际服务器/域名修改）
echo "VITE_SERVER_URL=https://your-domain" > .env

npm install
npm run build
```

构建产物位于 `frontend/dist/`。将其复制到部署目录：

```bash
# 在仓库根目录执行
cp -r frontend/dist backend/deploy/fim_web/
```

`backend/deploy/docker-compose.yaml` 中的 `fim_web` 服务会把该目录挂载到 nginx 容器的 `/usr/share/nginx/fim_web`。

### 2. 准备后端配置 / 证书

- 检查 `backend/deploy/<service>/<api|rpc>/*.yaml`：MySQL / Redis / etcd / kafka 默认使用容器内网 IP（`10.0.0.20`、`10.0.0.21`、`10.0.0.22`、`10.0.0.24`），通常无需修改
- 修改 `backend/deploy/nginx/nginx.conf` 中的 `server_name` 为自己的域名
- 将 SSL 证书放入 `backend/deploy/nginx/cert/`，并同步修改 nginx.conf 中 `ssl_certificate` / `ssl_certificate_key` 路径
- 文件上传持久化目录：`backend/deploy/fim_file/file_api/uploads/`（可按需挂卷）

### 3. 构建 Go 基础镜像

`backend/deploy/<service>/Dockerfile` 都是 `FROM fim_server`，因此必须先用仓库根 `backend/Dockerfile` 构建出统一编译产物的基础镜像：

```bash
cd backend
docker build -t fim_server .
```

该 Dockerfile 在 `golang:alpine` 中一次编译出所有服务二进制（auth/chat/file/gateway/group/logs/settings/user 的 api 与 rpc）。

### 4. 启动整个后端集群

```bash
cd backend/deploy
docker compose up -d
```

会按依赖顺序拉起：

```
mysql (10.0.0.20)  redis (10.0.0.21)  etcd (10.0.0.22)
zookeeper (10.0.0.23)  kafka (10.0.0.24)  kafka-map (10.0.0.25, :9001)
gateway (10.0.0.2, :8080)
auth_api / chat_api+chat_rpc / file_api+file_rpc /
group_api+group_rpc / logs_api / settings_api+settings_rpc /
user_api+user_rpc
fim_web (nginx, :80 + :443)
```

### 5. 验证

- 浏览器访问 `https://your-domain` 进入聊天界面
- 网关健康端口：`http://server:8080/api/...`
- Kafka 监控面板：`http://server:9001`（账号 / 密码：admin / admin）
- 日志：`backend/deploy/nginx/logs/`、容器 `docker logs <container>`

### 常用运维命令

```bash
# 查看所有容器状态
cd backend/deploy && docker compose ps

# 查看某个服务日志
docker compose logs -f gateway

# 重启某个微服务（修改配置后）
docker compose restart user_api

# 全量重新构建
docker compose build --no-cache && docker compose up -d
```

---

## 三、模块说明

| 服务            | 职责                                             | 端口（容器）   |
|-----------------|--------------------------------------------------|----------------|
| fim_gateway     | 统一 HTTP 网关，分发 REST 与 WebSocket           | 8080           |
| fim_auth        | 登录、Token 颁发与校验                           | 内部           |
| fim_user        | 用户、好友、好友验证、用户配置（api + rpc）      | 内部           |
| fim_chat        | 单聊消息、置顶、消息删除（api + rpc，含 ws）     | 内部           |
| fim_group       | 群聊、群成员、群验证、群消息（api + rpc，含 ws） | 内部           |
| fim_file        | 头像 / 群头像 / 普通文件上传下载（api + rpc）    | 内部           |
| fim_settings    | 系统设置（api + rpc）                            | 内部           |
| fim_logs        | 通过 kafka 异步收集日志                          | 内部           |

## 四、常见问题

- **MySQL 端口冲突**：`develop/docker-compose.yaml` 暴露在 3307；如本机有 MySQL 占用 3306 即可避免冲突
- **etcd 服务发现失败**：检查每个服务 `etc/*.yaml` 中的 etcd 地址是否一致并可达
- **前端 WebSocket 连接失败**：确认网关已启动；生产环境检查 nginx.conf 中 `/api/chat/ws/chat`、`/api/group/ws/chat` 的 `Upgrade` 头转发
- **生产环境镜像构建失败**：`deploy/<service>/Dockerfile` 都依赖本地已存在的 `fim_server` 镜像，必须先在 `backend/` 下执行 `docker build -t fim_server .`

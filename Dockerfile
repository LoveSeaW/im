FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
COPY . .
RUN go mod tidy

RUN go build -o im_auth/auth_api/auth im_auth/auth_api/auth.go

RUN go build -o im_chat/chat_api/chat im_chat/chat_api/chat.go
RUN go build -o im_chat/chat_rpc/chatrpc im_chat/chat_rpc/chatrpc.go

RUN go build -o fim_file/file_api/file fim_file/file_api/file.go
RUN go build -o fim_file/file_rpc/filerpc fim_file/file_rpc/filerpc.go

RUN go build -o fim_gateway/gateway fim_gateway/gateway.go

RUN go build -o fim_group/group_api/group fim_group/group_api/group.go
RUN go build -o fim_group/group_rpc/grouprpc fim_group/group_rpc/grouprpc.go

RUN go build -o im_logs/logs_api/logs im_logs/logs_api/logs.go

RUN go build -o im_settings/settings_api/settings im_settings/settings_api/settings.go
RUN go build -o im_settings/settings_rpc/settingsrpc im_settings/settings_rpc/settingsrpc.go

RUN go build -o im_user/user_api/users im_user/user_api/users.go
RUN go build -o im_user/user_rpc/userrpc im_user/user_rpc/userrpc.go
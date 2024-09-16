// Code generated by goctl. DO NOT EDIT.
// Source: chat_rpc.proto

package chat

import (
	"context"

	"im_server/im_chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	ChatCountMessage          = chat_rpc.ChatCountMessage
	UserChatRequest           = chat_rpc.UserChatRequest
	UserChatResponse          = chat_rpc.UserChatResponse
	UserListChatCountRequest  = chat_rpc.UserListChatCountRequest
	UserListChatCountResponse = chat_rpc.UserListChatCountResponse

	Chat interface {
		UserChat(ctx context.Context, in *UserChatRequest, opts ...grpc.CallOption) (*UserChatResponse, error)
		UserListChatCount(ctx context.Context, in *UserListChatCountRequest, opts ...grpc.CallOption) (*UserListChatCountResponse, error)
	}

	defaultChat struct {
		cli zrpc.Client
	}
)

func NewChat(cli zrpc.Client) Chat {
	return &defaultChat{
		cli: cli,
	}
}

func (m *defaultChat) UserChat(ctx context.Context, in *UserChatRequest, opts ...grpc.CallOption) (*UserChatResponse, error) {
	client := chat_rpc.NewChatClient(m.cli.Conn())
	return client.UserChat(ctx, in, opts...)
}

func (m *defaultChat) UserListChatCount(ctx context.Context, in *UserListChatCountRequest, opts ...grpc.CallOption) (*UserListChatCountResponse, error) {
	client := chat_rpc.NewChatClient(m.cli.Conn())
	return client.UserListChatCount(ctx, in, opts...)
}

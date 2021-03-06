// Code generated by goctl. DO NOT EDIT!
// Source: trace.proto

package traceclient

import (
	"context"

	"web_game/service/trace/rpc/trace"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	ReqTrace = trace.ReqTrace
	ResTrace = trace.ResTrace

	Trace interface {
		PushTrace(ctx context.Context, in *ReqTrace, opts ...grpc.CallOption) (*ResTrace, error)
	}

	defaultTrace struct {
		cli zrpc.Client
	}
)

func NewTrace(cli zrpc.Client) Trace {
	return &defaultTrace{
		cli: cli,
	}
}

func (m *defaultTrace) PushTrace(ctx context.Context, in *ReqTrace, opts ...grpc.CallOption) (*ResTrace, error) {
	client := trace.NewTraceClient(m.cli.Conn())
	return client.PushTrace(ctx, in, opts...)
}

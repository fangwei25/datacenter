package logic

import (
	"context"

	"web_game/service/trace/rpc/internal/svc"
	"web_game/service/trace/rpc/trace"

	"github.com/zeromicro/go-zero/core/logx"
)

type PushTraceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPushTraceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushTraceLogic {
	return &PushTraceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PushTraceLogic) PushTrace(in *trace.ReqTrace) (*trace.ResTrace, error) {
	//本地文件写入
	err := l.svcCtx.TraceFile.WriteTrace(in.TraceName, in.JsonData)
	if nil != err {
		logx.Error("record trace file failed, ", err)
	}

	//todo 推送EA

	return &trace.ResTrace{}, nil
}

package logic

import (
	"context"
	"strconv"
	"time"
	"web_game/common/counterx"

	"web_game/service/counter/rpc/counter"
	"web_game/service/counter/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryLogic {
	return &QueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryLogic) Query(in *counter.ReqQuery) (*counter.ResQuery, error) {
	key := GenKey(in.OwnerId, in.EventType, counterx.TimeDimension(in.TimeDimension), time.Unix(in.TimeFlag, 0))
	fieldCount := GenField(in.EventField, counterx.CTCount)
	fieldValue := GenField(in.EventField, counterx.CTValue)
	fieldMax := GenField(in.EventField, counterx.CTMax)
	fieldMin := GenField(in.EventField, counterx.CTMin)
	results, err := l.HMGetInt64(key, fieldCount, fieldValue, fieldMax, fieldMin)
	if nil != err {
		return nil, err
	}

	return &counter.ResQuery{
		Count: results[fieldCount],
		Value: results[fieldValue],
		Max:   results[fieldMax],
		Min:   results[fieldMin],
	}, nil
}

func (l *QueryLogic) HMGetInt64(key string, fields ...string) (map[string]int64, error) {
	res, err := l.svcCtx.RedisObj.Hmget(key, fields...)
	if nil != err {
		return nil, err
	}
	resMap := make(map[string]int64)
	for idx, oneRes := range res {
		if len(oneRes) == 0 {
			resMap[fields[idx]] = 0
		} else {
			num, err := strconv.ParseInt(oneRes, 10, 64)
			if err == nil {
				resMap[fields[idx]] = num
			}
		}
	}
	return resMap, nil
}

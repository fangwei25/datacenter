package logic

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
	"web_game/common/counterx"

	"web_game/service/counter/rpc/counter"
	"web_game/service/counter/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLogic) Update(in *counter.ReqUpdate) (*counter.ResUpdate, error) {
	l.doUpdate(in.OwnerId, in.EventType, in.EventField, in.Value)
	return &counter.ResUpdate{}, nil
}

// doUpdate 数值统计 按照给定的数值进行累加统计
func (l *UpdateLogic) doUpdate(ownerId int64, eventType string, field string, value int64) {
	//根据配置决定统计的计算维度（累加，计数，记录最大值，记录最小值）和时间维度(总，年，月，日，小时，分)
	nonce := time.Now()
	for _, TD := range counterx.TDList {
		key := GenKey(ownerId, eventType, TD, nonce)
		lifeTime := GetLifeTime(TD, 5)
		l.UpdateByTimeDimension(key, field, value, lifeTime, counterx.CalcList)
	}
}

func (l *UpdateLogic) UpdateByTimeDimension(key, field string, value int64, expire time.Duration, calcTypes []counterx.CalcType) {
	var err error
	for _, calcType := range calcTypes {
		fieldExt := GenField(field, calcType)
		switch calcType {
		case counterx.CTCount:
			_, err = l.svcCtx.RedisObj.Hincrby(key, fieldExt, 1)
		case counterx.CTValue:
			_, err = l.svcCtx.RedisObj.Hincrby(key, fieldExt, int(value))
		case counterx.CTMax:
			_, err = l.UpdateMax(key, fieldExt, value)
		case counterx.CTMin:
			_, err = l.UpdateMinus(key, fieldExt, value)
		}
		if err != nil {
			fmt.Printf("Engine.UpdateByTimeDimension failed, key=%s, field=%s, value=%d, err: %v", key, field, value, err)
		}
	}
	if expire != -1 {
		_ = l.svcCtx.RedisObj.Expire(key, int(expire/time.Second))
	}
}

// UpdateMax 更新最大值
func (l *UpdateLogic) UpdateMax(key, field string, v int64) (newV int64, err error) {
	var result string
	result, err = l.svcCtx.RedisObj.Hget(key, field)
	if nil != err && err != redis.Nil {
		return //查询失败，返回0，不限制生成，避免业务不能持续
	}

	var queryValue int64
	if result == "" {
		queryValue = 0
	} else {
		queryValue, err = strconv.ParseInt(result, 10, 64)
		if nil != err {
			return
		}
	}

	if queryValue >= v {
		return queryValue, nil
	}
	err = l.svcCtx.RedisObj.Hset(key, field, strconv.FormatInt(v, 10))
	return
}

// UpdateMinus 更新最小值
func (l *UpdateLogic) UpdateMinus(key, field string, v int64) (newV int64, err error) {
	var result string
	result, err = l.svcCtx.RedisObj.Hget(key, field)
	if nil != err && err != redis.Nil {
		return //查询失败，返回0，不限制生成，避免业务不能持续
	}

	var queryValue int64
	if result == "" {
		queryValue = 0
	} else {
		queryValue, err = strconv.ParseInt(result, 10, 64)
		if nil != err {
			return
		}
	}

	if queryValue <= v {
		return queryValue, nil
	}
	err = l.svcCtx.RedisObj.Hset(key, field, strconv.FormatInt(v, 10))
	return
}

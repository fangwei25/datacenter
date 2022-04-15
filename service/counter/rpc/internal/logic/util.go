package logic

import (
	"fmt"
	"time"
	"web_game/common/counterx"
)

const KeyPatton = "mdx-%d-%s-%s" //mdx-ownerId-eventType-time(根据time的格式可以识别时间维度)
const totalCounterField = "total"

const (
	timeFormatEver         = "0"
	timeFormatPattonYear   = "2006"
	timeFormatPattonMonth  = "2006-01"
	timeFormatPattonDay    = "2006-01-02"
	timeFormatPattonHour   = "2006-01-02-15"
	timeFormatPattonMinute = "2006-01-02-15-04"
)

var TFPMap map[counterx.TimeDimension]string

func init() {
	TFPMap = make(map[counterx.TimeDimension]string)
	//TFPMap[cfg.TDEver] = timeFormatEver
	TFPMap[counterx.TDYear] = timeFormatPattonYear
	TFPMap[counterx.TDMonth] = timeFormatPattonMonth
	TFPMap[counterx.TDDay] = timeFormatPattonDay
	TFPMap[counterx.TDHour] = timeFormatPattonHour
	TFPMap[counterx.TDMinute] = timeFormatPattonMinute
}

func GenKey(ownerId int64, eventType string, timeDimension counterx.TimeDimension, t time.Time) string {
	var key string
	if timeDimension == counterx.TDEver {
		key = fmt.Sprintf(KeyPatton, ownerId, eventType, timeFormatEver)
	} else {
		tfp := TFPMap[timeDimension]
		key = fmt.Sprintf(KeyPatton, ownerId, eventType, t.Format(tfp))
	}
	return key
}

func GenField(field string, calcType counterx.CalcType) string {
	return field + "-" + string(calcType)
}

func GetLifeTime(timeDimension counterx.TimeDimension, cfgValue int32) time.Duration {
	switch timeDimension {
	case counterx.TDEver:
		return -1
	case counterx.TDYear:
		return time.Duration(cfgValue) * time.Hour * 24 * 365
	case counterx.TDMonth:
		return time.Duration(cfgValue) * time.Hour * 24 * 30
	case counterx.TDDay:
		return time.Duration(cfgValue) * time.Hour * 24
	case counterx.TDHour:
		return time.Duration(cfgValue) * time.Hour
	case counterx.TDMinute:
		return time.Duration(cfgValue) * time.Minute
	}
	fmt.Printf("GetLifeTime failed, no timeDimension hit: %d", timeDimension)
	return -1
}

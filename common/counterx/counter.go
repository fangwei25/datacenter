package counterx

// TimeDimension 时间维度类型
type TimeDimension int

const (
	TDEver = TimeDimension(iota)
	TDYear
	TDMonth
	TDDay
	TDHour
	TDMinute
)

var TDList []TimeDimension

// CalcType 计算类型
type CalcType string

const (
	CTCount = CalcType("count")
	CTValue = CalcType("value")
	CTMax   = CalcType("max")
	CTMin   = CalcType("min")
)

var CalcList []CalcType

func init() {
	TDList = []TimeDimension{TDEver, TDYear, TDMonth, TDDay, TDHour, TDMinute}
	CalcList = []CalcType{CTCount, CTValue, CTMax, CTMin}
}

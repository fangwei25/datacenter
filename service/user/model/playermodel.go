package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	playerFieldNames          = builder.RawFieldNames(&Player{})
	playerRows                = strings.Join(playerFieldNames, ",")
	playerRowsExpectAutoSet   = strings.Join(stringx.Remove(playerFieldNames, "`create_time`", "`update_time`"), ",")
	playerRowsWithPlaceHolder = strings.Join(stringx.Remove(playerFieldNames, "`player_id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cachePlayerPlayerIdPrefix = "cache:player:playerId:"
)

type (
	PlayerModel interface {
		Insert(data *Player) (sql.Result, error)
		FindOne(playerId int64) (*Player, error)
		Update(data *Player) error
		Delete(playerId int64) error
	}

	defaultPlayerModel struct {
		sqlc.CachedConn
		table string
	}

	Player struct {
		PlayerId       int64     `db:"player_id"`
		Name           string    `db:"name"`
		Gender         int64     `db:"gender"`           // 性别
		AvatorUrl      string    `db:"avator_url"`       // 头像地址
		InvitationId   int64     `db:"invitation_id"`    // 邀请玩家id
		Channel        string    `db:"channel"`          // 渠道来源
		VipLv          int64     `db:"vip_lv"`           // vip等级
		VipExp         int64     `db:"vip_exp"`          // vip升级经验
		Level          int64     `db:"level"`            // 等级
		LevelExp       int64     `db:"level_exp"`        // 等级经验
		OriAccountType string    `db:"ori_account_type"` // 玩家原始注册方式
		OriAccount     string    `db:"ori_account"`      // 玩家原始注册账户名
		CreateTime     time.Time `db:"create_time"`
		UpdateTime     time.Time `db:"update_time"`
	}
)

func NewPlayerModel(conn sqlx.SqlConn, c cache.CacheConf) PlayerModel {
	return &defaultPlayerModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`player`",
	}
}

func (m *defaultPlayerModel) Insert(data *Player) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, playerRowsExpectAutoSet)
	ret, err := m.ExecNoCache(query, data.PlayerId, data.Name, data.Gender, data.AvatorUrl, data.InvitationId, data.Channel, data.VipLv, data.VipExp, data.Level, data.LevelExp, data.OriAccountType, data.OriAccount)

	return ret, err
}

func (m *defaultPlayerModel) FindOne(playerId int64) (*Player, error) {
	playerPlayerIdKey := fmt.Sprintf("%s%v", cachePlayerPlayerIdPrefix, playerId)
	var resp Player
	err := m.QueryRow(&resp, playerPlayerIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `player_id` = ? limit 1", playerRows, m.table)
		return conn.QueryRow(v, query, playerId)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultPlayerModel) Update(data *Player) error {
	playerPlayerIdKey := fmt.Sprintf("%s%v", cachePlayerPlayerIdPrefix, data.PlayerId)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `player_id` = ?", m.table, playerRowsWithPlaceHolder)
		return conn.Exec(query, data.Name, data.Gender, data.AvatorUrl, data.InvitationId, data.Channel, data.VipLv, data.VipExp, data.Level, data.LevelExp, data.OriAccountType, data.OriAccount, data.PlayerId)
	}, playerPlayerIdKey)
	return err
}

func (m *defaultPlayerModel) Delete(playerId int64) error {

	playerPlayerIdKey := fmt.Sprintf("%s%v", cachePlayerPlayerIdPrefix, playerId)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `player_id` = ?", m.table)
		return conn.Exec(query, playerId)
	}, playerPlayerIdKey)
	return err
}

func (m *defaultPlayerModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cachePlayerPlayerIdPrefix, primary)
}

func (m *defaultPlayerModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `player_id` = ? limit 1", playerRows, m.table)
	return conn.QueryRow(v, query, primary)
}

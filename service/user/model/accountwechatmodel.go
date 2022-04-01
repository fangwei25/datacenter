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
	accountWechatFieldNames          = builder.RawFieldNames(&AccountWechat{})
	accountWechatRows                = strings.Join(accountWechatFieldNames, ",")
	accountWechatRowsExpectAutoSet   = strings.Join(stringx.Remove(accountWechatFieldNames, "`create_time`", "`update_time`"), ",")
	accountWechatRowsWithPlaceHolder = strings.Join(stringx.Remove(accountWechatFieldNames, "`account`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheAccountWechatAccountPrefix  = "cache:accountWechat:account:"
	cacheAccountWechatPlayerIdPrefix = "cache:accountWechat:playerId:"
)

type (
	AccountWechatModel interface {
		Insert(data *AccountWechat) (sql.Result, error)
		FindOne(account string) (*AccountWechat, error)
		FindOneByPlayerId(playerId int64) (*AccountWechat, error)
		Update(data *AccountWechat) error
		Delete(account string) error
	}

	defaultAccountWechatModel struct {
		sqlc.CachedConn
		table string
	}

	AccountWechat struct {
		Account     string    `db:"account"`
		PlayerId    int64     `db:"player_id"`
		UnionId     string    `db:"union_id"`
		AccessToken string    `db:"access_token"`
		UserInfo    string    `db:"user_info"`
		CreateTime  time.Time `db:"create_time"`
		LastLogin   time.Time `db:"last_login"`
	}
)

func NewAccountWechatModel(conn sqlx.SqlConn, c cache.CacheConf) AccountWechatModel {
	return &defaultAccountWechatModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`account_wechat`",
	}
}

func (m *defaultAccountWechatModel) Insert(data *AccountWechat) (sql.Result, error) {
	accountWechatAccountKey := fmt.Sprintf("%s%v", cacheAccountWechatAccountPrefix, data.Account)
	accountWechatPlayerIdKey := fmt.Sprintf("%s%v", cacheAccountWechatPlayerIdPrefix, data.PlayerId)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, accountWechatRowsExpectAutoSet)
		return conn.Exec(query, data.Account, data.PlayerId, data.UnionId, data.AccessToken, data.UserInfo, data.LastLogin)
	}, accountWechatAccountKey, accountWechatPlayerIdKey)
	return ret, err
}

func (m *defaultAccountWechatModel) FindOne(account string) (*AccountWechat, error) {
	accountWechatAccountKey := fmt.Sprintf("%s%v", cacheAccountWechatAccountPrefix, account)
	var resp AccountWechat
	err := m.QueryRow(&resp, accountWechatAccountKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `account` = ? limit 1", accountWechatRows, m.table)
		return conn.QueryRow(v, query, account)
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

func (m *defaultAccountWechatModel) FindOneByPlayerId(playerId int64) (*AccountWechat, error) {
	accountWechatPlayerIdKey := fmt.Sprintf("%s%v", cacheAccountWechatPlayerIdPrefix, playerId)
	var resp AccountWechat
	err := m.QueryRowIndex(&resp, accountWechatPlayerIdKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `player_id` = ? limit 1", accountWechatRows, m.table)
		if err := conn.QueryRow(&resp, query, playerId); err != nil {
			return nil, err
		}
		return resp.Account, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultAccountWechatModel) Update(data *AccountWechat) error {
	accountWechatAccountKey := fmt.Sprintf("%s%v", cacheAccountWechatAccountPrefix, data.Account)
	accountWechatPlayerIdKey := fmt.Sprintf("%s%v", cacheAccountWechatPlayerIdPrefix, data.PlayerId)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `account` = ?", m.table, accountWechatRowsWithPlaceHolder)
		return conn.Exec(query, data.PlayerId, data.UnionId, data.AccessToken, data.UserInfo, data.LastLogin, data.Account)
	}, accountWechatAccountKey, accountWechatPlayerIdKey)
	return err
}

func (m *defaultAccountWechatModel) Delete(account string) error {
	data, err := m.FindOne(account)
	if err != nil {
		return err
	}

	accountWechatAccountKey := fmt.Sprintf("%s%v", cacheAccountWechatAccountPrefix, account)
	accountWechatPlayerIdKey := fmt.Sprintf("%s%v", cacheAccountWechatPlayerIdPrefix, data.PlayerId)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `account` = ?", m.table)
		return conn.Exec(query, account)
	}, accountWechatAccountKey, accountWechatPlayerIdKey)
	return err
}

func (m *defaultAccountWechatModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheAccountWechatAccountPrefix, primary)
}

func (m *defaultAccountWechatModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `account` = ? limit 1", accountWechatRows, m.table)
	return conn.QueryRow(v, query, primary)
}

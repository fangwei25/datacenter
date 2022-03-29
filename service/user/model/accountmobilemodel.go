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
	accountMobileFieldNames          = builder.RawFieldNames(&AccountMobile{})
	accountMobileRows                = strings.Join(accountMobileFieldNames, ",")
	accountMobileRowsExpectAutoSet   = strings.Join(stringx.Remove(accountMobileFieldNames, "`create_time`", "`update_time`"), ",")
	accountMobileRowsWithPlaceHolder = strings.Join(stringx.Remove(accountMobileFieldNames, "`account`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheAccountMobileAccountPrefix  = "cache:accountMobile:account:"
	cacheAccountMobilePlayerIdPrefix = "cache:accountMobile:playerId:"
)

type (
	AccountMobileModel interface {
		Insert(data *AccountMobile) (sql.Result, error)
		FindOne(account string) (*AccountMobile, error)
		FindOneByPlayerId(playerId int64) (*AccountMobile, error)
		Update(data *AccountMobile) error
		Delete(account string) error
	}

	defaultAccountMobileModel struct {
		sqlc.CachedConn
		table string
	}

	AccountMobile struct {
		Account    string    `db:"account"`
		PlayerId   int64     `db:"player_id"`
		CreateTime time.Time `db:"create_time"`
		LastLogin  time.Time `db:"last_login"`
	}
)

func NewAccountMobileModel(conn sqlx.SqlConn, c cache.CacheConf) AccountMobileModel {
	return &defaultAccountMobileModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`account_mobile`",
	}
}

func (m *defaultAccountMobileModel) Insert(data *AccountMobile) (sql.Result, error) {
	accountMobileAccountKey := fmt.Sprintf("%s%v", cacheAccountMobileAccountPrefix, data.Account)
	accountMobilePlayerIdKey := fmt.Sprintf("%s%v", cacheAccountMobilePlayerIdPrefix, data.PlayerId)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?)", m.table, accountMobileRowsExpectAutoSet)
		return conn.Exec(query, data.Account, data.PlayerId, data.LastLogin)
	}, accountMobileAccountKey, accountMobilePlayerIdKey)
	return ret, err
}

func (m *defaultAccountMobileModel) FindOne(account string) (*AccountMobile, error) {
	accountMobileAccountKey := fmt.Sprintf("%s%v", cacheAccountMobileAccountPrefix, account)
	var resp AccountMobile
	err := m.QueryRow(&resp, accountMobileAccountKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `account` = ? limit 1", accountMobileRows, m.table)
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

func (m *defaultAccountMobileModel) FindOneByPlayerId(playerId int64) (*AccountMobile, error) {
	accountMobilePlayerIdKey := fmt.Sprintf("%s%v", cacheAccountMobilePlayerIdPrefix, playerId)
	var resp AccountMobile
	err := m.QueryRowIndex(&resp, accountMobilePlayerIdKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `player_id` = ? limit 1", accountMobileRows, m.table)
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

func (m *defaultAccountMobileModel) Update(data *AccountMobile) error {
	accountMobilePlayerIdKey := fmt.Sprintf("%s%v", cacheAccountMobilePlayerIdPrefix, data.PlayerId)
	accountMobileAccountKey := fmt.Sprintf("%s%v", cacheAccountMobileAccountPrefix, data.Account)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `account` = ?", m.table, accountMobileRowsWithPlaceHolder)
		return conn.Exec(query, data.PlayerId, data.LastLogin, data.Account)
	}, accountMobileAccountKey, accountMobilePlayerIdKey)
	return err
}

func (m *defaultAccountMobileModel) Delete(account string) error {
	data, err := m.FindOne(account)
	if err != nil {
		return err
	}

	accountMobileAccountKey := fmt.Sprintf("%s%v", cacheAccountMobileAccountPrefix, account)
	accountMobilePlayerIdKey := fmt.Sprintf("%s%v", cacheAccountMobilePlayerIdPrefix, data.PlayerId)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `account` = ?", m.table)
		return conn.Exec(query, account)
	}, accountMobileAccountKey, accountMobilePlayerIdKey)
	return err
}

func (m *defaultAccountMobileModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheAccountMobileAccountPrefix, primary)
}

func (m *defaultAccountMobileModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `account` = ? limit 1", accountMobileRows, m.table)
	return conn.QueryRow(v, query, primary)
}

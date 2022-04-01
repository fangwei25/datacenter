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
	accountPwdFieldNames          = builder.RawFieldNames(&AccountPwd{})
	accountPwdRows                = strings.Join(accountPwdFieldNames, ",")
	accountPwdRowsExpectAutoSet   = strings.Join(stringx.Remove(accountPwdFieldNames, "`create_time`", "`update_time`"), ",")
	accountPwdRowsWithPlaceHolder = strings.Join(stringx.Remove(accountPwdFieldNames, "`account`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheAccountPwdAccountPrefix  = "cache:accountPwd:account:"
	cacheAccountPwdPlayerIdPrefix = "cache:accountPwd:playerId:"
)

type (
	AccountPwdModel interface {
		Insert(data *AccountPwd) (sql.Result, error)
		FindOne(account string) (*AccountPwd, error)
		FindOneByPlayerId(playerId int64) (*AccountPwd, error)
		Update(data *AccountPwd) error
		Delete(account string) error
	}

	defaultAccountPwdModel struct {
		sqlc.CachedConn
		table string
	}

	AccountPwd struct {
		Account    string    `db:"account"`
		PlayerId   int64     `db:"player_id"`
		Pwd        string    `db:"pwd"`
		CreateTime time.Time `db:"create_time"`
		LastLogin  time.Time `db:"last_login"`
	}
)

func NewAccountPwdModel(conn sqlx.SqlConn, c cache.CacheConf) AccountPwdModel {
	return &defaultAccountPwdModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`account_pwd`",
	}
}

func (m *defaultAccountPwdModel) Insert(data *AccountPwd) (sql.Result, error) {
	accountPwdAccountKey := fmt.Sprintf("%s%v", cacheAccountPwdAccountPrefix, data.Account)
	accountPwdPlayerIdKey := fmt.Sprintf("%s%v", cacheAccountPwdPlayerIdPrefix, data.PlayerId)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, accountPwdRowsExpectAutoSet)
		return conn.Exec(query, data.Account, data.PlayerId, data.Pwd, data.LastLogin)
	}, accountPwdAccountKey, accountPwdPlayerIdKey)
	return ret, err
}

func (m *defaultAccountPwdModel) FindOne(account string) (*AccountPwd, error) {
	accountPwdAccountKey := fmt.Sprintf("%s%v", cacheAccountPwdAccountPrefix, account)
	var resp AccountPwd
	err := m.QueryRow(&resp, accountPwdAccountKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `account` = ? limit 1", accountPwdRows, m.table)
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

func (m *defaultAccountPwdModel) FindOneByPlayerId(playerId int64) (*AccountPwd, error) {
	accountPwdPlayerIdKey := fmt.Sprintf("%s%v", cacheAccountPwdPlayerIdPrefix, playerId)
	var resp AccountPwd
	err := m.QueryRowIndex(&resp, accountPwdPlayerIdKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `player_id` = ? limit 1", accountPwdRows, m.table)
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

func (m *defaultAccountPwdModel) Update(data *AccountPwd) error {
	accountPwdAccountKey := fmt.Sprintf("%s%v", cacheAccountPwdAccountPrefix, data.Account)
	accountPwdPlayerIdKey := fmt.Sprintf("%s%v", cacheAccountPwdPlayerIdPrefix, data.PlayerId)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `account` = ?", m.table, accountPwdRowsWithPlaceHolder)
		return conn.Exec(query, data.PlayerId, data.Pwd, data.LastLogin, data.Account)
	}, accountPwdAccountKey, accountPwdPlayerIdKey)
	return err
}

func (m *defaultAccountPwdModel) Delete(account string) error {
	data, err := m.FindOne(account)
	if err != nil {
		return err
	}

	accountPwdAccountKey := fmt.Sprintf("%s%v", cacheAccountPwdAccountPrefix, account)
	accountPwdPlayerIdKey := fmt.Sprintf("%s%v", cacheAccountPwdPlayerIdPrefix, data.PlayerId)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `account` = ?", m.table)
		return conn.Exec(query, account)
	}, accountPwdAccountKey, accountPwdPlayerIdKey)
	return err
}

func (m *defaultAccountPwdModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheAccountPwdAccountPrefix, primary)
}

func (m *defaultAccountPwdModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `account` = ? limit 1", accountPwdRows, m.table)
	return conn.QueryRow(v, query, primary)
}

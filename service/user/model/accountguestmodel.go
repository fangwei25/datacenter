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
	accountGuestFieldNames          = builder.RawFieldNames(&AccountGuest{})
	accountGuestRows                = strings.Join(accountGuestFieldNames, ",")
	accountGuestRowsExpectAutoSet   = strings.Join(stringx.Remove(accountGuestFieldNames, "`create_time`", "`update_time`"), ",")
	accountGuestRowsWithPlaceHolder = strings.Join(stringx.Remove(accountGuestFieldNames, "`account`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheAccountGuestAccountPrefix  = "cache:accountGuest:account:"
	cacheAccountGuestPlayerIdPrefix = "cache:accountGuest:playerId:"
)

type (
	AccountGuestModel interface {
		Insert(data *AccountGuest) (sql.Result, error)
		FindOne(account string) (*AccountGuest, error)
		FindOneByPlayerId(playerId int64) (*AccountGuest, error)
		Update(data *AccountGuest) error
		Delete(account string) error
	}

	defaultAccountGuestModel struct {
		sqlc.CachedConn
		table string
	}

	AccountGuest struct {
		Account     string         `db:"account"`
		PlayerId    int64          `db:"player_id"`
		AccessToken sql.NullString `db:"access_token"`
		CreateTime  time.Time      `db:"create_time"`
		LastLogin   time.Time      `db:"last_login"`
	}
)

func NewAccountGuestModel(conn sqlx.SqlConn, c cache.CacheConf) AccountGuestModel {
	return &defaultAccountGuestModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`account_guest`",
	}
}

func (m *defaultAccountGuestModel) Insert(data *AccountGuest) (sql.Result, error) {
	accountGuestAccountKey := fmt.Sprintf("%s%v", cacheAccountGuestAccountPrefix, data.Account)
	accountGuestPlayerIdKey := fmt.Sprintf("%s%v", cacheAccountGuestPlayerIdPrefix, data.PlayerId)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, accountGuestRowsExpectAutoSet)
		return conn.Exec(query, data.Account, data.PlayerId, data.AccessToken, data.LastLogin)
	}, accountGuestAccountKey, accountGuestPlayerIdKey)
	return ret, err
}

func (m *defaultAccountGuestModel) FindOne(account string) (*AccountGuest, error) {
	accountGuestAccountKey := fmt.Sprintf("%s%v", cacheAccountGuestAccountPrefix, account)
	var resp AccountGuest
	err := m.QueryRow(&resp, accountGuestAccountKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `account` = ? limit 1", accountGuestRows, m.table)
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

func (m *defaultAccountGuestModel) FindOneByPlayerId(playerId int64) (*AccountGuest, error) {
	accountGuestPlayerIdKey := fmt.Sprintf("%s%v", cacheAccountGuestPlayerIdPrefix, playerId)
	var resp AccountGuest
	err := m.QueryRowIndex(&resp, accountGuestPlayerIdKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `player_id` = ? limit 1", accountGuestRows, m.table)
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

func (m *defaultAccountGuestModel) Update(data *AccountGuest) error {
	accountGuestAccountKey := fmt.Sprintf("%s%v", cacheAccountGuestAccountPrefix, data.Account)
	accountGuestPlayerIdKey := fmt.Sprintf("%s%v", cacheAccountGuestPlayerIdPrefix, data.PlayerId)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `account` = ?", m.table, accountGuestRowsWithPlaceHolder)
		return conn.Exec(query, data.PlayerId, data.AccessToken, data.LastLogin, data.Account)
	}, accountGuestAccountKey, accountGuestPlayerIdKey)
	return err
}

func (m *defaultAccountGuestModel) Delete(account string) error {
	data, err := m.FindOne(account)
	if err != nil {
		return err
	}

	accountGuestAccountKey := fmt.Sprintf("%s%v", cacheAccountGuestAccountPrefix, account)
	accountGuestPlayerIdKey := fmt.Sprintf("%s%v", cacheAccountGuestPlayerIdPrefix, data.PlayerId)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `account` = ?", m.table)
		return conn.Exec(query, account)
	}, accountGuestAccountKey, accountGuestPlayerIdKey)
	return err
}

func (m *defaultAccountGuestModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheAccountGuestAccountPrefix, primary)
}

func (m *defaultAccountGuestModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `account` = ? limit 1", accountGuestRows, m.table)
	return conn.QueryRow(v, query, primary)
}

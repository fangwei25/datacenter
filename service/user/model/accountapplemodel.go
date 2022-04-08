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
	accountAppleFieldNames          = builder.RawFieldNames(&AccountApple{})
	accountAppleRows                = strings.Join(accountAppleFieldNames, ",")
	accountAppleRowsExpectAutoSet   = strings.Join(stringx.Remove(accountAppleFieldNames, "`create_time`", "`update_time`"), ",")
	accountAppleRowsWithPlaceHolder = strings.Join(stringx.Remove(accountAppleFieldNames, "`account`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheAccountAppleAccountPrefix  = "cache:accountApple:account:"
	cacheAccountApplePlayerIdPrefix = "cache:accountApple:playerId:"
)

type (
	AccountAppleModel interface {
		Insert(data *AccountApple) (sql.Result, error)
		FindOne(account string) (*AccountApple, error)
		FindOneByPlayerId(playerId int64) (*AccountApple, error)
		Update(data *AccountApple) error
		Delete(account string) error
	}

	defaultAccountAppleModel struct {
		sqlc.CachedConn
		table string
	}

	AccountApple struct {
		Account       string    `db:"account"`
		PlayerId      int64     `db:"player_id"`
		Email         string    `db:"email"`
		IdentityToken string    `db:"identity_token"`
		CreateTime    time.Time `db:"create_time"`
		LastLogin     time.Time `db:"last_login"`
	}
)

func NewAccountAppleModel(conn sqlx.SqlConn, c cache.CacheConf) AccountAppleModel {
	return &defaultAccountAppleModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`account_apple`",
	}
}

func (m *defaultAccountAppleModel) Insert(data *AccountApple) (sql.Result, error) {
	accountAppleAccountKey := fmt.Sprintf("%s%v", cacheAccountAppleAccountPrefix, data.Account)
	accountApplePlayerIdKey := fmt.Sprintf("%s%v", cacheAccountApplePlayerIdPrefix, data.PlayerId)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, accountAppleRowsExpectAutoSet)
		return conn.Exec(query, data.Account, data.PlayerId, data.Email, data.IdentityToken, data.LastLogin)
	}, accountAppleAccountKey, accountApplePlayerIdKey)
	return ret, err
}

func (m *defaultAccountAppleModel) FindOne(account string) (*AccountApple, error) {
	accountAppleAccountKey := fmt.Sprintf("%s%v", cacheAccountAppleAccountPrefix, account)
	var resp AccountApple
	err := m.QueryRow(&resp, accountAppleAccountKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `account` = ? limit 1", accountAppleRows, m.table)
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

func (m *defaultAccountAppleModel) FindOneByPlayerId(playerId int64) (*AccountApple, error) {
	accountApplePlayerIdKey := fmt.Sprintf("%s%v", cacheAccountApplePlayerIdPrefix, playerId)
	var resp AccountApple
	err := m.QueryRowIndex(&resp, accountApplePlayerIdKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `player_id` = ? limit 1", accountAppleRows, m.table)
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

func (m *defaultAccountAppleModel) Update(data *AccountApple) error {
	accountAppleAccountKey := fmt.Sprintf("%s%v", cacheAccountAppleAccountPrefix, data.Account)
	accountApplePlayerIdKey := fmt.Sprintf("%s%v", cacheAccountApplePlayerIdPrefix, data.PlayerId)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `account` = ?", m.table, accountAppleRowsWithPlaceHolder)
		return conn.Exec(query, data.PlayerId, data.Email, data.IdentityToken, data.LastLogin, data.Account)
	}, accountApplePlayerIdKey, accountAppleAccountKey)
	return err
}

func (m *defaultAccountAppleModel) Delete(account string) error {
	data, err := m.FindOne(account)
	if err != nil {
		return err
	}

	accountAppleAccountKey := fmt.Sprintf("%s%v", cacheAccountAppleAccountPrefix, account)
	accountApplePlayerIdKey := fmt.Sprintf("%s%v", cacheAccountApplePlayerIdPrefix, data.PlayerId)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `account` = ?", m.table)
		return conn.Exec(query, account)
	}, accountAppleAccountKey, accountApplePlayerIdKey)
	return err
}

func (m *defaultAccountAppleModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheAccountAppleAccountPrefix, primary)
}

func (m *defaultAccountAppleModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `account` = ? limit 1", accountAppleRows, m.table)
	return conn.QueryRow(v, query, primary)
}

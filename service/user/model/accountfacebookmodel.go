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
	accountFacebookFieldNames          = builder.RawFieldNames(&AccountFacebook{})
	accountFacebookRows                = strings.Join(accountFacebookFieldNames, ",")
	accountFacebookRowsExpectAutoSet   = strings.Join(stringx.Remove(accountFacebookFieldNames, "`create_time`", "`update_time`"), ",")
	accountFacebookRowsWithPlaceHolder = strings.Join(stringx.Remove(accountFacebookFieldNames, "`account`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheAccountFacebookAccountPrefix  = "cache:accountFacebook:account:"
	cacheAccountFacebookPlayerIdPrefix = "cache:accountFacebook:playerId:"
)

type (
	AccountFacebookModel interface {
		Insert(data *AccountFacebook) (sql.Result, error)
		FindOne(account string) (*AccountFacebook, error)
		FindOneByPlayerId(playerId int64) (*AccountFacebook, error)
		Update(data *AccountFacebook) error
		Delete(account string) error
	}

	defaultAccountFacebookModel struct {
		sqlc.CachedConn
		table string
	}

	AccountFacebook struct {
		Account       string         `db:"account"`
		PlayerId      int64          `db:"player_id"`
		IdentityToken sql.NullString `db:"identity_token"`
		AccessToken   sql.NullString `db:"access_token"`
		CreateTime    time.Time      `db:"create_time"`
		LastLogin     time.Time      `db:"last_login"`
	}
)

func NewAccountFacebookModel(conn sqlx.SqlConn, c cache.CacheConf) AccountFacebookModel {
	return &defaultAccountFacebookModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`account_facebook`",
	}
}

func (m *defaultAccountFacebookModel) Insert(data *AccountFacebook) (sql.Result, error) {
	accountFacebookAccountKey := fmt.Sprintf("%s%v", cacheAccountFacebookAccountPrefix, data.Account)
	accountFacebookPlayerIdKey := fmt.Sprintf("%s%v", cacheAccountFacebookPlayerIdPrefix, data.PlayerId)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, accountFacebookRowsExpectAutoSet)
		return conn.Exec(query, data.Account, data.PlayerId, data.IdentityToken, data.AccessToken, data.LastLogin)
	}, accountFacebookAccountKey, accountFacebookPlayerIdKey)
	return ret, err
}

func (m *defaultAccountFacebookModel) FindOne(account string) (*AccountFacebook, error) {
	accountFacebookAccountKey := fmt.Sprintf("%s%v", cacheAccountFacebookAccountPrefix, account)
	var resp AccountFacebook
	err := m.QueryRow(&resp, accountFacebookAccountKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `account` = ? limit 1", accountFacebookRows, m.table)
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

func (m *defaultAccountFacebookModel) FindOneByPlayerId(playerId int64) (*AccountFacebook, error) {
	accountFacebookPlayerIdKey := fmt.Sprintf("%s%v", cacheAccountFacebookPlayerIdPrefix, playerId)
	var resp AccountFacebook
	err := m.QueryRowIndex(&resp, accountFacebookPlayerIdKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `player_id` = ? limit 1", accountFacebookRows, m.table)
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

func (m *defaultAccountFacebookModel) Update(data *AccountFacebook) error {
	accountFacebookAccountKey := fmt.Sprintf("%s%v", cacheAccountFacebookAccountPrefix, data.Account)
	accountFacebookPlayerIdKey := fmt.Sprintf("%s%v", cacheAccountFacebookPlayerIdPrefix, data.PlayerId)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `account` = ?", m.table, accountFacebookRowsWithPlaceHolder)
		return conn.Exec(query, data.PlayerId, data.IdentityToken, data.AccessToken, data.LastLogin, data.Account)
	}, accountFacebookAccountKey, accountFacebookPlayerIdKey)
	return err
}

func (m *defaultAccountFacebookModel) Delete(account string) error {
	data, err := m.FindOne(account)
	if err != nil {
		return err
	}

	accountFacebookAccountKey := fmt.Sprintf("%s%v", cacheAccountFacebookAccountPrefix, account)
	accountFacebookPlayerIdKey := fmt.Sprintf("%s%v", cacheAccountFacebookPlayerIdPrefix, data.PlayerId)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `account` = ?", m.table)
		return conn.Exec(query, account)
	}, accountFacebookAccountKey, accountFacebookPlayerIdKey)
	return err
}

func (m *defaultAccountFacebookModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheAccountFacebookAccountPrefix, primary)
}

func (m *defaultAccountFacebookModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `account` = ? limit 1", accountFacebookRows, m.table)
	return conn.QueryRow(v, query, primary)
}

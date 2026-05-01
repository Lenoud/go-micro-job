package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	OpLogTableName = "b_op_log"
	opLogFields    = "id, IFNULL(request_id,'') AS request_id, IFNULL(user_id,'') AS user_id, IFNULL(re_ip,'') AS re_ip, IFNULL(re_time,0) AS re_time, IFNULL(re_ua,'') AS re_ua, IFNULL(re_url,'') AS re_url, IFNULL(re_method,'') AS re_method, IFNULL(re_content,'') AS re_content, IFNULL(success,'1') AS success, IFNULL(biz_code,0) AS biz_code, IFNULL(biz_msg,'') AS biz_msg, IFNULL(access_time,0) AS access_time"
)

var loginLogPaths = []string{"/api/user/login", "/api/user/userLogin"}

type OpLog struct {
	Id         string `db:"id"          json:"id"`
	RequestId  string `db:"request_id"  json:"requestId"`
	UserId     string `db:"user_id"     json:"userId"`
	ReIp       string `db:"re_ip"       json:"reIp"`
	ReTime     int64  `db:"re_time"     json:"reTime"`
	ReUa       string `db:"re_ua"       json:"reUa"`
	ReUrl      string `db:"re_url"      json:"reUrl"`
	ReMethod   string `db:"re_method"   json:"reMethod"`
	ReContent  string `db:"re_content"  json:"reContent"`
	Success    string `db:"success"     json:"success"`
	BizCode    int64  `db:"biz_code"    json:"bizCode"`
	BizMsg     string `db:"biz_msg"     json:"bizMsg"`
	AccessTime int64  `db:"access_time" json:"accessTime"`
}

type OpLogModel interface {
	BatchInsert(ctx context.Context, logs []*OpLog) error
	CountOpLogList(ctx context.Context) (int64, error)
	CountLoginLogList(ctx context.Context) (int64, error)
	FindOpLogList(ctx context.Context, page, pageSize int64) ([]*OpLog, error)
	FindLoginLogList(ctx context.Context, page, pageSize int64) ([]*OpLog, error)
	DeleteBefore(ctx context.Context, beforeTime int64) error
}

type defaultOpLogModel struct {
	conn  sqlx.SqlConn
	table string
}

func NewOpLogModel(conn sqlx.SqlConn) OpLogModel {
	return &defaultOpLogModel{
		conn:  conn,
		table: OpLogTableName,
	}
}

func (m *defaultOpLogModel) BatchInsert(ctx context.Context, logs []*OpLog) error {
	if len(logs) == 0 {
		return nil
	}
	placeholders := make([]string, 0, len(logs))
	args := make([]interface{}, 0, len(logs)*12)
	for _, log := range logs {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		args = append(args,
			log.RequestId,
			log.UserId,
			log.ReIp,
			log.ReTime,
			log.ReUa,
			log.ReUrl,
			log.ReMethod,
			log.ReContent,
			log.Success,
			log.BizCode,
			log.BizMsg,
			log.AccessTime,
		)
	}
	query := fmt.Sprintf("INSERT INTO %s (request_id, user_id, re_ip, re_time, re_ua, re_url, re_method, re_content, success, biz_code, biz_msg, access_time) VALUES %s",
		m.table, strings.Join(placeholders, ", "))
	_, err := m.conn.ExecCtx(ctx, query, args...)
	return err
}

func (m *defaultOpLogModel) CountOpLogList(ctx context.Context) (int64, error) {
	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE re_url NOT IN (?, ?)", m.table)
	err := m.conn.QueryRowCtx(ctx, &count, query, loginLogPaths[0], loginLogPaths[1])
	return count, err
}

func (m *defaultOpLogModel) CountLoginLogList(ctx context.Context) (int64, error) {
	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE re_url IN (?, ?)", m.table)
	err := m.conn.QueryRowCtx(ctx, &count, query, loginLogPaths[0], loginLogPaths[1])
	return count, err
}

func (m *defaultOpLogModel) FindOpLogList(ctx context.Context, page, pageSize int64) ([]*OpLog, error) {
	offset := (page - 1) * pageSize

	var resp []*OpLog
	query := fmt.Sprintf("SELECT %s FROM %s WHERE re_url NOT IN (?, ?) ORDER BY re_time DESC LIMIT ? OFFSET ?", opLogFields, m.table)
	if err := m.conn.QueryRowsCtx(ctx, &resp, query, loginLogPaths[0], loginLogPaths[1], pageSize, offset); err != nil {
		return nil, err
	}
	if resp == nil {
		resp = []*OpLog{}
	}
	return resp, nil
}

func (m *defaultOpLogModel) FindLoginLogList(ctx context.Context, page, pageSize int64) ([]*OpLog, error) {
	offset := (page - 1) * pageSize

	var resp []*OpLog
	query := fmt.Sprintf("SELECT %s FROM %s WHERE re_url IN (?, ?) ORDER BY re_time DESC LIMIT ? OFFSET ?", opLogFields, m.table)
	if err := m.conn.QueryRowsCtx(ctx, &resp, query, loginLogPaths[0], loginLogPaths[1], pageSize, offset); err != nil {
		return nil, err
	}
	if resp == nil {
		resp = []*OpLog{}
	}
	return resp, nil
}

func (m *defaultOpLogModel) DeleteBefore(ctx context.Context, beforeTime int64) error {
	for {
		query := fmt.Sprintf("DELETE FROM %s WHERE re_time < ? LIMIT 10000", m.table)
		result, err := m.conn.ExecCtx(ctx, query, beforeTime)
		if err != nil {
			return err
		}
		affected, _ := result.RowsAffected()
		if affected < 10000 {
			return nil
		}
	}
}


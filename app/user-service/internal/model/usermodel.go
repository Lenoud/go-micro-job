package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	UserTableName = "b_user"
	userFields    = "id, username, password, IFNULL(nickname,'') AS nickname, IFNULL(mobile,'') AS mobile, IFNULL(email,'') AS email, IFNULL(role,'1') AS role, IFNULL(status,'0') AS status, create_time"
)

type User struct {
	Id         string   `db:"id"          json:"id"`
	Username   string   `db:"username"    json:"username"`
	Password   string   `db:"password"    json:"password"`
	Nickname   string   `db:"nickname"    json:"nickname"`
	Mobile     string   `db:"mobile"      json:"mobile"`
	Email      string   `db:"email"       json:"email"`
	Role       string   `db:"role"        json:"role"`
	Status     string   `db:"status"      json:"status"`
	CreateTime NullTime `db:"create_time" json:"createTime"`
}

type UserModel interface {
	Insert(ctx context.Context, data *User) (sql.Result, error)
	FindOne(ctx context.Context, id string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindAdminUser(ctx context.Context, username, password string) (*User, error)
	FindNormalUser(ctx context.Context, username, password string) (*User, error)
	FindList(ctx context.Context, keyword string, page, pageSize int64) ([]*User, int64, error)
	Delete(ctx context.Context, id string) error
	DeleteBatch(ctx context.Context, ids []string) error
	Update(ctx context.Context, data *User) error
	UpdatePassword(ctx context.Context, id, hashedPassword string) error
}

type defaultUserModel struct {
	conn  sqlx.SqlConn
	table string
}

func NewUserModel(conn sqlx.SqlConn) UserModel {
	return &defaultUserModel{
		conn:  conn,
		table: UserTableName,
	}
}

func (m *defaultUserModel) Insert(ctx context.Context, data *User) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (username, password, nickname, mobile, email, role, status, create_time) VALUES (?, ?, ?, ?, ?, ?, ?, NOW())", m.table)
	return m.conn.ExecCtx(ctx, query,
		data.Username, data.Password, data.Nickname, data.Mobile, data.Email,
		data.Role, data.Status,
	)
}

func (m *defaultUserModel) FindOne(ctx context.Context, id string) (*User, error) {
	var resp User
	query := fmt.Sprintf("SELECT %s FROM %s WHERE id = ? LIMIT 1", userFields, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch {
	case err == nil:
		return &resp, nil
	case err == sqlx.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *defaultUserModel) FindByUsername(ctx context.Context, username string) (*User, error) {
	var resp User
	query := fmt.Sprintf("SELECT %s FROM %s WHERE username = ? LIMIT 1", userFields, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, username)
	switch {
	case err == nil:
		return &resp, nil
	case err == sqlx.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *defaultUserModel) FindAdminUser(ctx context.Context, username, password string) (*User, error) {
	var resp User
	query := fmt.Sprintf("SELECT %s FROM %s WHERE username = ? AND password = ? AND role = '3' LIMIT 1", userFields, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, username, password)
	switch {
	case err == nil:
		return &resp, nil
	case err == sqlx.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *defaultUserModel) FindNormalUser(ctx context.Context, username, password string) (*User, error) {
	var resp User
	query := fmt.Sprintf("SELECT %s FROM %s WHERE username = ? AND password = ? AND role IN ('1','2') LIMIT 1", userFields, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, username, password)
	switch {
	case err == nil:
		return &resp, nil
	case err == sqlx.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *defaultUserModel) FindList(ctx context.Context, keyword string, page, pageSize int64) ([]*User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var total int64
	if strings.TrimSpace(keyword) != "" {
		countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE username LIKE ?", m.table)
		_ = m.conn.QueryRowCtx(ctx, &total, countQuery, "%"+keyword+"%")
	} else {
		countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", m.table)
		_ = m.conn.QueryRowCtx(ctx, &total, countQuery)
	}

	var resp []*User
	if strings.TrimSpace(keyword) != "" {
		query := fmt.Sprintf("SELECT %s FROM %s WHERE username LIKE ? ORDER BY create_time DESC LIMIT ? OFFSET ?", userFields, m.table)
		err := m.conn.QueryRowsCtx(ctx, &resp, query, "%"+keyword+"%", pageSize, offset)
		if err != nil {
			return nil, 0, err
		}
	} else {
		query := fmt.Sprintf("SELECT %s FROM %s ORDER BY create_time DESC LIMIT ? OFFSET ?", userFields, m.table)
		err := m.conn.QueryRowsCtx(ctx, &resp, query, pageSize, offset)
		if err != nil {
			return nil, 0, err
		}
	}
	return resp, total, nil
}

func (m *defaultUserModel) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultUserModel) DeleteBatch(ctx context.Context, ids []string) error {
	inClause, args := buildInClauseArgs(ids)
	if inClause == "" {
		return nil
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE id IN (%s)", m.table, inClause)
	_, err := m.conn.ExecCtx(ctx, query, args...)
	return err
}

func (m *defaultUserModel) Update(ctx context.Context, data *User) error {
	query := fmt.Sprintf("UPDATE %s SET username=?, nickname=?, mobile=?, email=?, role=?, status=? WHERE id=?", m.table)
	_, err := m.conn.ExecCtx(ctx, query,
		data.Username, data.Nickname, data.Mobile, data.Email, data.Role, data.Status, data.Id,
	)
	return err
}

func (m *defaultUserModel) UpdatePassword(ctx context.Context, id, hashedPassword string) error {
	query := fmt.Sprintf("UPDATE %s SET password=? WHERE id=?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, hashedPassword, id)
	return err
}

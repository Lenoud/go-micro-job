package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	DepartmentTableName = "b_department"
	departmentFields    = "id, IFNULL(title,'') AS title, IFNULL(description,'') AS description, IFNULL(parent_id,'0') AS parent_id, create_time"
)

type Department struct {
	Id          string   `db:"id"          json:"id"`
	Title       string   `db:"title"       json:"title"`
	Description string   `db:"description" json:"description"`
	ParentId    string   `db:"parent_id"   json:"parentId"`
	CreateTime  NullTime `db:"create_time" json:"createTime"`
}

type DepartmentModel interface {
	Insert(ctx context.Context, data *Department) (sql.Result, error)
	FindOne(ctx context.Context, id string) (*Department, error)
	FindList(ctx context.Context, keyword string, page, pageSize int64) ([]*Department, int64, error)
	Update(ctx context.Context, data *Department) error
	Delete(ctx context.Context, ids string) error
}

type defaultDepartmentModel struct {
	conn  sqlx.SqlConn
	table string
}

func NewDepartmentModel(conn sqlx.SqlConn) DepartmentModel {
	return &defaultDepartmentModel{
		conn:  conn,
		table: DepartmentTableName,
	}
}

func (m *defaultDepartmentModel) Insert(ctx context.Context, data *Department) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (title, description, parent_id, create_time) VALUES (?, ?, ?, NOW())", m.table)
	return m.conn.ExecCtx(ctx, query, data.Title, data.Description, data.ParentId)
}

func (m *defaultDepartmentModel) FindOne(ctx context.Context, id string) (*Department, error) {
	var resp Department
	query := fmt.Sprintf("SELECT %s FROM %s WHERE id = ? LIMIT 1", departmentFields, m.table)
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

func (m *defaultDepartmentModel) FindList(ctx context.Context, keyword string, page, pageSize int64) ([]*Department, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	keyword = strings.TrimSpace(keyword)

	var total int64
	if keyword != "" {
		countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE title LIKE ?", m.table)
		_ = m.conn.QueryRowCtx(ctx, &total, countQuery, "%"+keyword+"%")
	} else {
		countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", m.table)
		_ = m.conn.QueryRowCtx(ctx, &total, countQuery)
	}

	var resp []*Department
	if keyword != "" {
		query := fmt.Sprintf("SELECT %s FROM %s WHERE title LIKE ? ORDER BY create_time DESC LIMIT ? OFFSET ?", departmentFields, m.table)
		if err := m.conn.QueryRowsCtx(ctx, &resp, query, "%"+keyword+"%", pageSize, offset); err != nil {
			return nil, 0, err
		}
	} else {
		query := fmt.Sprintf("SELECT %s FROM %s ORDER BY create_time DESC LIMIT ? OFFSET ?", departmentFields, m.table)
		if err := m.conn.QueryRowsCtx(ctx, &resp, query, pageSize, offset); err != nil {
			return nil, 0, err
		}
	}
	if resp == nil {
		resp = []*Department{}
	}
	return resp, total, nil
}

func (m *defaultDepartmentModel) Update(ctx context.Context, data *Department) error {
	query := fmt.Sprintf("UPDATE %s SET title=?, description=?, parent_id=? WHERE id=?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, data.Title, data.Description, data.ParentId, data.Id)
	return err
}

func (m *defaultDepartmentModel) Delete(ctx context.Context, ids string) error {
	idList := splitIDs(ids)
	if len(idList) == 0 {
		return nil
	}
	placeholders := make([]string, len(idList))
	args := make([]interface{}, len(idList))
	for i, id := range idList {
		placeholders[i] = "?"
		args[i] = id
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE id IN (%s)", m.table, strings.Join(placeholders, ","))
	_, err := m.conn.ExecCtx(ctx, query, args...)
	return err
}

func splitIDs(raw string) []string {
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

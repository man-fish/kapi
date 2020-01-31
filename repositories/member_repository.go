package repositories

import (
	"Kapi/config"
	"Kapi/models"
	"Kapi/utils"
	"database/sql"
	"fmt"
	"errors"
)

type IMemberManager interface {
	Conn() error
	InsertOne(*models.Member) (int64, error)
	DeleteOne(int64, int64) (bool, error)
	UpdateOne(*models.Member) (int64, error)
	SelectAllByGid(int64) ([]*models.UserMember, error)
	SelectOne(int64, int64) (*models.Member, error)
	IsAllowed(int64, int64) (bool, error)
}

type MemberManager struct {
	mysqlCon	*sql.DB
	tableName	string
}

func NewMemberManager(db *sql.DB, tableName string) IMemberManager {
	return &MemberManager{
		mysqlCon:  db,
		tableName: tableName,
	}
}

func(m *MemberManager) Conn() (err error) {
	if m.mysqlCon == nil {
		AppConfig, err := config.GetConfig(utils.RootPath()+"/config/config.json")
		if err != nil {
			return err
		}
		mysql, err := utils.NewMysqlConn(AppConfig.MysqlDsn)
		if err != nil {
			return err
		}
		m.mysqlCon = mysql
	}
	if m.tableName == "" {
		m.tableName = "Kapi_member"
	}
	return
}

func (m *MemberManager) InsertOne(member *models.Member) (mid int64, err error) {
	if err = m.Conn(); err != nil {
		return
	}
	sql := fmt.Sprintf("INSERT %v SET uid = ?, gid = ?, role = ?",m.tableName)
	stmt, err := m.mysqlCon.Prepare(sql)
	if err != nil {
		return
	}
	result, err := stmt.Exec(member.Uid, member.Gid, member.Role)
	if err != nil {
		return
	}
	return result.LastInsertId()
}

func (m *MemberManager) DeleteOne(gid, uid int64) (ok bool, err error) {
	if err = m.Conn(); err != nil {
		return
	}
	sql := "delete from Kapi_member where gid = ? and uid = ?"
	stmt, err := m.mysqlCon.Prepare(sql)
	if err != nil {
		return
	}
	result, err := stmt.Exec(gid,uid)
	count, err := result.RowsAffected()
	if count == 0 {
		err = errors.New("不存在的成员")
		return
	}
	return true, err
}

func (m *MemberManager) UpdateOne(member *models.Member) (mid int64, err error) {
	if err = m.Conn(); err != nil {
		return
	}
	sql := fmt.Sprintf("UPDATE %v SET uid = ?, gid = ?, role = ? WHERE id = ?", m.tableName)
	stmt, err := m.mysqlCon.Prepare(sql)
	if err != nil {
		return
	}
	result, err := stmt.Exec(member.Uid,member.Gid,member.Role,member.ID)
	if err != nil {
		return
	}
	return result.LastInsertId()
}

func (m *MemberManager) SelectAllByGid(gid int64) (userMembers []*models.UserMember, err error) {
	if err = m.Conn(); err != nil {
		return
	}
	sql := "select * from Kapi_member LEFT OUTER JOIN Kapi_user ON Kapi_member.gid = ?"
	stmt, err := m.mysqlCon.Prepare(sql)
	if err != nil {
		return
	}
	rows, err := stmt.Query(gid)
	if err != nil {
		return
	}
	results := utils.RowsToMap(rows)
	if len(results) == 0 {
		err = errors.New("没有查询到有效信息")
		return
	}
	for _, result := range results {
		userMemberResult := new(models.UserMember)
		err := utils.MapToStructByTagSql(result,userMemberResult)
		if err != nil {
			continue
		}
		userMembers = append(userMembers,userMemberResult)
	}
	return
}

func (m *MemberManager) SelectOne(uid,gid int64) (member *models.Member, err error) {
	if err = m.Conn(); err != nil {
		return
	}
	sql := "select * from Kapi_member where uid = ? and gid = ?"
	stmt, err := m.mysqlCon.Prepare(sql)
	if err != nil {
		return
	}
	rows, err := stmt.Query(uid, gid)
	if err != nil {
		return
	}
	result := utils.RowToMap(rows)
	if len(result) == 0 {
		err = errors.New("没有查询到有效信息")
		return
	}
	member = new(models.Member)
	err = utils.MapToStructByTagSql(result,member)
	if err != nil {
		return
	}
	return
}

func (m *MemberManager) IsAllowed(uid, gid int64) (allowed bool, err error) {
	if err = m.Conn(); err != nil {
		return
	}
	member, err := m.SelectOne(uid, gid)
	if err != nil {
		return
	}
	return member.Role == "leader", err
}



package repositories

import (
	"Kapi/config"
	"Kapi/models"
	"Kapi/utils"
	validators "Kapi/validator"
	"database/sql"
	"errors"
	"fmt"
)

type IGroupManager interface {
	Conn() error
	InsertOne(*validators.CreateGroupValidator, *models.Member) (int64, error)	/* 新增组的时候用户自动成为leader，事务操作 */
	UpdateOne(*validators.UpdateGroupValidator) (int64, error)
	SelectOne(int64) (*models.Group, error)
	SelectAllByUid(int64) ([]*models.Group, error)
}

type GroupManager struct {
	mysqlCon	*sql.DB
	tableName	string
}

func NewGroupManager(db *sql.DB, tableName string) IGroupManager {
	return &GroupManager{
		mysqlCon:  db,
		tableName: tableName,
	}
}

func(u *GroupManager) Conn() (err error) {
	if u.mysqlCon == nil {
		AppConfig, err := config.GetConfig(utils.RootPath()+"/config/config.json")
		if err != nil {
			return err
		}
		mysql, err := utils.NewMysqlConn(AppConfig.MysqlDsn)
		if err != nil {
			return err
		}
		u.mysqlCon = mysql
	}
	if u.tableName == "" {
		u.tableName = "Kapi_group"
	}
	return
}

func(u *GroupManager) InsertOne (group *validators.CreateGroupValidator, member *models.Member) (gid int64, err error) {
	if err = u.Conn(); err != nil {
		return
	}
	tx, err := u.mysqlCon.Begin()
	if err != nil {
		if tx != nil {
			_ = tx.Rollback()
		}
		return
	}
	sql := fmt.Sprintf("INSERT %v SET group_name = ?, group_desc = ?, type = ?",u.tableName)
	stxt, err := tx.Prepare(sql)
	if err != nil {
		_ = tx.Rollback()
		return
	}
	result, err := stxt.Exec(group.GroupName, group.GroupDesc, group.Type)
	if err != nil {
		_ = tx.Rollback()
		return
	}
	gid, err = result.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return
	}

	sql = fmt.Sprintf("INSERT %v SET uid = ?, gid = ?, role = ?", "Kapi_member")
	stxt, err = tx.Prepare(sql)
	if err != nil {
		_ = tx.Rollback()
		return
	}
	result, err = stxt.Exec(member.Uid, gid, member.Role)
	if err != nil {
		_ = tx.Rollback()
		return
	}
	_, err = result.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
	}
	return
}

func(u *GroupManager) UpdateOne(group *validators.UpdateGroupValidator) (gid int64,err error) {
	if err = u.Conn(); err != nil {
		return
	}
	sql := "UPDATE Kapi_group SET group_name = ?, group_desc = ?, custom_field = ?, custom_field_enable = ? WHERE id = ?"
	stmt, err := u.mysqlCon.Prepare(sql)
	if err != nil {
		return
	}
	result, err := stmt.Exec(group.GroupName, group.GroupDesc, group.CustomField, group.CustomFieldEnabled, group.ID)
	if err != nil {
		return
	}
	return result.LastInsertId()
}

func(u *GroupManager) SelectAllByUid(uid int64) (groups []*models.Group, err error) {
	if err = u.Conn(); err != nil {
		return
	}
	sql := "SELECT * FROM Kapi_group WHERE id IN (SELECT gid FROM Kapi_member WHERE uid = ?)"
	stmt, err := u.mysqlCon.Prepare(sql)
	if err != nil {
		return
	}
	rows, err := stmt.Query(uid)
	if err != nil {
		return
	}
	results := utils.RowsToMap(rows)
	if len(results) == 0 {
		return
	}
	for _, result := range results {
		groupResult := new(models.Group)
		err := utils.MapToStructByTagSql(result,groupResult)
		if err != nil {
			continue
		}
		groups = append(groups,groupResult)
	}
	return
}

func(u *GroupManager) SelectOne(gid int64) (group *models.Group, err error) {
	if err = u.Conn(); err != nil {
		return
	}
	sql := "SELECT * FROM Kapi_group WHERE id = ?"
	stmt, err := u.mysqlCon.Prepare(sql)
	if err != nil {
		return
	}
	row, err := stmt.Query(gid)
	if err != nil {
		return
	}
	result := utils.RowToMap(row)
	if len(result) == 0 {
		err = errors.New("没有查询到有效信息")
		return
	}
	group = new(models.Group)
	err = utils.MapToStructByTagSql(result, group)
	return
}
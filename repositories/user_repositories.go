package repositories

import (
	"Kapi/config"
	"Kapi/models"
	"Kapi/utils"
	"database/sql"
	"errors"
	"fmt"
)

type IUserManager interface {
	Conn() error
	SelectOne(string) (*models.User, error)
	InsertOne (*models.User) (int64, error)
}

type UserManager struct {
	mysqlCon *sql.DB
	tableName string
}

func NewUserManager(db *sql.DB, tableName string) IUserManager {
	return &UserManager{db,tableName}
}

func(u *UserManager) Conn() (err error) {
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
		u.tableName = "Kapi_user"
	}
	return
}

func (u *UserManager) SelectOne(email string) (userResult *models.User, err error) {
	if err = u.Conn(); err != nil {
		return
	}
	sql := "SELECT * FROM " + u.tableName + " WHERE email = ?"
	stmt, err := u.mysqlCon.Prepare(sql)
	if err != nil {
		return
	}
	rows, err := stmt.Query(email)
	if err != nil {
		return
	}
	results := utils.RowToMap(rows)
	if len(results) == 0 {
		err = errors.New("没有查询到有效信息")
		return
	}
	userResult = new(models.User)
	err = utils.MapToStructByTagSql(results, userResult)
	if err != nil {
		return
	}
	return
}

func (u *UserManager) InsertOne(user *models.User) (uid int64, err error) {
	if err = u.Conn(); err != nil {
		return
	}
	tx, err := u.mysqlCon.Begin()
	if err != nil {
		if tx != nil {
			_ = tx.Rollback()
			return
		}
	}
	sql := fmt.Sprintf("INSERT %v SET username = ?, email = ?, password = ?, pass_salt = ?, ip = ?",u.tableName)
	stxt, err := tx.Prepare(sql)
	if err != nil {
		_ = tx.Rollback()
		return
	}
	result, err := stxt.Exec(user.Username,user.Email,user.Password,user.PassSalt,user.IP)
	if err != nil {
		_ = tx.Rollback()
		return
	}

	uid, err =  result.LastInsertId()

	sql = fmt.Sprintf("INSERT %v SET group_name = ?, group_desc = ?, type = ?","Kapi_group")
	stxt, err = tx.Prepare(sql)
	if err != nil {
		_ = tx.Rollback()
		return
	}
	result, err = stxt.Exec("个人空间", "私有个人项目", "private")
	if err != nil {
		_ = tx.Rollback()
		return
	}
	gid, err := result.LastInsertId()
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
	result, err = stxt.Exec(uid, gid, "leader")
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
package repositories

import (
	"Kapi/config"
	"Kapi/models"
	"Kapi/utils"
	"database/sql"
	"errors"
)

type IUserManager interface {
	Conn() error
	SelectForLogin(string) (*models.User, error)
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
		u.tableName = "user"
	}
	return
}

func (u *UserManager) SelectForLogin(email string) (userResult *models.User, err error) {
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
	err = utils.MapToStructByTagSql(results, userResult)
	if err != nil {
		return
	}
	return
}
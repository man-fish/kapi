package main

import (
	"Kapi/config"
	"Kapi/models"
	"Kapi/utils"
	"fmt"
)

func main() {
	//conn, err := utils.NewMysqlConn()
	//if err != nil {
	//	panic(err)
	//}

	//username := "admin"
	//pass_salt := utils.MD5("")
	//password := utils.MD5(pass_salt)
	//email := "admin@admin.com"
	//
	//sql := "INSERT INTO user(username,password,email,pass_salt,add_time,up_time) values (?,?,?,?,?,?)"
	//stmt, err := conn.Prepare(sql)
	//if err != nil {
	//	panic(err)
	//}
	//result, err := stmt.Exec(username,password,pass_salt,email,time.Now(),nil)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(result.LastInsertId())
	//sql := "SELECT * FROM user"
	//
	//stmt, err := conn.Prepare(sql)
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//rows, err := stmt.Query()
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//defer rows.Close()
	//
	//fmt.Println(utils.SQLToMap(rows))
	//fmt.Println("love")
	//scan([]int{1,23,3,4}...)
	//fmt.Println(utils.RootPath())
	//try()

	// Token from another example.  This token is expired
	tokenstring, err := utils.DefaultToken(1)
	if err != nil {

	}
	fmt.Println(tokenstring)
	c, err := config.GetConfig("./../config/config.json")
	utils.VerifyToken(tokenstring,c.SecurityKey)

}

func try() (user *models.User) {
	fmt.Println(user)
	return user
}

func scan(nums ...int) {
	fmt.Println(nums)
}
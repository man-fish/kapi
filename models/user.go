package models

type User struct{
	ID			int64	`json:"id" sql:"id" Kapi:"id"`
	Username 	string	`json:"username" sql:"username" Kapi:"username"`
	Password 	string	`json:"password" sql:"password" Kapi:"password"`
	Email	 	string	`json:"email" sql:"email" Kapi:"email"`
	PassSalt 	string 	`json:"pass_salt" sql:"pass_salt"`
	AddTime		string	`json:"add_time" sql:"add_time" Kapi:"add_time"`
	UpTime		string	`json:"up_time" sql:"up_time" Kapi:"up_time"`
}
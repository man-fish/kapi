package models

type Member struct {
	ID					int64	`json:"id" sql:"id" `
	Gid					int64	`json:"gid" sql:"gid"`
	Uid					int64	`json:"uid" sql:"uid"`
	Role				string	`json:"role" sql:"role"`
	AddTime				string	`json:"add_time" sql:"add_time"`
	UpTime				string	`json:"up_time" sql:"up_time"`
}


type UserMember struct {
	Member
	User
}
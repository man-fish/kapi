package services

import (
	"Kapi/models"
	"Kapi/repositories"
	"Kapi/utils"
	"fmt"
)

type IUserService interface {
	LoginByEmail(email string, password string) (token string, err error)
	RegisterByEmail(username string, email string, password string,ip string) (uid int64, err error)
}

type UserService struct {
	userRepository repositories.IUserManager
}

func NewUserService(userRepository repositories.IUserManager) IUserService {
	us := &UserService{userRepository}
	return us
}

func (us *UserService) LoginByEmail(email string, password string) (token string, err error) {
	user, err := us.userRepository.SelectOne(email)
	if err != nil {
		fmt.Println(err)
		err = utils.NewError(400, "邮箱不存在")
		return
	}
	result := verifyPassword(password, user.PassSalt, user.Password)
	if !result {
		fmt.Println(err)
		err = utils.NewError(400, "密码不正确")
		return
	}
	token, err = utils.DefaultToken(user.ID)
	if err != nil {
		fmt.Println(err)
		err = utils.NewError(500, "token制作失败")
		return
	}
	return
}

func (uc *UserService) RegisterByEmail(username string, email string, password string,ip string) (uid int64, err error) {
	_, err = uc.userRepository.SelectOne(email)
	if err == nil {
		err = utils.NewError(400, "邮箱已存在")
		return
	}
	passSalt := utils.MD5("")
	passwordCrypto := utils.MD5(passSalt+password)
	user := &models.User{
		Username: username,
		Password: passwordCrypto,
		Email:    email,
		PassSalt: passSalt,
		IP:		  ip,
	}
	uid, err = uc.userRepository.InsertOne(user)
	if err != nil {
		fmt.Println(err)
		err = utils.NewError(500, "插入数据失败")
	}
	return
}

func verifyPassword(inputPass string, passSalt string,excPass string) bool {
	MD5pass := utils.MD5(passSalt+inputPass)
	if MD5pass == excPass {
		return true
	}else{
		return false
	}
}


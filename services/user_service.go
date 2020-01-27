package services

import (
	"Kapi/repositories"
	"Kapi/utils"
	"errors"
)

type IUserService interface {
	LoginByEmail(string, string) (string, error)
}

type UserService struct {
	userRepository repositories.IUserManager
}

func NewUserService(userRepository repositories.IUserManager) IUserService {
	us := &UserService{userRepository}
	return us
}

func (us *UserService) LoginByEmail(email string, password string) (token string, err error) {
	user, err := us.userRepository.SelectForLogin(email)
	if err != nil {
		return
	}
	result := verifyPassword(password, user.PassSalt, user.Password)
	if !result {
		err = errors.New("密码错误")
		return
	}
	token, err = utils.DefaultToken(user.ID)
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



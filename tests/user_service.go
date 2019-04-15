package tests

import "github.com/eddieowens/axon"

type UserService interface {
	GetUser() string
}

type UserServiceImpl struct {
}

func (UserServiceImpl) GetUser() string {
	return "im a user"
}

func UserServiceFactory() axon.Instance {
	return axon.StructPtr(new(UserServiceImpl))
}

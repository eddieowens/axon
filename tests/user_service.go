package tests

import "axon"

type UserService interface {
    axon.Instance
    GetUser() string
}

type UserServiceImpl struct {

}

func (UserServiceImpl) GetInstanceName() string {
    return "userService"
}

func (UserServiceImpl) GetUser() string {
    return "im a user"
}

func UserServiceFactory() axon.Instance {
    return new(UserServiceImpl)
}

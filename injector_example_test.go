package axon_test

import (
	"fmt"
	"github.com/eddieowens/axon"
	"os"
)

// A full-scale example of injecting values into a struct
func Example() {
	axon.Add(axon.NewKey("secret"), os.Getenv("LOGIN_SECRET"))
	axon.Add(axon.NewTypeKey[LoginServiceClient](new(loginServiceClient)))
	axon.Add(axon.NewTypeKey[UserService](new(userService)))
	axon.Add(axon.NewTypeKeyFactory[DBClient](axon.NewFactory[DBClient](func(_ axon.Injector) (DBClient, error) {
		// inject username, DB name, password, etc.
		return &dbClient{}, nil
	})))

	api := new(ApiGateway)
	_ = axon.Inject(api)
	api.UserService.DeleteUser("user")

	// Output:
	// Deleting user from DB!
	// Logout for user
	// Successfully deleted user
}

type DBClient interface {
	DeleteUser(username string)
}

type dbClient struct {
}

func (d *dbClient) DeleteUser(username string) {
	fmt.Println("Deleting", username, "from DB!")
}

type LoginServiceClient interface {
	Logout(username string)
}

type loginServiceClient struct {
	// Injects the secret via a key called "secret"
	LoginSecret string `inject:"secret"`
}

func (l *loginServiceClient) Logout(username string) {
	fmt.Println("Logout for", username)
}

type UserService interface {
	DeleteUser(username string)
}

type userService struct {
	// Injects the default LoginServiceClient implementation.
	LoginServiceClient LoginServiceClient `inject:",type"`
	DBClient           DBClient           `inject:",type"`
}

func (u *userService) DeleteUser(username string) {
	u.DBClient.DeleteUser(username)
	u.LoginServiceClient.Logout(username)
	fmt.Println("Successfully deleted user")
}

type ApiGateway struct {
	UserService UserService `inject:",type"`
}

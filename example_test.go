package axon_test

import (
	"fmt"
	"github.com/eddieowens/axon"
)

// Whenever Injector.Add is called, the value is updated within the Injector only. If you've already called Get from the
// Injector before an Add is called, the value you hold is not updated, for example
func ExampleNewProvider() {
	type Server struct {
		ApiKey *axon.Provider[string] `inject:"api_key"`
	}

	axon.Add("api_key", axon.NewProvider("123"))

	server := new(Server)
	_ = axon.Inject(server)

	fmt.Println(server.ApiKey.Get())

	axon.Add("api_key", axon.NewProvider("456"))

	fmt.Println(server.ApiKey.Get())

	// Output:
	// 123
	// 456
}

func ExampleNewFactory() {
	type Service struct {
		DBClient DBClient `inject:",type"`
	}

	axon.Add(axon.NewTypeKeyFactory[DBClient](axon.NewFactory[DBClient](func(_ axon.Injector) (DBClient, error) {
		// construct the DB client.
		return &dbClient{}, nil
	})))

	s := new(Service)
	_ = axon.Inject(s)
	s.DBClient.DeleteUser("user")

	// Output:
	// Deleting user from DB!
}

// Simple example of adding values to the injector.
func ExampleAdd() {
	// Use a string key
	axon.Add("key", "val")
	val := axon.MustGet[string](axon.WithKey("key"))
	fmt.Println(val)

	// Use a type as a key
	axon.Add(axon.NewTypeKey(1))
	i := axon.MustGet[int]()
	fmt.Println(i)

	// Output:
	// val
	// 1
}

// Simple example of injecting values into a struct
func ExampleInject() {
	type ExampleStruct struct {
		Val string `inject:"val"`
		Int int    `inject:",type"`
	}

	axon.Add("val", "val")
	axon.Add(axon.NewTypeKey[int](1))

	// Only struct pointers allowed.
	out := new(ExampleStruct)
	_ = axon.Inject(out)

	fmt.Println(out.Val)
	fmt.Println(out.Int)

	// Output:
	// val
	// 1
}

// Simple example of getting values from the injector.
func ExampleGet() {
	axon.Add("key", "val")
	val := axon.MustGet[string](axon.WithKey("key"))
	fmt.Println(val)

	// Output:
	// val
}

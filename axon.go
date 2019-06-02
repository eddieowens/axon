// A simple, lightweight, and lazy-loaded DI (really just a singleton management) library.
package axon

// Creates and instantiates an Injector via a Binder provided by the user.
//   binder := axon.NewBinder(axon.NewPackage(
//     axon.Bind("Car").To().Factory(CarFactory).WithoutArgs()),
//     axon.Bind("Engine").To().Instance(axon.StructPtr(new(EngineImpl)),
//     axon.Bind("LockCode").To().String(os.Getenv("LOCK_CODE")),
//   ))
//
//   injector := axon.NewInjector(binder)
//
// The user can now retrieve their Instances via a call to the Injector's Get*(string) methods.
//   fmt.Println(injector.GetString("LockCode")) // Prints the value of the env var "LOCK_CODE"
//
// Or the Binding can be replaced before running a test via the injector.Add(...) method.
//   injector.Add("Engine", axon.StructPtr(new(EngineMock)))
func NewInjector(binder Binder) Injector {
	return newInjector(binder)
}

// A simple, lightweight, lazy-loaded, reflectionless, and concurrent DI (really just a singleton management) library.
package axon

// Creates and instantiates an Injector via a Binder provided by the user e.g.
//   binder := []axon.BinderEntry{
//     {Instance: testServiceImpl{}},
//   }
//
//   injector := axon.NewInjector(binder)
func NewInjector(binder Binder) Injector {
	return newInjector(binder)
}

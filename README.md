# axon

[![Go Report Card](https://goreportcard.com/badge/github.com/eddieowens/axon)](https://goreportcard.com/report/github.com/eddieowens/axon)
[![License](https://img.shields.io/badge/License-Apache%202.0-yellowgreen.svg)](https://github.com/eddieowens/axon/blob/master/LICENSE)
[![godoc](https://img.shields.io/badge/godoc-reference-blue)](https://pkg.go.dev/github.com/eddieowens/axon?tab=doc)

A simple, lightweight, and lazy-loaded DI (really just a singleton management) library that supports generics.
Influenced by multiple DI
libraries but more specifically Google's [Guice](https://github.com/google/guice).

## Install

```bash
go get github.com/eddieowens/axon
```

## Usage

### Basic

Simple add and get with a string key

```go
package main

import (
  "fmt"
  "github.com/eddieowens/axon"
)

func main() {
  axon.Add("AnswerToTheUltimateQuestion", 42)
  answer := axon.MustGet(axon.WithKey("AnswerToTheUltimateQuestion"))
  fmt.Println(answer) // prints 42
}

```

You can also use a type as a key

```go
package main

import (
  "fmt"
  "github.com/eddieowens/axon"
)

func main() {
  axon.Add(axon.NewTypeKey[int](42))
  answer := axon.MustGet[int]()
  fmt.Println(answer) // prints 42
}
```

### Injecting dependencies

To inject dependencies to a struct, you can use the `Inject` func.

```go
package main

import (
  "fmt"
  "github.com/eddieowens/axon"
)

type Struct struct {
  Answer int `inject:"AnswerToTheUltimateQuestion"`
}

func main() {
  axon.Add("AnswerToTheUltimateQuestion", 42)

  s := new(Struct)
  _ = axon.Inject(s)
  fmt.Println(s.Answer) // prints 42
}
```

A more full fledged example

```go
package main

import (
  "fmt"
  "github.com/eddieowens/axon"
  "os"
)

type DatabaseClient interface {
  DeleteUser(user string) error
}

type databaseClient struct {
}

func (d *databaseClient) DeleteUser(_ string) error {
  fmt.Println("Deleting user!")
  return nil
}

type ServerConfig struct {
  Port int `inject:"port"`
}

type Server struct {
  // inject whatever is the default for the DBClient type
  DB           DatabaseClient `inject:",type"`
  ServerConfig ServerConfig   `inject:"config"`
}

func main() {
  axon.Add("port", os.Getenv("PORT"))

  // default implementation for DatabaseClient
  axon.Add(axon.NewTypeKey[DatabaseClient](new(databaseClient)))

  // construct the Config whenever it's needed (only ever called once)
  axon.Add("config", axon.NewFactory[ServerConfig](func(_ axon.Injector) (ServerConfig, error) {
    return ServerConfig{}, nil
  }))

  s := new(Server)
  _ = axon.Inject(s)
  fmt.Println(s.ServerConfig.Port)     // prints value of env var PORT
  fmt.Println(s.DB.DeleteUser("user")) // prints Deleting user!
}
```

For more examples and info, check out the [GoDoc](https://pkg.go.dev/github.com/eddieowens/axon?tab=doc)

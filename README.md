# go-env
[![LICENSE](https://img.shields.io/github/license/mashape/apistatus.svg?style=flat-square&label=License)](https://github.com/czasg/go-env/blob/master/LICENSE)
[![codecov](https://codecov.io/gh/czasg/go-env/branch/main/graph/badge.svg?token=OkiSH6DMqf)](https://codecov.io/gh/czasg/go-env)
[![GitHub Stars](https://img.shields.io/github/stars/czasg/go-env.svg?style=flat-square&label=Stars&logo=github)](https://github.com/czasg/go-env/stargazers)

## Env Parse
```go
package main

import (
	"github.com/czasg/go-env"
	"os"
)

type Config struct {
	Env     string
	Postgres
}

type Postgres struct {
	Addr     string
	User     string
	Password string
	Database string
}

func main() {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}
	if cfg.Env != "test" {
		panic("fail")
	}
	if cfg.Postgres.User != "postgres" {
		panic("fail")
	}
}

func init() {
	_ = os.Setenv("ENV", "test")
	_ = os.Setenv("POSTGRES_ADDR", "localhost:5432")
	_ = os.Setenv("POSTGRES_USER", "postgres")
	_ = os.Setenv("POSTGRES_PASSWORD", "postgres")
	_ = os.Setenv("POSTGRES_DATABASE", "postgres")
}
```

## Env Tag
|tag|comment|
|---|---|
|env:"fieldName"|default is struct field name, you can also point a new fieldName.|
|env:",default=value"|set default env value.|
|env:",require"|it return an err when env is not found.|
|env:",empty"|set current fieldName to an empty string like "".|
|env:",sep=_"|when struct into struct, sep is the connector, default is "_".|

```go
package main

import (
	"github.com/czasg/go-env"
	"os"
)

type Config struct {
	RPC      `env:"GRPC"`
	Redis    `env:"RDS"`
	Postgres `env:"PG"`
}

type RPC struct {
	gRPC `env:",empty"`
}

type gRPC struct {
	Addr string
}

type Redis struct {
	Addr     string `env:",sep=!"`
	Password string `env:",sep=@"`
	DB       int    `env:",sep=#"`
}

type Postgres struct {
	Addr     string `env:",default=localhost:5432"`
	User     string `env:",default=postgres"`
	Password string `env:",default=postgres"`
	Database string `env:",default=postgres"`
}

func main() {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}
	if cfg.RPC.Addr != "localhost:9000" {
		panic(err)
	}
	if cfg.Redis.Addr != "localhost:6379" {
		panic(err)
	}
}

func init() {
	_ = os.Setenv("GRPC_ADDR", "localhost:9000")
	_ = os.Setenv("RDS!ADDR", "localhost:6379")
	_ = os.Setenv("RDS@PASSWORD", "123456")
	_ = os.Setenv("RDS#DB", "10")
}
```

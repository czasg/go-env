package env

import (
	"fmt"
	"os"
	"testing"
	"time"
)

type TestStructParseEnv struct {
	TestStructParseEnv1 TestStructParseEnv1 `env:"TestStructParseEnv1"`
	TestStructParseEnv2 TestStructParseEnv2 `env:"TestStructParseEnv2"`
}

type TestStructParseEnv1 struct {
	A                   string
	B                   int
	TestStructParseEnv2 TestStructParseEnv2 `env:"xxx,sep=--"`
}

type TestStructParseEnv2 struct {
	A string
	B int
}

func TestStructParse(t *testing.T) {
	assert := assertWrap(t)
	{
		nowUnix := int(time.Now().Unix())
		randomA := fmt.Sprintf("%d", nowUnix)
		_ = os.Setenv("TestStructParseEnv1_A", randomA)
		_ = os.Setenv("TestStructParseEnv1_B", randomA)
		_ = os.Setenv("TestStructParseEnv1--xxx_A", randomA)
		_ = os.Setenv("TestStructParseEnv1--xxx_B", randomA)
		_ = os.Setenv("TestStructParseEnv2_A", randomA)
		_ = os.Setenv("TestStructParseEnv2_B", randomA)
		test := TestStructParseEnv{}
		err := Parse(&test)
		assert("TestStructParse", test.TestStructParseEnv1.A, randomA)
		assert("TestStructParse", test.TestStructParseEnv1.B, nowUnix)
		assert("TestStructParse", test.TestStructParseEnv1.TestStructParseEnv2.A, randomA)
		assert("TestStructParse", test.TestStructParseEnv1.TestStructParseEnv2.B, nowUnix)
		assert("TestStructParse", test.TestStructParseEnv2.A, randomA)
		assert("TestStructParse", test.TestStructParseEnv2.B, nowUnix)
		assert("TestStructParse", err, nil)
	}
	{
		_ = os.Setenv("RPC", "rpc")
		_ = os.Setenv("RPC_ADDR", "rpc")
		_ = os.Setenv("RPC_USER_NAME", "rpc")
		_ = os.Setenv("mysql_ADDR", "mysql")
		_ = os.Setenv("mysql_PASSWORD", "mysql")
		_ = os.Setenv("mysql_DB", "mysql")
		_ = os.Setenv("mysql_NAME", "mysqlname")
		_ = os.Setenv("RPC_USER_NAME", "rpc")
		_ = os.Setenv("RDS_ADDR", "redis")
		_ = os.Setenv("RDS_PASSWORD", "redis")
		_ = os.Setenv("RDS_DB", "1")
		_ = os.Setenv("RDS-UU_NAME", "redis")
		_ = os.Setenv("PG__ADDR", "postgres")
		_ = os.Setenv("PG__PASSWORD", "postgres")
		_ = os.Setenv("PG__DB", "postgres")
		_ = os.Setenv("PG_USER_NAME", "postgres")
		test := StructConfig{}
		err := Parse(&test)
		assert("TestStructParse", test.RPC.Addr, "rpc")
		assert("TestStructParse", test.RPC.Name, "rpc")
		assert("TestStructParse", test.RPC.User.Name, "rpc")
		assert("TestStructParse", test.MySQL.Addr, "mysql")
		assert("TestStructParse", test.MySQL.Password, "mysql")
		assert("TestStructParse", test.MySQL.DB, "mysql")
		assert("TestStructParse", test.MySQL.User.Name, "mysqlname")
		assert("TestStructParse", test.Redis.Addr, "redis")
		assert("TestStructParse", test.Redis.Password, "redis")
		assert("TestStructParse", test.Redis.DB, 1)
		assert("TestStructParse", test.Redis.User.Name, "redis")
		assert("TestStructParse", test.Postgres.Addr, "postgres")
		assert("TestStructParse", test.Postgres.Password, "postgres")
		assert("TestStructParse", test.Postgres.DB, "postgres")
		assert("TestStructParse", test.Postgres.User.Name, "postgres")
		assert("TestStructParse", err, nil)
	}
}

type StructConfig struct {
	RPC
	MySQL    `env:"mysql"`
	Redis    `env:"RDS"`
	Postgres `env:"PG"`
}

type RPC struct {
	Addr string
	Name string `env:",empty"`
	User
}

type MySQL struct {
	Addr     string
	Password string
	DB       string
	User     `env:",empty"`
}

type Redis struct {
	Addr     string
	Password string
	DB       int
	User     `env:"UU,sep=-"`
}

type Postgres struct {
	Addr     string `env:"ADDR,sep=__"`
	Password string `env:"PASSWORD,sep=__"`
	DB       string `env:"DB,sep=__"`
	User
}

type User struct {
	Name string
}

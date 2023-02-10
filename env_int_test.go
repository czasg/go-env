package env

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"
)

type TestIntParseEnv struct {
	A int   `env:",default=1"`
	B int8  `env:",default=1"`
	C int16 `env:",default=1"`
	D int32 `env:",default=1"`
	E int64 `env:",default=1"`
	F int   `env:",require"`
}

type TestIntParseDefault struct {
	A int   `default:"1"`
	B int8  `default:"1"`
	C int16 `default:"1"`
	D int32 `default:"1"`
	E int64 `default:"xxx"`
}

type TestIntParseEmpty struct {
	A int
	B int8
	C int16
	D int32
	E int64
}

type TestIntParseNotEmpty struct {
	A int
	B int8
	C int16
	D int32
	E int64
}

func TestIntParse(t *testing.T) {
	assert := assertWrap(t)
	nowUnix := time.Now().Unix()
	_ = os.Setenv("A", fmt.Sprintf("%d", nowUnix))
	{
		test := TestIntParseEnv{}
		err := Parse(&test)
		assert("TestIntParse", test.A, int(nowUnix))
		assert("TestIntParse", test.B, int8(1))
		assert("TestIntParse", test.C, int16(1))
		assert("TestIntParse", test.D, int32(1))
		assert("TestIntParse", test.E, int64(1))
		assert("TestIntParse", err, errors.New("F require"))
	}
	{
		test := TestIntParseDefault{}
		err := ParseEntity(Entity{Value: &test, Opt: OptDefault})
		assert("TestIntParse", test.A, 1)
		assert("TestIntParse", test.B, int8(1))
		assert("TestIntParse", test.C, int16(1))
		assert("TestIntParse", test.D, int32(1))
		assert("TestIntParse", test.E, int64(0))
		assert("TestIntParse", err, errors.New("E invalid [xxx]"))
	}
	{
		test := TestIntParseEmpty{}
		err := ParseEntity(Entity{Value: &test})
		assert("TestIntParse", test.A, 0)
		assert("TestIntParse", test.B, int8(0))
		assert("TestIntParse", test.C, int16(0))
		assert("TestIntParse", test.D, int32(0))
		assert("TestIntParse", test.E, int64(0))
		assert("TestIntParse", err, nil)
	}
	{
		test := TestIntParseNotEmpty{
			A: 1,
			B: 1,
			C: 1,
			D: 1,
			E: 1,
		}
		err := Parse(&test)
		assert("TestIntParse", test.A, 1)
		assert("TestIntParse", test.B, int8(1))
		assert("TestIntParse", test.C, int16(1))
		assert("TestIntParse", test.D, int32(1))
		assert("TestIntParse", test.E, int64(1))
		assert("TestIntParse", err, nil)
	}
}

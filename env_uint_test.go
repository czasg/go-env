package env

import (
	"errors"
	"os"
	"testing"
)

type TestUintParseEnv struct {
	A uint   `env:",default=1"`
	B uint8  `env:",default=1"`
	C uint16 `env:",default=1"`
	D uint32 `env:",default=1"`
	E uint64 `env:",default=1"`
	F uint   `env:",require"`
}

type TestUintParseDefault struct {
	A uint   `default:"1"`
	B uint8  `default:"1"`
	C uint16 `default:"1"`
	D uint32 `default:"1"`
	E uint64 `default:"xxx"`
}

type TestUintParseEmpty struct {
	A uint
	B uint8
	C uint16
	D uint32
	E uint64
}

type TestUiIntParseNotEmpty struct {
	A uint
	B uint8
	C uint16
	D uint32
	E uint64
}

func TestUintParse(t *testing.T) {
	assert := assertWrap(t)
	_ = os.Setenv("A", "")
	_ = os.Setenv("B", "")
	_ = os.Setenv("C", "")
	_ = os.Setenv("D", "")
	_ = os.Setenv("E", "")
	{
		test := TestUintParseEnv{}
		err := Parse(&test)
		assert("TestUintParse", test.A, uint(1))
		assert("TestUintParse", test.B, uint8(1))
		assert("TestUintParse", test.C, uint16(1))
		assert("TestUintParse", test.D, uint32(1))
		assert("TestUintParse", test.E, uint64(1))
		assert("TestUintParse", err, errors.New("F require"))
	}
	{
		test := TestUintParseDefault{}
		err := ParseWithOpt(&test, OptDefault)
		assert("TestUintParse", test.A, uint(1))
		assert("TestUintParse", test.B, uint8(1))
		assert("TestUintParse", test.C, uint16(1))
		assert("TestUintParse", test.D, uint32(1))
		assert("TestUintParse", test.E, uint64(0))
		assert("TestUintParse", err, errors.New("E invalid [xxx]"))
	}
	{
		test := TestUintParseEmpty{}
		err := ParseWithOpt(&test)
		assert("TestUintParse", test.A, uint(0))
		assert("TestUintParse", test.B, uint8(0))
		assert("TestUintParse", test.C, uint16(0))
		assert("TestUintParse", test.D, uint32(0))
		assert("TestUintParse", test.E, uint64(0))
		assert("TestUintParse", err, nil)
	}
	{
		test := TestUiIntParseNotEmpty{
			A: 2,
			B: 2,
			C: 2,
			D: 2,
			E: 2,
		}
		err := Parse(&test)
		assert("TestUintParse", test.A, uint(2))
		assert("TestUintParse", test.B, uint8(2))
		assert("TestUintParse", test.C, uint16(2))
		assert("TestUintParse", test.D, uint32(2))
		assert("TestUintParse", test.E, uint64(2))
		assert("TestUintParse", err, nil)
	}
}

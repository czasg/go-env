package env

import (
	"errors"
	"testing"
)

type TestFloatParseEnv struct {
	A float32 `env:",default=1"`
	B float64 `env:",default=1"`
	F float32 `env:",require"`
}

type TestFloatParseDefault struct {
	A float32 `default:"1"`
	B float64 `default:"1"`
	E float32 `default:"xxx"`
}

type TestFloatParseEmpty struct {
	A float32 `default:"1"`
	B float64 `default:"1"`
}

type TestFloatParseNotEmpty struct {
	A float32 `default:"1"`
	B float64 `default:"1"`
}

func TestFloatParse(t *testing.T) {
	assert := assertWrap(t)
	{
		test := TestFloatParseEnv{}
		err := Parse(&test)
		assert("TestFloatParse", test.A, float32(1))
		assert("TestFloatParse", test.B, float64(1))
		assert("TestFloatParse", err, errors.New("F require"))
	}
	{
		test := TestFloatParseDefault{}
		err := Parse(&test)
		assert("TestFloatParse", test.A, float32(1))
		assert("TestFloatParse", test.B, float64(1))
		assert("TestFloatParse", err, errors.New("E invalid [xxx]"))
	}
	{
		test := TestFloatParseEmpty{}
		err := ParseEntity(Entity{Value: &test})
		assert("TestFloatParse", test.A, float32(0))
		assert("TestFloatParse", test.B, float64(0))
		assert("TestFloatParse", err, nil)
	}
	{
		test := TestFloatParseNotEmpty{
			A: 1,
			B: 1,
		}
		err := Parse(&test)
		assert("TestFloatParse", test.A, float32(1))
		assert("TestFloatParse", test.B, float64(1))
		assert("TestFloatParse", err, nil)
	}
}

package env

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"
)

type TestStringParseEnv struct {
	A string `env:"A"`
	B string `env:",default=test"`
	C string `env:",require"`
}

type TestStringParseDefault struct {
	A string `default:"test"`
}

type TestStringParseEmpty struct {
	A string
}

type TestStringParseNotEmpty struct {
	A string
}

func TestStringParse(t *testing.T) {
	assert := assertWrap(t)
	randomA := fmt.Sprintf("%d", time.Now().Unix())
	_ = os.Setenv("A", randomA)
	{
		test := TestStringParseEnv{}
		err := Parse(&test)
		assert("TestStringParse", test.A, randomA)
		assert("TestStringParse", test.B, "test")
		assert("TestStringParse", err, errors.New("C require"))
	}
	{
		test := TestStringParseDefault{}
		err := ParseWithOpt(&test, OptDefault)
		assert("TestStringParse", test.A, "test")
		assert("TestStringParse", err, nil)
	}
	{
		test := TestStringParseEmpty{}
		err := ParseWithOpt(&test)
		assert("TestStringParse", test.A, "")
		assert("TestStringParse", err, nil)
	}
	{
		test := TestStringParseNotEmpty{A: "not empty"}
		err := Parse(&test)
		assert("TestStringParse", test.A, "not empty")
		assert("TestStringParse", err, nil)
	}
}

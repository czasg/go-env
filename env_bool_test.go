package env

import (
	"errors"
	"testing"
)

type TestParseBoolEnv struct {
	A bool
	B bool `env:",default=true"`
	C bool `env:",default=false"`
	D bool `env:",require"`
}

func TestParseBool(t *testing.T) {
	assert := assertWrap(t)
	{
		test := TestParseBoolEnv{C: true}
		err := Parse(&test)
		assert("TestParseBool", test.A, false)
		assert("TestParseBool", test.B, true)
		assert("TestParseBool", test.C, true)
		assert("TestParseBool", err, errors.New("D require"))
	}
}

package env

import (
	"testing"
)

type TestPtrParseNil struct {
	A *string
	B *int
	C map[string]string
	D []string
	E chan struct{}
	F *[2]string
}

type TestPtrParseEnv struct {
	TestPtrParseEnv1 *TestPtrParseEnv1
	TestPtrParseEnv2 *TestPtrParseEnv2
}

type TestPtrParseEnv1 struct {
	A                string `env:",default=test"`
	B                int    `env:",default=1"`
	TestPtrParseEnv2 *TestPtrParseEnv2
}

type TestPtrParseEnv2 struct {
	A string `env:",default=test2"`
	B int    `env:",default=2"`
}

func TestPtrParse(t *testing.T) {
	assert := assertWrap(t)
	{
		test := TestPtrParseNil{}
		err := Parse(&test)
		assert("TestPtrParse", *test.A, "")
		assert("TestPtrParse", *test.B, 0)
		assert("TestPtrParse", test.C, make(map[string]string))
		assert("TestPtrParse", test.D, make([]string, 0, 0))
		assert("TestPtrParse", len(test.E), 0)
		assert("TestPtrParse", err, nil)
	}
	{
		name := "test"
		test := TestPtrParseNil{A: &name}
		err := Parse(&test)
		assert("TestPtrParse", *test.A, "test")
		assert("TestPtrParse", *test.F, [2]string{})
		assert("TestPtrParse", err, nil)
	}
	{
		test := TestPtrParseEnv{}
		err := Parse(&test)
		assert("TestPtrParse", test.TestPtrParseEnv1.A, "")
		assert("TestPtrParse", test.TestPtrParseEnv1.B, 0)
		assert("TestPtrParse", test.TestPtrParseEnv2.A, "")
		assert("TestPtrParse", test.TestPtrParseEnv2.B, 0)
		assert("TestPtrParse", err, nil)
	}
}

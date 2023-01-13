package env

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

type TestString struct {
	Env      string `env:"TEST_STRING_ENV"`
	DF1      string `env:",default=test_string"`
	DF2      string `default:"test_string"`
	Empty    string
	NotEmpty string `default:"invalid"`
}

type TestStringRequire struct {
	Require string `env:"REQUIRE,require"`
}

type TestInt struct {
	A int   `env:"TEST_INT_A" default:"1"`
	B int8  `env:"TEST_INT_B,default=1"`
	C int16 `env:"TEST_INT_C" default:"1"`
	D int32 `env:"TEST_INT_D" default:"1"`
	E int64 `env:"TEST_INT_E" default:"1"`
}

type TestInvalidInt struct {
	A int `env:"ZZ,default=xxx"`
}

type TestStructStringInt struct {
	String TestString
	Int    TestInt
}

type TestNotSupportType struct {
	A interface{}
	B func()
}

func TestParse(t *testing.T) {
	assert := assertWrap(t)
	{
		assert("NotPointerStructErr", Parse(""), NotPointerStructErr)
	}
	{
		testString := TestString{NotEmpty: "NotEmpty"}
		_ = os.Setenv("TEST_STRING_ENV", "test_string")
		assert("String.Parse", Parse(&testString), nil)
		assert("String.Parse.Env", testString.Env, "test_string")
		assert("String.Parse.DF1", testString.DF1, "test_string")
		assert("String.Parse.DF2", testString.DF2, "test_string")
		assert("String.Parse.Empty", testString.Empty, "")
		assert("String.Parse.NotEmpty", testString.NotEmpty, "NotEmpty")
	}
	{
		testStringRequire := TestStringRequire{}
		assert("", Parse(&testStringRequire), errors.New("REQUIRE require"))
	}
	{
		testInt := TestInt{A: 1}
		assert("Int.Parse", Parse(&testInt), nil)
		assert("Int.Parse.A", testInt.A, 1)
		assert("Int.Parse.B", testInt.B, int8(1))
		assert("Int.Parse.C", testInt.C, int16(1))
		assert("Int.Parse.D", testInt.D, int32(1))
		assert("Int.Parse.E", testInt.E, int64(1))
	}
	{
		testInvalidInt := TestInvalidInt{}
		assert("Int.Invalid", Parse(&testInvalidInt), errors.New("A invalid [xxx]"))
	}
	{
		testStructStringInt := TestStructStringInt{}
		assert("Struct.Parse", Parse(&testStructStringInt), nil)
	}
	{
		test := TestNotSupportType{}
		err := Parse(&test)
		assert("Struct.Parse", err, errors.New("unsupport field [A] kind [interface]"))
	}
}

func assertWrap(t *testing.T) func(name string, a, b interface{}) {
	nw := nameWrap()
	return func(name string, a, b interface{}) {
		name = nw(name)
		if !reflect.DeepEqual(a, b) {
			t.Errorf("%s failure! [%v] != [%v]", name, a, b)
		} else {
			t.Logf("%s - pass", name)
		}
	}
}

func nameWrap() func(name string) string {
	count := 0
	lastName := ""
	return func(name string) string {
		if name == lastName {
			count++
		} else {
			count = 0
			lastName = name
		}
		name = strings.ReplaceAll(name, " ", "_")
		return fmt.Sprintf("%s-%d", name, count)
	}
}

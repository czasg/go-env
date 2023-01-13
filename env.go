package env

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Opt int

const (
	OptEnv Opt = 1 << iota
	OptDefault
	OptSilent // ignore err and iterate over all fields.
)

type Payload struct {
	Value       reflect.Value
	Prefix      string
	Opt         Opt
	Field       reflect.Value
	StructField reflect.StructField
}

var (
	NotPointerStructErr = errors.New("only supported pointer to `struct`")
)

func Parse(v interface{}) error {
	return ParseWithOpt(v, OptEnv, OptDefault)
}

func ParseWithOpt(v interface{}, opts ...Opt) error {
	ind := reflect.Indirect(reflect.ValueOf(v))
	if reflect.ValueOf(v).Kind() != reflect.Ptr || ind.Kind() != reflect.Struct {
		return NotPointerStructErr
	}
	var opt Opt = 0
	for _, o := range opts {
		opt = opt | o
	}
	return parse(Payload{Value: ind, Opt: opt})
}

func parse(payload Payload) error {
	for i := 0; i < payload.Value.NumField(); i++ {
		payload.Field = payload.Value.Field(i)
		payload.StructField = payload.Value.Type().Field(i)
		err := parseField(payload)
		if err != nil && payload.Opt&OptSilent != OptSilent {
			return err
		}
	}
	return nil
}

func parseField(payload Payload) error {
	switch payload.Field.Kind() {
	case reflect.Ptr:
		return parsePtr(payload)
	case reflect.Struct:
		return parseStruct(payload)
	case reflect.String:
		return parseString(payload)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return parseInt(payload)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return parseUint(payload)
	case reflect.Float32, reflect.Float64:
		return parseFloat(payload)
	case reflect.Chan:
		return parseChan(payload)
	case reflect.Map:
		return parseMap(payload)
	case reflect.Slice:
		return parseSlice(payload)
	case reflect.Array:
		return parseArray(payload)
	case reflect.Bool:
		return parseBool(payload)
	default:
		return fmt.Errorf("unsupport field [%s] kind [%v]", payload.StructField.Name, payload.Field.Kind())
	}
}

func parseValue(payload Payload) (string, error) {
	var value string
	if payload.Opt&OptEnv == OptEnv {
		envName := strings.ToUpper(payload.StructField.Name)
		envDefault := ""
		envRequire := false
		sep := "_"
		envStr, exist := payload.StructField.Tag.Lookup("env")
		if exist {
			for index, str := range strings.Split(envStr, ",") {
				if index == 0 && str != "" {
					envName = str
				} else if strings.Contains(str, "require") {
					envRequire = true
				} else if strings.Contains(str, "default=") {
					envDefault = strings.TrimPrefix(str, "default=")
				} else if strings.Contains(str, "sep=") {
					sep = strings.TrimPrefix(str, "sep=")
				} else if strings.Contains(str, "empty") {
					envName = ""
				}
			}
		}
		if payload.Prefix != "" {
			if envName == "" {
				envName = payload.Prefix
			} else {
				envName = fmt.Sprintf("%s%s%s", payload.Prefix, sep, envName)
			}
		}
		envValue, exist := os.LookupEnv(envName)
		if !exist && envRequire {
			return "", fmt.Errorf("%s require", envName)
		}
		if envValue == "" {
			envValue = envDefault
		}
		value = envValue
	}
	if value == "" && payload.Opt&OptDefault == OptDefault {
		value = payload.StructField.Tag.Get("default")
	}
	return value, nil
}

func parsePtr(payload Payload) error {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	if payload.Field.IsNil() {
		field := reflect.New(payload.Field.Type().Elem())
		payload.Field.Set(field)
		return nil
	}
	return nil
}

/*
`env:"field,sep=_,default=df,require,empty"`
*/
func parseStruct(payload Payload) error {
	payload.Value = payload.Field
	fieldName := strings.ToUpper(payload.StructField.Name)
	sep := "_"
	envStr, exist := payload.StructField.Tag.Lookup("env")
	if exist {
		for index, str := range strings.Split(envStr, ",") {
			if index == 0 && str != "" {
				fieldName = str
			} else if strings.Contains(str, "sep=") {
				sep = strings.TrimPrefix(str, "sep=")
			} else if strings.Contains(str, "empty") {
				fieldName = ""
			}
		}
	}
	if payload.Prefix == "" {
		payload.Prefix = fieldName
		return parse(payload)
	}
	if fieldName == "" {
		return parse(payload)
	}
	payload.Prefix = fmt.Sprintf("%s%s%s", payload.Prefix, sep, fieldName)
	return parse(payload)
}

func parseString(payload Payload) error {
	if payload.Field.String() != "" {
		return nil
	}
	value, err := parseValue(payload)
	if err != nil {
		return err
	}
	if value == "" {
		return nil
	}
	payload.Field.SetString(value)
	return nil
}

func parseInt(payload Payload) error {
	if payload.Field.Int() != 0 {
		return nil
	}
	value, err := parseValue(payload)
	if err != nil {
		return err
	}
	if value == "" {
		return nil
	}
	iv, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("%s invalid [%s]", payload.StructField.Name, value)
	}
	payload.Field.SetInt(int64(iv))
	return nil
}

func parseUint(payload Payload) error {
	if payload.Field.Uint() != 0 {
		return nil
	}
	value, err := parseValue(payload)
	if err != nil {
		return err
	}
	if value == "" {
		return nil
	}
	iv, err := strconv.ParseUint(value, 0, 64)
	if err != nil {
		return fmt.Errorf("%s invalid [%s]", payload.StructField.Name, value)
	}
	payload.Field.SetUint(iv)
	return nil
}

func parseFloat(payload Payload) error {
	if payload.Field.Float() != 0 {
		return nil
	}
	value, err := parseValue(payload)
	if err != nil {
		return err
	}
	if value == "" {
		return nil
	}
	iv, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("%s invalid [%s]", payload.StructField.Name, value)
	}
	payload.Field.SetFloat(iv)
	return nil
}

func parseChan(payload Payload) error {
	if payload.Field.IsNil() {
		payload.Field.Set(reflect.MakeChan(payload.Field.Type(), 0))
	}
	return nil
}

func parseMap(payload Payload) error {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	if payload.Field.IsNil() {
		payload.Field.Set(reflect.MakeMap(payload.Field.Type()))
	}
	return nil
}

func parseSlice(payload Payload) error {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	if payload.Field.IsNil() {
		payload.Field.Set(reflect.MakeSlice(payload.Field.Type(), 0, 0))
	}
	return nil
}

func parseArray(_ Payload) error {
	return nil
}

func parseBool(payload Payload) error {
	if payload.Field.Bool() {
		return nil
	}
	value, err := parseValue(payload)
	if err != nil {
		return err
	}
	if value == "" {
		return nil
	}
	b, _ := strconv.ParseBool(value)
	payload.Field.SetBool(b)
	return nil
}

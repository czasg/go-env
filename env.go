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

func (o Opt) Enable(opt Opt) bool {
	return o&opt == opt
}

const (
	OptEnv Opt = 1 << iota
	OptDefault
	OptSilent // ignore err and iterate over all fields.
)

type Entity struct {
	Value  interface{}
	Prefix string
	Opt    Opt
}

type Payload struct {
	Value       reflect.Value
	Prefix      string
	Opt         Opt
	Field       reflect.Value
	StructField reflect.StructField
}

var NotPointerStructErr = errors.New("only supported pointer to `struct`")

func Parse(v interface{}) error {
	return ParseEntity(Entity{
		Value: v,
		Opt:   OptEnv | OptDefault,
	})
}

func ParseEntity(e Entity) error {
	ind := reflect.Indirect(reflect.ValueOf(e.Value))
	if reflect.ValueOf(e.Value).Kind() != reflect.Ptr || ind.Kind() != reflect.Struct {
		return NotPointerStructErr
	}
	return parse(Payload{Value: ind, Prefix: e.Prefix, Opt: e.Opt})
}

func parse(p Payload) error {
	for i := 0; i < p.Value.NumField(); i++ {
		p.Field = p.Value.Field(i)
		p.StructField = p.Value.Type().Field(i)
		err := parseField(p)
		if err != nil && !p.Opt.Enable(OptSilent) {
			return err
		}
	}
	return nil
}

func parseField(p Payload) error {
	switch p.Field.Kind() {
	case reflect.Ptr:
		return parsePtr(p)
	case reflect.Struct:
		return parseStruct(p)
	case reflect.String:
		return parseString(p)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return parseInt(p)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return parseUint(p)
	case reflect.Float32, reflect.Float64:
		return parseFloat(p)
	case reflect.Chan:
		return parseChan(p)
	case reflect.Map:
		return parseMap(p)
	case reflect.Slice:
		return parseSlice(p)
	case reflect.Array:
		return parseArray(p)
	case reflect.Bool:
		return parseBool(p)
	default:
		return fmt.Errorf("unsupport field [%s] kind [%v]", p.StructField.Name, p.Field.Kind())
	}
}

func parseValue(p Payload) (string, error) {
	var value string
	if p.Opt.Enable(OptEnv) {
		envName := strings.ToUpper(p.StructField.Name)
		envDefault := ""
		envRequire := false
		sep := "_"
		envStr, exist := p.StructField.Tag.Lookup("env")
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
		if p.Prefix != "" {
			if envName == "" {
				envName = p.Prefix
			} else {
				envName = fmt.Sprintf("%s%s%s", p.Prefix, sep, envName)
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
	if value == "" && p.Opt.Enable(OptDefault) {
		value = p.StructField.Tag.Get("default")
	}
	return value, nil
}

func parsePtr(p Payload) error {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	if p.Field.IsNil() {
		field := reflect.New(p.Field.Type().Elem())
		p.Field.Set(field)
		return nil
	}
	return nil
}

/*
`env:"field,sep=_,default=df,require,empty"`
*/
func parseStruct(p Payload) error {
	p.Value = p.Field
	fieldName := strings.ToUpper(p.StructField.Name)
	sep := "_"
	envStr, exist := p.StructField.Tag.Lookup("env")
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
	if p.Prefix == "" {
		p.Prefix = fieldName
		return parse(p)
	}
	if fieldName == "" {
		return parse(p)
	}
	p.Prefix = fmt.Sprintf("%s%s%s", p.Prefix, sep, fieldName)
	return parse(p)
}

func parseString(p Payload) error {
	if p.Field.String() != "" {
		return nil
	}
	value, err := parseValue(p)
	if err != nil {
		return err
	}
	if value == "" {
		return nil
	}
	p.Field.SetString(value)
	return nil
}

func parseInt(p Payload) error {
	if p.Field.Int() != 0 {
		return nil
	}
	value, err := parseValue(p)
	if err != nil {
		return err
	}
	if value == "" {
		return nil
	}
	iv, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("%s invalid [%s]", p.StructField.Name, value)
	}
	p.Field.SetInt(int64(iv))
	return nil
}

func parseUint(p Payload) error {
	if p.Field.Uint() != 0 {
		return nil
	}
	value, err := parseValue(p)
	if err != nil {
		return err
	}
	if value == "" {
		return nil
	}
	iv, err := strconv.ParseUint(value, 0, 64)
	if err != nil {
		return fmt.Errorf("%s invalid [%s]", p.StructField.Name, value)
	}
	p.Field.SetUint(iv)
	return nil
}

func parseFloat(p Payload) error {
	if p.Field.Float() != 0 {
		return nil
	}
	value, err := parseValue(p)
	if err != nil {
		return err
	}
	if value == "" {
		return nil
	}
	iv, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("%s invalid [%s]", p.StructField.Name, value)
	}
	p.Field.SetFloat(iv)
	return nil
}

func parseChan(p Payload) error {
	if p.Field.IsNil() {
		p.Field.Set(reflect.MakeChan(p.Field.Type(), 0))
	}
	return nil
}

func parseMap(p Payload) error {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	if p.Field.IsNil() {
		p.Field.Set(reflect.MakeMap(p.Field.Type()))
	}
	return nil
}

func parseSlice(p Payload) error {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	if p.Field.IsNil() {
		p.Field.Set(reflect.MakeSlice(p.Field.Type(), 0, 0))
	}
	return nil
}

func parseArray(_ Payload) error {
	return nil
}

func parseBool(p Payload) error {
	if p.Field.Bool() {
		return nil
	}
	value, err := parseValue(p)
	if err != nil {
		return err
	}
	if value == "" {
		return nil
	}
	b, _ := strconv.ParseBool(value)
	p.Field.SetBool(b)
	return nil
}

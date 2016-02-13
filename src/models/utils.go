package models

import (
	"bytes"
	"io"
	"strings"

    "database/sql/driver"
    "encoding/json"

    "github.com/satori/go.uuid"
    "github.com/golang/glog"
)

func NewInterfaceMap() map[string]interface{} {
	return make(map[string]interface{})
}

type InterfaceMap map[string]interface{}


func (m *InterfaceMap) Scan(value interface{}) error {
    return json.Unmarshal(value.([]byte), m)
}

func (m InterfaceMap) Value() (driver.Value, error) {
	if len(m) == 0 {
		return string("{}"), nil
	}
	
    b, err := json.Marshal(m)
    if err != nil {
        return string("{}"), err
    }
    return string(b), nil
}

type StringMap map[string]string

func NewStringMap() map[string]string {
	return make(map[string]string)
}

func (m *StringMap) Scan(value interface{}) error {
    return json.Unmarshal(value.([]byte), m)
}

func (m StringMap) Value() (driver.Value, error) {
	if len(m) == 0 {
		return string("{}"), nil
	}
	
    b, err := json.Marshal(m)
    if err != nil {
        return string("{}"), err
    }
    return string(b), nil
}

type StringArray []string

func (c *StringArray) FromArray(v interface{}) *StringArray {
	switch v.(type) {
	case []uuid.UUID:
		for _, _v := range v.([]uuid.UUID) {
			c.Add(_v.String())
		}
	case []string:
		for _, _v := range v.([]string) {
			c.Add(_v)
		}
	case StringArray:
		for _, _v := range []string(v.(StringArray)) {
			c.Add(_v)
		}
	default:
		glog.Warningf("Not supported type %T", v)
	}

	return c
}

func (c StringArray) Array() []string {
	return []string(c)
}

func (c *StringArray) AddAsArray(in []string) {
	for _, str := range in {
		c.Add(str)
	}
}

func (c *StringArray) AddAsFields(in string) {
	for _, str := range strings.Fields(in) {
		c.Add(str)
	}
}

func (c *StringArray) Del(str string) *StringArray {
	str = strings.TrimSpace(str)

	if !c.IsExist(str) {
		return c
	}

	for index, value := range *c {

		if bytes.Equal([]byte(str), []byte(value)) {
			(*c)[index] = (*c)[len((*c))-1]
			(*c) = (*c)[:len((*c))-1]

			return c
		}
	}

	return c
}

func (c *StringArray) IsExist(str string) bool {
	str = strings.TrimSpace(str)

	for _, value := range *c {

		if bytes.Equal([]byte(str), []byte(value)) {

			return true
		}
	}

	return false
}

func (c *StringArray) Add(str string) *StringArray {
	str = strings.TrimSpace(str)

	if len(str) == 0 {
		return c
	}

	if !c.IsExist(str) {
		*c = append(*c, str)
	}

	return c
}

// extractFieldsFromMap if len(without) == 0 return all fileds
func ExtractFieldsFromMap(m map[string]interface{}, without ...string) (keys []string, fields []interface{}) {
	_without := make(map[string]bool)
	var flagAllFields = len(without) == 0

	for _, v := range without {
		_without[v] = true
	}

	for fieldName, field := range m {
		if !flagAllFields {
			if !_without[fieldName] {
				continue
			}
		}

		keys = append(keys, fieldName)
		fields = append(fields, field)
	}

	return
}

// FromJson extract object from data (io.Reader OR []byte)
func FromJson(obj interface{}, data interface{}) error {
	switch data.(type) {
	case io.Reader:
		decoder := json.NewDecoder(data.(io.Reader))
		return decoder.Decode(obj)
	case []byte:
		return json.Unmarshal(data.([]byte), obj)
	}

	return ErrNotSupported
}

package store

import (
	"testing"
	"models"
	"github.com/satori/go.uuid"
	"strings"
)

func TestUpdateValue(t *testing.T) {
	key1 := uuid.NewV1().String()
	key2 := uuid.NewV1().String()
	key3 := uuid.NewV1().String()


	dto := models.NewValueDTO()
	dto.Keys.Add(key1)
	dto.Keys.Add(key2)
	dto.Keys.Add(key3)
	dto.Value = strings.Join([]string{key1, key2, key3}, ":")

	value , err := _s.Get("value").(*ValueStore).Create(dto)

	if err != nil {
		t.Error(err)
		return
	}

	dto.ValueId = value.PrimaryValue().String()
	dto.Value = "new value"
	dto.UpdateFields = []string{"value"}
	_, err = _s.Get("value").(*ValueStore).Update(dto)

	if err != nil {
		t.Error(err)
		return
	}

	dto.Keys = models.StringArray{}
	dto.Keys.Add(key3)
	dto.Keys.Add(key1)
	dto.Keys.Add(key2)
	values, err := _s.Get("value").(*ValueStore).FindByKeys(dto, "=")

	if err != nil {
		t.Error(err)
		return
	}

	if len(values) != 1 {
		t.Error("count of values ​​is not expected")
		return
	}

	if values[0].Value != dto.Value {
		t.Error("is not expected 'value'")
		return	
	}
}

func TestCreateUniqValue(t *testing.T) {
	key1 := uuid.NewV1().String()
	key2 := uuid.NewV1().String()
	key3 := uuid.NewV1().String()


	dto := models.NewValueDTO()
	dto.Keys.Add(key1)
	dto.Keys.Add(key2)
	dto.Keys.Add(key3)
	dto.Keys.Add("uniq")
	dto.Value = strings.Join([]string{key1, key2, key3}, ":")

	_ , err := _s.Get("value").(*ValueStore).Create(dto)

	if err != nil {
		t.Error(err)
		return
	}

	_ , err = _s.Get("value").(*ValueStore).Create(dto)

	if err == nil || !strings.HasPrefix(err.Error(), "ERROR: duplicate key value violates unique constraint \"values_keys_ifuniq_idx\""){
		t.Error(err)
		return
	}
}

func TestCreateValue(t *testing.T) {
	key1 := uuid.NewV1().String()
	key2 := uuid.NewV1().String()
	key3 := uuid.NewV1().String()


	dto := models.NewValueDTO()
	dto.Keys.Add(key1)
	dto.Keys.Add(key2)
	dto.Keys.Add(key3)
	dto.Value = strings.Join([]string{key1, key2, key3}, ":")

	_ , err := _s.Get("value").(*ValueStore).Create(dto)

	if err != nil {
		t.Error(err)
		return
	}

	dto.ValueId = ""
	values, err := _s.Get("value").(*ValueStore).FindByKeys(dto, "=")

	if err != nil {
		t.Error(err)
		return
	}

	value := values[0]
	_keys := models.StringArray(value.Keys)

	if !(_keys.IsExist(key1) && _keys.IsExist(key2) && _keys.IsExist(key3)) {
		t.Error("not valid keys")
		return
	}

	// 

	dto.ValueId = ""
	dto.Keys = models.StringArray{}
	dto.Keys.Add(key1)
	values, err = _s.Get("value").(*ValueStore).FindByKeys(dto, "=")

	if err != nil {
		t.Error(err)
		return
	}

	if len(values) != 0 {
		t.Error("error search mode =")
		return
	}

	// 

	dto.ValueId = ""
	dto.Keys = models.StringArray{}
	dto.Keys.Add(key1)
	values, err = _s.Get("value").(*ValueStore).FindByKeys(dto, "&&")

	if err != nil {
		t.Error(err)
		return
	}

	if len(values) != 1 {
		t.Error("error search mode &&")
		return
	}

	value = values[0]
	_keys = models.StringArray(value.Keys)

	if !(_keys.IsExist(key1) && _keys.IsExist(key2) && _keys.IsExist(key3)) {
		t.Error("not valid keys")
		return
	}
}
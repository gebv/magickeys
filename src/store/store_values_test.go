package store

import (
	"testing"
	"models"
	"github.com/satori/go.uuid"
	"strings"
	"strconv"
)

func TestModeEventsCreated(t *testing.T) {
	key1 := uuid.NewV1().String()
	key2 := uuid.NewV1().String()

	var origValueIndexes = make(map[uuid.UUID]int)

	for i, _ := range strings.Repeat("+", 10) {
		dto := models.NewValueDTO()
		dto.Keys.Add(key2)
		dto.Keys.Add(key1)
		dto.Value = strings.Join([]string{key1, key2, strconv.Itoa(i)}, ":")
		dto.Props["iter"] = i
		dto.Props["iter_inc"] = i+1
		dto.Props[strconv.Itoa(i)+":special_value"] = i+2
		dto.Flags.Add("tag:"+strconv.Itoa(i))
		dto.IsEnabled = true

		value , err := _s.Get("value").(*ValueStore).Create(dto)

		if err != nil {
			t.Error(err)
			return 
		}

		origValueIndexes[value.PrimaryValue()] = i
	}

	dto := models.NewValueDTO()
	dto.Keys.Add(key2)
	dto.Keys.Add(key1)
	values, err := _s.Get("value").(*ValueStore).FindByKeys(dto, "=")

	if err != nil {
		t.Error(err)
		return
	}
	
	for _, value := range values {
		iterValue, exist := origValueIndexes[value.PrimaryValue()]

		if !exist {
			t.Fatal("not exist item")
			return	
		}

		if value.Value != strings.Join([]string{key1, key2, strconv.Itoa(iterValue)}, ":") {
			t.Errorf("values Value is not correct for %v", value.PrimaryValue().String())
			return	
		}	

		if int(value.Props["iter"].(float64)) != iterValue {
			t.Errorf("values Props[\"iter\"] is not correct for %v", value.PrimaryValue().String())
			return	
		}

		if int(value.Props["iter_inc"].(float64)) != iterValue+1 {
			t.Errorf("values Props[\"iter_inc\"] is not correct for %v", value.PrimaryValue().String())
			return	
		}

		if int(value.Props[strconv.Itoa(iterValue)+":special_value"].(float64)) != iterValue+2 {
			t.Errorf("values Props[\"#:special_value\"] is not correct for %v", value.PrimaryValue().String())
			return	
		}

		_tags := models.StringArray{}
		_tags.FromArray(value.Flags)

		if !_tags.IsExist("tag:"+strconv.Itoa(iterValue)) {
			t.Errorf("values Flags is not correct for %v", value.PrimaryValue().String())
			return		
		}

		if value.IsEnabled != true {
			t.Errorf("values IsEnabled is not correct for %v", value.PrimaryValue().String())
			return			
		}
	}
}

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
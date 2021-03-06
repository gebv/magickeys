// Code generated. DO NOT EDIT.
package models

import (
	"github.com/golang/glog"
	"github.com/satori/go.uuid"
	"time"
)

// ValueDTO
func NewValueDTO() *ValueDTO {
	model := new(ValueDTO)
	// Custom factory code
	model.Value = NewInterfaceMap()
	return model
}

type ValueDTO struct {
	DTO
	// ValueId
	ValueId string `json:"-" `
	// UpdateFields	Какие поля следует
	UpdateFields StringArray `json:"-" `
	// Keys
	Keys StringArray `json:"keys" `
	// Value
	Value InterfaceMap `json:"value" `
	// IsRemoved
	IsRemoved bool `json:"is_removed" `
}

func (model ValueDTO) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *ValueDTO) TransformFrom(in interface{}) error {
	switch in.(type) {
	default:
		glog.Errorf("Not supported type %v", in)
		return ErrNotSupported
	}
	return nil

}

//
// Helpful functions
//

func (v *ValueDTO) Maps() map[string]interface{} {
	return map[string]interface{}{
		// ValueId
		"value_id": &v.ValueId,
		// UpdateFields	Какие поля следует
		"update_fields": &v.UpdateFields,
		// Keys
		"keys": &v.Keys,
		// Value
		"value": &v.Value,
		// IsRemoved
		"is_removed": &v.IsRemoved,
	}
}

// Fields extract of fields from map
func (v *ValueDTO) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(v.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (v *ValueDTO) FromJson(data interface{}) error {
	return FromJson(v, data)
}

// Value
func NewValue() *Value {
	model := new(Value)
	// Custom factory code
	model.Value = NewInterfaceMap()
	return model
}

type Value struct {
	// ValueId
	ValueId uuid.UUID `json:"value_id" `
	// Keys
	Keys []string `json:"keys" `
	// Value
	Value InterfaceMap `json:"value" `
	// IsRemoved
	IsRemoved bool `json:"is_removed" `
	// UpdatedAt
	UpdatedAt time.Time `json:"updated_at" sql:"type:timestamp;default:now()" `
	// CreatedAt
	CreatedAt time.Time `json:"created_at" sql:"type:timestamp;default:null" `
}

func (model Value) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *Value) TransformFrom(in interface{}) error {
	switch in.(type) {
	case *Value:
		dto := in.(*Value)
		model.CreatedAt = dto.CreatedAt
		model.ValueId = dto.ValueId
		model.Keys = dto.Keys
		model.Value = dto.Value
		model.IsRemoved = dto.IsRemoved
		model.UpdatedAt = dto.UpdatedAt
	case *ValueDTO:
		dto := in.(*ValueDTO)
		model.IsRemoved = dto.IsRemoved
		model.Keys = dto.Keys.Array()
		model.Value = dto.Value
		model.ValueId = uuid.FromStringOrNil(dto.ValueId)
	default:
		glog.Errorf("Not supported type %v", in)
		return ErrNotSupported
	}
	return nil

}

//
// Helpful functions
//

func (v *Value) Maps() map[string]interface{} {
	return map[string]interface{}{
		// ValueId
		"value_id": &v.ValueId,
		// Keys
		"keys": &v.Keys,
		// Value
		"value": &v.Value,
		// IsRemoved
		"is_removed": &v.IsRemoved,
		// UpdatedAt
		"updated_at": &v.UpdatedAt,
		// CreatedAt
		"created_at": &v.CreatedAt,
	}
}

// Fields extract of fields from map
func (v *Value) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(v.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (v *Value) FromJson(data interface{}) error {
	return FromJson(v, data)
}

func (Value) TableName() string {
	return "values"
}

// PrimaryName primary field name
func (Value) PrimaryName() string {
	return "value_id"
}

// PrimaryValue primary value
func (v Value) PrimaryValue() uuid.UUID {
	return v.ValueId
}

// model

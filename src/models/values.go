package models

import (
	"time"
)

func (v *Value) BeforeCreate() {
	v.CreatedAt = time.Now()
}

func (v *Value) BeforeSave() {
	v.UpdatedAt = time.Now()	
}

func (v *Value) BeforeDelete() {
	v.UpdatedAt = time.Now()
	v.IsRemoved = true
}
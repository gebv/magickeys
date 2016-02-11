package models

import (
	"github.com/satori/go.uuid"
	"github.com/jackc/pgx"
)

type Model interface {
	TransformTo(interface{}) error
	TransformFrom(interface{}) error

	Maps() map[string]interface{}
	Fields(fields ...string) ([]string, []interface{})
	FromJson(data interface{}) error

	BeforeCreate()
	BeforeSave()
	BeforeDelete()

	PrimaryName() string
	PrimaryValue() uuid.UUID
	TableName() string
}

type DTO struct {
	Tx *pgx.Tx `json:"-" `
}
// Code generated. DO NOT EDIT.
package models

import (
	"github.com/golang/glog"
)

// Response
func NewResponse() *Response {
	model := new(Response)
	// Custom factory code
	model.StatusCode = 400 // http.StatusBadRequest
	return model
}

type Response struct {
	// StatusCode
	StatusCode int `json:"status_code" `
	// Message
	Message string `json:"message,omitempty" `
	// DevMessage
	DevMessage string `json:"dev_message,omitempty" `
	// Data
	Data interface{} `json:"data,omitempty" `
}

func (model Response) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *Response) TransformFrom(in interface{}) error {
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

func (r *Response) Maps() map[string]interface{} {
	return map[string]interface{}{
		// StatusCode
		"status_code": &r.StatusCode,
		// Message
		"message": &r.Message,
		// DevMessage
		"dev_message": &r.DevMessage,
		// Data
		"data": &r.Data,
	}
}

// Fields extract of fields from map
func (r *Response) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(r.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (r *Response) FromJson(data interface{}) error {
	return FromJson(r, data)
}

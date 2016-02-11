package models

import (
	"encoding/json"
	"github.com/golang/glog"
)

// Response

func (c *Response) ToJson() []byte {
	b, err := json.Marshal(c)

	if err != nil {
		glog.Errorf("Marshal ResponseDTO error, %s", err)
		return []byte(`{}`)
	}

	return b
}

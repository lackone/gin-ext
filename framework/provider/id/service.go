package id

import "github.com/rs/xid"

type ExtId struct {
}

func NewExtId(params ...interface{}) (interface{}, error) {
	return &ExtId{}, nil
}

func (e *ExtId) NewID() string {
	return xid.New().String()
}

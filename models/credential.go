package models

import (
	"encoding/json"
	"time"
)

type Credential struct {
	ID       int64
	Name     string `xorm:"index"`
	Hash     string `xorm:"index"`
	Pass     string `xorm:"index"`
	Additive json.RawMessage
	Created  time.Time `xorm:"CREATED"`
	Updated  time.Time `xorm:"UPDATED"`
	Deleted  time.Time `xorm:"deleted"`
}

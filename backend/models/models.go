package models

import "time"

type Session struct {
	Userid string
	Expiry time.Time
}

package models

import "time"

type SolutionStatusHelper struct {
	SolutionId uint64 `orm:"pk"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
	StatusHelper string `orm:"type(text)"`
}

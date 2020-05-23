package models

import "time"

const SolutionCodeTableName = "solution_codes"

type SolutionCode struct {
	SolutionId uint64 `orm:"pk"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
	Code string `orm:"type(text)"`
}

func (sc *SolutionCode) TableName() string {
	return SolutionCodeTableName
}
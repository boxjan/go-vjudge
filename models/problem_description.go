package models

import "time"

const ProblemDescriptionTableName = "problem_descriptions"

type ProblemDescription struct {
	ProblemId uint64 `orm:"pk" json:"problem_id"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)" json:"updated_at"`
	Description string `orm:"type(text)" json:"description"`
	Input string `orm:"type(text)" json:"input"`
	Output string `orm:"type(text)" json:"output"`
	Sample string `orm:"type(text)" json:"sample"`
	Hint string `orm:"type(text)" json:"hint"`
}

func (pd *ProblemDescription)TableName() string {
	return ProblemDescriptionTableName
}

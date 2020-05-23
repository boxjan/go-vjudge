package models

import "time"

const ProblemTableName = "problems"

type Problem struct {
	Id uint64 `orm:"pk;auto" json:"id"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)" json:"updated_at"`
	OriOj string `orm:"size(32)" json:"ori_oj"`
	OriProblemId string `orm:"size(32)" json:"ori_id"`
	OriUrl string `orm:"size(512)" json:"ori_url"`
	Title string `orm:"size(128)" json:"title"`
	Source string `orm:"size(128)" json:"source"`
	MemoryLimit uint `json:"memory_limit"`// limit in kb
	TimeLimit uint `json:"time_limit"`// limit in ms
}

func (p *Problem) TableName() string {
	return ProblemTableName
}

func (p *Problem) TableUnique() [][]string {
	return [][]string{
		[]string{"OriOj", "OriProblemId"},
	}
}

func (p *Problem) GetProblemDesc() (*ProblemDescription, error) {
	res := &ProblemDescription{ProblemId: p.Id}
	err := o.Read(res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *Problem) InsertProblemDesc(pd *ProblemDescription) error {
	pd.ProblemId = p.Id
	_, err := o.Insert(pd)
	return err
}

func (p *Problem) UpdateProblemDesc(pd *ProblemDescription) error {
	pd.ProblemId = p.Id
	_, err := o.Update(pd)
	return err
}
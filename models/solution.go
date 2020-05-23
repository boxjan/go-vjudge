package models

import (
	"time"
)

const SolutionTableName = "solutions"

type Solution struct {
	Id uint64 `orm:"pk;auto"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
	UserId uint64
	ProblemId uint64
	ContestId uint64 `orm:"null"`// unused now
	OriOj string `orm:"size(32)"`
	OriProblemId string `orm:"size(32)"`
	OriSolutionId string `orm:"size(32)"`
	OriAccount string `orm:"size(64)"`
	Language string `orm:"size(64)"`
	Length uint
	Hash string `orm:"size(128)"`
	RawStatus string `orm:"size(255)"`
	Status int
	TimeUsed uint  `orm:"null"` //ms
	MemoryUsed uint  `orm:"null"` // kb
	QueryCount uint  `orm:"null"`
	OriSubmitTime time.Time `orm:"type(datetime)"`
}

func (s *Solution)TableName() string {
	return SolutionTableName
}

func (s *Solution) GetCode() (*SolutionCode, error) {
	res := SolutionCode{SolutionId: s.Id}
	err := o.Read(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *Solution) GetStatusHelper() (*SolutionStatusHelper, error) {
	res := SolutionStatusHelper{SolutionId: s.Id}
	err := o.Read(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *Solution) InsertCode(sc *SolutionCode) error {
	sc.SolutionId = s.Id
	_, err := o.Insert(sc)
	return err
}

func (s *Solution) InsertStatusHelper(ssh *SolutionStatusHelper) error {
	ssh.SolutionId = s.Id
	_, err := o.Insert(ssh)
	return err
}

func (s *Solution) UpdateCode(sc *SolutionCode) (int64, error) {
	sc.SolutionId = s.Id
	return o.Update(sc)
}

func (s *Solution) UpdateStatus(ssh *SolutionStatusHelper) (int64, error) {
	ssh.SolutionId = s.Id
	return o.Update(ssh)
}


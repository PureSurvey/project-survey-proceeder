package objects

import "time"

type Survey struct {
	Id          int
	Name        string
	UserId      int
	TargetingId int

	dateBy time.Time
}

func NewSurvey(id int, name string, dateBy time.Time, userId int, targetingId int) *Survey {
	return &Survey{
		Id:          id,
		Name:        name,
		dateBy:      dateBy,
		UserId:      userId,
		TargetingId: targetingId,
	}
}

func (s *Survey) IsActiveOnDate(date time.Time) bool {
	return s.dateBy.After(date)
}

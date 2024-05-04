package objects

import "time"

type Survey struct {
	Id          int
	Name        string
	DateBy      time.Time
	UserId      int
	TargetingId int
}

func NewSurvey(id int, name string, dateBy time.Time, userId int, targetingId int) *Survey {
	return &Survey{
		Id:          id,
		Name:        name,
		DateBy:      dateBy,
		UserId:      userId,
		TargetingId: targetingId,
	}
}

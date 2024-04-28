package objects

import "time"

type Survey struct {
	Id          int
	Name        int
	DateBy      time.Time
	UserId      int
	TargetingId int
}

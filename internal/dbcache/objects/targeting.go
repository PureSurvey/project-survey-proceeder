package objects

type Targeting struct {
	Id int
}

func NewTargeting(id int) *Targeting {
	return &Targeting{Id: id}
}

package objects

type Option struct {
	Id          int
	QuestionId  int
	OrderNumber int
}

func NewOption(id int, questionId int, orderNumber int) *Option {
	return &Option{
		Id:          id,
		QuestionId:  questionId,
		OrderNumber: orderNumber,
	}
}

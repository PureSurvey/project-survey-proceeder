package objects

type Translation struct {
	Id             int
	Translation    string
	Language       string
	ParentId       int
	QuestionLineId int
}

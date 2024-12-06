package contracts

type ISurveyMarkupService interface {
	Init() error
	GetMarkup(unitId int, surveyIds []int, language string) (string, error)
	Close() error
}

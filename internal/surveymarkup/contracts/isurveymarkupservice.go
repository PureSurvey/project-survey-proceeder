package contracts

type ISurveyMarkupService interface {
	GetMarkup(unitId int, surveyIds []int, language string) (string, error)
}

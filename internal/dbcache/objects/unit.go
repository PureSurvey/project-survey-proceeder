package objects

type Unit struct {
	Id           int
	Name         string
	UserId       int
	AppearanceId int

	SurveyTakesPerDevice    int
	MaximumSurveysPerDevice int
	HideAfterNoSurveys      bool
	MessageAfterNoSurveys   string
}

func NewUnit(id int, name string, userId int, appearanceId int,
	surveyTakesPerDevice int, maxSurveysPerDevice int, hideAfterNoSurveys bool,
	message string) *Unit {
	return &Unit{
		Id:                      id,
		Name:                    name,
		UserId:                  userId,
		AppearanceId:            appearanceId,
		SurveyTakesPerDevice:    surveyTakesPerDevice,
		MaximumSurveysPerDevice: maxSurveysPerDevice,
		HideAfterNoSurveys:      hideAfterNoSurveys,
		MessageAfterNoSurveys:   message,
	}
}

package objects

type Unit struct {
	Id           int
	Name         string
	UserId       int
	AppearanceId int

	OneSurveyPerDevice      bool
	MaximumSurveysPerDevice int
	HideAfterNoSurveys      bool
	MessageAfterNoSurveys   string
}

func NewUnit(id int, name string, userId int, appearanceId int,
	oneSurveyPerDevice bool, maxSurveysPerDevice int, hideAfterNoSurveys bool,
	message string) *Unit {
	return &Unit{
		Id:                      id,
		Name:                    name,
		UserId:                  userId,
		AppearanceId:            appearanceId,
		OneSurveyPerDevice:      oneSurveyPerDevice,
		MaximumSurveysPerDevice: maxSurveysPerDevice,
		HideAfterNoSurveys:      hideAfterNoSurveys,
		MessageAfterNoSurveys:   message,
	}
}

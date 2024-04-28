package objects

import "project-survey-proceeder/enums"

type Appearance struct {
	Id         int
	Type       enums.EnumAppearanceType
	TemplateId int
	Params     string
}

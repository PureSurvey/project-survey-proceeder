package objects

import "project-survey-proceeder/enums"

type Appearance struct {
	Id         int
	Type       enums.EnumAppearanceType
	TemplateId int
	Params     string
}

func NewAppearance(id int, aType enums.EnumAppearanceType, templateId int, params string) *Appearance {
	return &Appearance{
		Id:         id,
		Type:       aType,
		TemplateId: templateId,
		Params:     params,
	}
}

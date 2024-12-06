package objects

import (
	"project-survey-proceeder/internal/enums"
)

type Appearance struct {
	Id         int
	Type       enums.EnumAppearanceType
	TemplateId int
}

func NewAppearance(id int, aType enums.EnumAppearanceType, templateId int) *Appearance {
	return &Appearance{
		Id:         id,
		Type:       aType,
		TemplateId: templateId,
	}
}

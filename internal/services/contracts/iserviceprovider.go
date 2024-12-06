package contracts

import (
	contextcontracts "project-survey-proceeder/internal/context/contracts"
	"project-survey-proceeder/internal/dbcache"
	eventcontracts "project-survey-proceeder/internal/events/contracts"
	"project-survey-proceeder/internal/respondents"
	surveymarkupcontracts "project-survey-proceeder/internal/surveymarkup/contracts"
	"project-survey-proceeder/internal/targeting/contracts"
)

type IServiceProvider interface {
	GetUnitContextFiller() contextcontracts.IRequestFiller
	GetEventContextFiller() contextcontracts.IRequestFiller
	GetDbRepo() *dbcache.Repo
	GetTargetingService() contracts.ITargetingService
	GetSurveyMarkupService() surveymarkupcontracts.ISurveyMarkupService
	GetEventProducer() eventcontracts.IEventProducer
	GetRespondentsService() *respondents.Service

	Dispose() error
}

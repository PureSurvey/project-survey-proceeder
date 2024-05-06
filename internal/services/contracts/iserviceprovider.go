package contracts

import (
	contextcontracts "project-survey-proceeder/internal/context/contracts"
	"project-survey-proceeder/internal/dbcache"
	surveymarkupcontracts "project-survey-proceeder/internal/surveymarkup/contracts"
	"project-survey-proceeder/internal/targeting/contracts"
)

type IServiceProvider interface {
	GetContextFiller() contextcontracts.IRequestFiller
	GetDbRepo() *dbcache.Repo
	GetTargetingService() contracts.ITargetingService
	GetSurveyMarkupService() surveymarkupcontracts.ISurveyMarkupService
}

package services

import (
	"github.com/valyala/fastjson"
	"project-survey-proceeder/internal/configuration"
	contextcontracts "project-survey-proceeder/internal/context/contracts"
	"project-survey-proceeder/internal/context/filler"
	"project-survey-proceeder/internal/dbcache"
	"project-survey-proceeder/internal/dbcache/reader"
	"project-survey-proceeder/internal/geolocation"
	geolocationcontracts "project-survey-proceeder/internal/geolocation/contracts"
	"project-survey-proceeder/internal/pools"
	servicescontracts "project-survey-proceeder/internal/services/contracts"
	"project-survey-proceeder/internal/surveymarkup"
	surveymarkupcontracts "project-survey-proceeder/internal/surveymarkup/contracts"
	"project-survey-proceeder/internal/targeting"
	"project-survey-proceeder/internal/targeting/contracts"
)

type Provider struct {
	parserPool          *fastjson.ParserPool
	userAgentPool       *pools.UserAgentPool
	geolocationService  geolocationcontracts.IGeolocationService
	targetingService    contracts.ITargetingService
	contextFiller       contextcontracts.IRequestFiller
	surveyMarkupService surveymarkupcontracts.ISurveyMarkupService
	dbRepo              *dbcache.Repo
}

func NewProvider(appConfiguration *configuration.AppConfiguration) servicescontracts.IServiceProvider {
	dbReader := reader.NewSqlReader(appConfiguration.SurveyGeneratorAddress)

	provider := &Provider{
		parserPool:          &fastjson.ParserPool{},
		userAgentPool:       &pools.UserAgentPool{},
		geolocationService:  geolocation.NewService(),
		targetingService:    targeting.NewTargetingService(),
		surveyMarkupService: surveymarkup.NewService(appConfiguration.SurveyGeneratorAddress),
		dbRepo:              dbcache.NewRepo(dbReader),
	}

	provider.dbRepo.Reload()

	provider.contextFiller = filler.NewBaseFiller(provider.userAgentPool, provider.geolocationService)

	return provider
}

func (p *Provider) GetDbRepo() *dbcache.Repo {
	return p.dbRepo
}

func (p *Provider) GetTargetingService() contracts.ITargetingService {
	return p.targetingService
}

func (p *Provider) GetContextFiller() contextcontracts.IRequestFiller {
	return p.contextFiller
}

func (p *Provider) GetSurveyMarkupService() surveymarkupcontracts.ISurveyMarkupService {
	return p.surveyMarkupService
}

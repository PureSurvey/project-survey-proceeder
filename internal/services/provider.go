package services

import (
	"github.com/valyala/fastjson"
	"log"
	"project-survey-proceeder/internal/configuration"
	contextcontracts "project-survey-proceeder/internal/context/contracts"
	"project-survey-proceeder/internal/context/filler"
	"project-survey-proceeder/internal/dbcache"
	"project-survey-proceeder/internal/dbcache/reader"
	eventcontracts "project-survey-proceeder/internal/events/contracts"
	"project-survey-proceeder/internal/events/kafka"
	"project-survey-proceeder/internal/geolocation"
	geolocationcontracts "project-survey-proceeder/internal/geolocation/contracts"
	"project-survey-proceeder/internal/pools"
	"project-survey-proceeder/internal/respondents"
	servicescontracts "project-survey-proceeder/internal/services/contracts"
	"project-survey-proceeder/internal/surveymarkup"
	surveymarkupcontracts "project-survey-proceeder/internal/surveymarkup/contracts"
	"project-survey-proceeder/internal/targeting"
	"project-survey-proceeder/internal/targeting/contracts"
	"project-survey-proceeder/internal/trackers"
)

type Provider struct {
	parserPool          *fastjson.ParserPool
	userAgentPool       *pools.UserAgentPool
	geolocationService  geolocationcontracts.IGeolocationService
	targetingService    contracts.ITargetingService
	unitContextFiller   contextcontracts.IRequestFiller
	eventContextFiller  contextcontracts.IRequestFiller
	surveyMarkupService surveymarkupcontracts.ISurveyMarkupService
	dbRepo              *dbcache.Repo
	decryptor           *trackers.Decryptor
	eventProducer       eventcontracts.IEventProducer
	usersService        *respondents.Service
}

func NewProvider(appConfiguration *configuration.AppConfiguration) servicescontracts.IServiceProvider {
	dbReader := reader.NewSqlReader(appConfiguration.DbCacheConfiguration)
	dbRepo := dbcache.NewRepo(appConfiguration.DbCacheConfiguration, dbReader)
	go dbRepo.RunReloadCycle()

	userAgentPool := &pools.UserAgentPool{}
	geolocationService := geolocation.NewService()
	err := geolocationService.Init()
	if err != nil {
		log.Fatalf(err.Error())
	}

	decryptor := trackers.NewDecryptor(appConfiguration)

	eventProducer := kafka.NewProducer(appConfiguration.EventsConfiguration)
	err = eventProducer.Init()
	if err != nil {
		log.Fatalf(err.Error())
	}

	surveyMarkupService := surveymarkup.NewService(appConfiguration.SurveyGeneratorAddress)
	err = surveyMarkupService.Init()
	if err != nil {
		log.Fatalf(err.Error())
	}

	provider := &Provider{
		parserPool:          &fastjson.ParserPool{},
		userAgentPool:       userAgentPool,
		geolocationService:  geolocationService,
		targetingService:    targeting.NewTargetingService(dbRepo),
		surveyMarkupService: surveyMarkupService,
		dbRepo:              dbRepo,
		decryptor:           decryptor,
		eventProducer:       eventProducer,
		unitContextFiller:   filler.NewUnitFiller(userAgentPool, geolocationService),
		eventContextFiller:  filler.NewEventFiller(userAgentPool, geolocationService, decryptor),
		usersService:        respondents.NewService(),
	}

	return provider
}

func (p *Provider) GetDbRepo() *dbcache.Repo {
	return p.dbRepo
}

func (p *Provider) GetTargetingService() contracts.ITargetingService {
	return p.targetingService
}

func (p *Provider) GetUnitContextFiller() contextcontracts.IRequestFiller {
	return p.unitContextFiller
}

func (p *Provider) GetEventContextFiller() contextcontracts.IRequestFiller {
	return p.eventContextFiller
}

func (p *Provider) GetSurveyMarkupService() surveymarkupcontracts.ISurveyMarkupService {
	return p.surveyMarkupService
}

func (p *Provider) GetEventProducer() eventcontracts.IEventProducer {
	return p.eventProducer
}

func (p *Provider) GetRespondentsService() *respondents.Service {
	return p.usersService
}

func (p *Provider) Dispose() error {
	if p.surveyMarkupService != nil {
		err := p.surveyMarkupService.Close()
		if err != nil {
			return err
		}
	}

	if p.eventProducer != nil {
		p.eventProducer.CloseConnection()
	}

	return nil
}

package services

import (
	"github.com/valyala/fastjson"
	contextcontracts "project-survey-proceeder/internal/context/contracts"
	"project-survey-proceeder/internal/context/filler"
	"project-survey-proceeder/internal/geolocation"
	geolocationcontracts "project-survey-proceeder/internal/geolocation/contracts"
	"project-survey-proceeder/internal/pools"
	servicescontracts "project-survey-proceeder/internal/services/contracts"
)

type Provider struct {
	ParserPool         *fastjson.ParserPool
	UserAgentPool      *pools.UserAgentPool
	GeolocationService geolocationcontracts.IGeolocationService
}

func NewProvider() servicescontracts.IServiceProvider {
	return &Provider{
		ParserPool:         &fastjson.ParserPool{},
		UserAgentPool:      &pools.UserAgentPool{},
		GeolocationService: geolocation.NewService(),
	}
}

func (p *Provider) GetContextFiller() contextcontracts.IRequestFiller {
	return filler.NewBaseFiller(p.UserAgentPool, p.GeolocationService)
}

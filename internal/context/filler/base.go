package filler

import (
	"fmt"
	"github.com/ferluci/fast-realip"
	"github.com/valyala/fasthttp"
	"hash/fnv"
	"project-survey-proceeder/internal/context"
	geolocationcontracts "project-survey-proceeder/internal/geolocation/contracts"
	"project-survey-proceeder/internal/pools"
	"strconv"
	"time"
)

type Base struct {
	UserAgentPool      *pools.UserAgentPool
	GeolocationService geolocationcontracts.IGeolocationService
}

func NewBaseFiller(pool *pools.UserAgentPool, geoservice geolocationcontracts.IGeolocationService) *Base {
	return &Base{UserAgentPool: pool, GeolocationService: geoservice}
}

func (f *Base) FillFromRequest(prCtx *context.ProceederContext, ctx *fasthttp.RequestCtx) error {
	if ctx.QueryArgs() == nil {
		return fmt.Errorf("nil query args")
	}

	prCtx.RequestTimestamp = time.Now().UTC()

	f.fillUserAgent(prCtx, ctx)
	f.fillIp(prCtx, ctx)
	f.fillCountry(prCtx)
	f.fillUaIpHash(prCtx)

	return nil
}

func (f *Base) fillUserAgent(prCtx *context.ProceederContext, ctx *fasthttp.RequestCtx) {
	prCtx.UserAgent = string(ctx.Request.Header.UserAgent())

	if prCtx.UserAgent != `` {
		ua := f.UserAgentPool.Get()
		defer f.UserAgentPool.Put(ua)

		ua.Parse(prCtx.UserAgent)
		prCtx.IsMobile = ua.Mobile()
		prCtx.Platform = ua.Platform()
	}
}

func (f *Base) fillIp(prCtx *context.ProceederContext, ctx *fasthttp.RequestCtx) {
	prCtx.Ip = realip.FromRequest(ctx)
}

func (f *Base) fillCountry(prCtx *context.ProceederContext) {
	country, err := f.GeolocationService.GetCountryByIp(prCtx.Ip)
	if err != nil {
		return
	}
	prCtx.Country = country

	lang, err := f.GeolocationService.GetLanguageByCountry(prCtx.Country)
	if err != nil {
		return
	}
	prCtx.Language = lang
}

func (f *Base) fillUaIpHash(prCtx *context.ProceederContext) {
	stringToHash := prCtx.UserAgent + prCtx.Ip

	h := fnv.New32a()
	h.Write([]byte(stringToHash))
	prCtx.UaIpHash = strconv.Itoa(int(h.Sum32()))
}

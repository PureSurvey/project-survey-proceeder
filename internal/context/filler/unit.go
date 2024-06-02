package filler

import (
	"github.com/valyala/fasthttp"
	"project-survey-proceeder/internal/context"
	"project-survey-proceeder/internal/context/contracts"
	geolocationcontracts "project-survey-proceeder/internal/geolocation/contracts"
	"project-survey-proceeder/internal/pools"
)

type Unit struct {
	*Base
}

func NewUnitFiller(pool *pools.UserAgentPool, geoservice geolocationcontracts.IGeolocationService) contracts.IRequestFiller {
	return &Unit{Base: NewBaseFiller(pool, geoservice)}
}

func (f *Unit) FillFromRequest(prCtx *context.ProceederContext, ctx *fasthttp.RequestCtx) error {
	err := f.Base.FillFromRequest(prCtx, ctx)
	if err != nil {
		return err
	}

	f.fillEntityIds(prCtx, ctx)

	return nil
}

func (f *Unit) fillEntityIds(prCtx *context.ProceederContext, ctx *fasthttp.RequestCtx) {
	unitId, err := ctx.QueryArgs().GetUint(`unitId`)
	if err != nil {
		return
	}

	prCtx.UnitId = unitId
}

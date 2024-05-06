package contracts

import (
	"project-survey-proceeder/internal/context"
	dbcache "project-survey-proceeder/internal/dbcache/objects"
)

type ITargetingService interface {
	IsMatched(survey *dbcache.Survey, prCtx *context.ProceederContext) bool
}

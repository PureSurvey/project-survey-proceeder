package contract

import dbcache "project-survey-proceeder/internal/dbcache/objects"

type IService interface {
	IsMatched(survey *dbcache.Survey) bool
}

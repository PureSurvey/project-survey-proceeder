package targeting

import (
	"project-survey-proceeder/internal/context"
	"project-survey-proceeder/internal/dbcache"
	"project-survey-proceeder/internal/dbcache/objects"
)

type Service struct {
	prCtx  *context.ProceederContext
	dbRepo *dbcache.Repo
}

func NewTargetingService(prCtx *context.ProceederContext) *Service {
	return &Service{prCtx: prCtx}
}

func (s *Service) IsMatched(survey *objects.Survey) bool {
	return s.isMatchedByCountry(survey)
}

func (s *Service) isMatchedByCountry(survey *objects.Survey) bool {
	countries := s.dbRepo.GetCountriesByTargetingId(survey.TargetingId)
	for _, country := range countries {
		if s.prCtx.Country == country {
			return true
		}
	}

	return false
}

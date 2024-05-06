package targeting

import (
	"project-survey-proceeder/internal/context"
	"project-survey-proceeder/internal/dbcache"
	"project-survey-proceeder/internal/dbcache/objects"
	"project-survey-proceeder/internal/targeting/contracts"
)

type Service struct {
	dbRepo *dbcache.Repo
}

func NewTargetingService() contracts.ITargetingService {
	return &Service{}
}

func (s *Service) IsMatched(survey *objects.Survey, prCtx *context.ProceederContext) bool {
	return s.isMatchedByCountry(survey, prCtx.Country)
}

func (s *Service) isMatchedByCountry(survey *objects.Survey, country string) bool {
	countries := s.dbRepo.GetCountriesByTargetingId(survey.TargetingId)
	for _, ctr := range countries {
		if country == ctr {
			return true
		}
	}

	return false
}

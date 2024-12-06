package respondents

import (
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"project-survey-proceeder/internal/context"
	"sync"
	"time"
)

type Respondent struct {
	Id                     string
	UaIpHash               string
	CompletedSurveysByUnit map[int]map[int]int
}

type Service struct {
	idCache   *cache.Cache
	hashCache *cache.Cache

	incrementLock *sync.RWMutex
}

func NewService() *Service {
	return &Service{
		idCache:       cache.New(3*24*time.Hour, 2*time.Hour),
		hashCache:     cache.New(1*24*time.Hour, 2*time.Hour),
		incrementLock: &sync.RWMutex{},
	}
}

func (c *Service) UpdateContext(context *context.ProceederContext) {
	respondent := c.getRespondent(context.UserCookie, context.UaIpHash)
	if respondent == nil {
		respondent = &Respondent{
			Id:                     uuid.New().String(),
			UaIpHash:               context.UaIpHash,
			CompletedSurveysByUnit: make(map[int]map[int]int),
		}
	}

	c.idCache.Set(respondent.Id, respondent, cache.DefaultExpiration)
	c.hashCache.Set(respondent.UaIpHash, respondent, cache.DefaultExpiration)

	context.UserCookie = respondent.Id
}

func (c *Service) GetCompletedSurveysByUnit(context *context.ProceederContext) int {
	respondent := c.getRespondent(context.UserCookie, context.UaIpHash)
	if respondent == nil {
		return 0
	}
	if respondent.CompletedSurveysByUnit == nil {
		return 0
	}
	surveysByUnit := respondent.CompletedSurveysByUnit[context.UnitId]
	if surveysByUnit == nil {
		return 0
	}
	return len(surveysByUnit)
}

func (c *Service) GetSurveyTakesByUnit(context *context.ProceederContext, surveyId int) int {
	respondent := c.getRespondent(context.UserCookie, context.UaIpHash)
	if respondent == nil {
		return 0
	}
	if respondent.CompletedSurveysByUnit == nil {
		return 0
	}
	surveysByUnit := respondent.CompletedSurveysByUnit[context.UnitId]
	if surveysByUnit == nil {
		return 0
	}

	return surveysByUnit[surveyId]
}

func (c *Service) IncrementSurveyTakeByUnit(context *context.ProceederContext) {
	respondent := c.getRespondent(context.UserCookie, context.UaIpHash)
	if respondent == nil {
		return
	}
	c.incrementLock.Lock()
	defer c.incrementLock.Unlock()

	if respondent.CompletedSurveysByUnit == nil {
		respondent.CompletedSurveysByUnit = make(map[int]map[int]int)
	}
	surveysByUnit := respondent.CompletedSurveysByUnit[context.UnitId]
	if surveysByUnit == nil {
		respondent.CompletedSurveysByUnit[context.UnitId] = make(map[int]int)
		surveysByUnit = respondent.CompletedSurveysByUnit[context.UnitId]
	}

	surveysByUnit[context.SurveyId] += 1
}

func (c *Service) getRespondent(respondentId string, uaIpHash string) *Respondent {
	respondentInterface, found := c.idCache.Get(respondentId)
	if !found {
		respondentInterface, found = c.hashCache.Get(uaIpHash)
		if !found {
			return nil
		}
	}
	return respondentInterface.(*Respondent)
}

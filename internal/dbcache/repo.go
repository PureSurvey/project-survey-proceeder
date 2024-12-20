package dbcache

import (
	"project-survey-proceeder/internal/configuration"
	"project-survey-proceeder/internal/dbcache/contracts"
	"project-survey-proceeder/internal/dbcache/objects"
	"time"
)

const (
	storedProcedureName = "objects_01"
)

type Repo struct {
	config *configuration.DbCacheConfiguration
	reader contracts.IReader
	cache  *Cache
}

func NewRepo(config *configuration.DbCacheConfiguration, reader contracts.IReader) *Repo {
	return &Repo{reader: reader, config: config}
}

func (r *Repo) RunReloadCycle() {
	for {
		r.reload()
		time.Sleep(time.Duration(r.config.ReloadSleepTime) * time.Second)
	}
}

func (r *Repo) reload() {
	err := r.reader.Connect()
	defer r.reader.CloseConnection()
	if err != nil {
		return
	}

	newCache := &Cache{
		Users: map[int]*objects.User{},

		Units:       map[int]*objects.Unit{},
		Surveys:     map[int]*objects.Survey{},
		Templates:   map[int]*objects.Template{},
		Appearances: map[int]*objects.Appearance{},
		Questions:   map[int]*objects.Question{},
		Options:     map[int]*objects.Option{},

		Targetings:         map[int]*objects.Targeting{},
		CountryInTargeting: map[int][]string{},

		SurveysByUnit: map[int][]*objects.Survey{},

		TranslationsByQuestionLine: map[int]map[string]*objects.Translation{},
		TranslationsByOption:       map[int]map[string]*objects.Translation{},
	}

	res, err := r.reader.GetStoredProcedureResult(storedProcedureName)
	i := 0
	for cont := true; cont; cont = res.NextResultSet() {
		switch i {
		case 0:
			err := newCache.fillUsers(res)
			if err != nil {
				return
			}
			break
		case 1:
			err := newCache.fillUnits(res)
			if err != nil {
				return
			}
			break
		case 2:
			err := newCache.fillSurveys(res)
			if err != nil {
				return
			}
			break
		case 3:
			err := newCache.fillSurveyInUnits(res)
			if err != nil {
				return
			}
			break
		case 4:
			err := newCache.fillTemplates(res)
			if err != nil {
				return
			}
			break
		case 5:
			err := newCache.fillTargetings(res)
			if err != nil {
				return
			}
			break
		case 6:
			err := newCache.fillCountryInTargetings(res)
			if err != nil {
				return
			}
			break
		case 7:
			err := newCache.fillAppearances(res)
			if err != nil {
				return
			}
			break
		case 8:
			err := newCache.fillQuestions(res)
			if err != nil {
				return
			}
			break
		case 9:
			err := newCache.fillOptions(res)
			if err != nil {
				return
			}
			break
		case 10:
			err := newCache.fillQuestionTranslations(res)
			if err != nil {
				return
			}
			break
		case 11:
			err := newCache.fillOptionTranslations(res)
			if err != nil {
				return
			}
			break
		}
		i++
	}

	r.cache = newCache
}

func (r *Repo) GetUnitById(id int) *objects.Unit {
	return r.cache.Units[id]
}

func (r *Repo) GetUserById(id int) *objects.User {
	return r.cache.Users[id]
}

func (r *Repo) GetSurveysByUnitId(id int) []*objects.Survey {
	return r.cache.SurveysByUnit[id]
}

func (r *Repo) GetCountriesByTargetingId(id int) []string {
	return r.cache.CountryInTargeting[id]
}

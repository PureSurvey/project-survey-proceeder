package dbcache

import (
	"database/sql"
	"project-survey-proceeder/internal/dbcache/contracts"
	"project-survey-proceeder/internal/dbcache/objects"
)

const (
	storedProcedureName = "objects_01"
)

type Cache struct {
	Users map[int]*objects.User

	Units       map[int]*objects.Unit
	Surveys     map[int]*objects.Survey
	Templates   map[int]*objects.Template
	Appearances map[int]*objects.Appearance
	Questions   map[int]*objects.Question
	Options     map[int]*objects.Option

	Targetings         map[int]*objects.Targeting
	CountryInTargeting map[int]*string

	SurveysByUnit map[int][]*objects.Survey

	TranslationsByQuestionLine map[int]map[string]*objects.Translation
	TranslationsByOption       map[int]map[string]*objects.Translation
}

type Repo struct {
	reader contracts.IReader
	cache  *Cache
}

func NewRepo(reader contracts.IReader) *Repo {
	return &Repo{reader: reader}
}

func (r *Repo) Reload() {
	err := r.reader.Connect()
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
		CountryInTargeting: map[int]*string{},

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
		}
		i++
	}

	r.cache = newCache
}

func (r *Repo) GetUnitById(id int) *objects.Unit {
	return r.cache.Units[id]
}

func (r *Repo) GetSurveysByUnitId(id int) []*objects.Survey {
	return r.cache.SurveysByUnit[id]
}

func (c *Cache) fillUsers(rows *sql.Rows) error {
	for rows.Next() {
		var id int
		var role string
		err := rows.Scan(&id, &role)
		if err != nil {
			return err
		}

		user := objects.NewUser(id, role)
		c.Users[id] = user
	}

	return nil
}

func (c *Cache) fillUnits(rows *sql.Rows) error {
	for rows.Next() {
		var id, userId, appearanceId, maxPerDevice int
		var oneTakePerDevice, hideAfterNoSurveys bool
		var name, message string

		err := rows.Scan(&id, &name, &userId, &appearanceId, &oneTakePerDevice,
			&maxPerDevice, &hideAfterNoSurveys, &message)
		if err != nil {
			return err
		}

		unit := objects.NewUnit(id, name, userId, appearanceId, oneTakePerDevice,
			maxPerDevice, hideAfterNoSurveys, message)
		c.Units[id] = unit
	}

	return nil
}

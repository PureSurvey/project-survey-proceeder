package geolocation

import (
	"github.com/ip2location/ip2location-go/v9"
	"path/filepath"
	"project-survey-proceeder/internal/geolocation/contracts"
	"strings"
)

type Service struct {
	db            *ip2location.DB
	countryInfoDb *ip2location.CI
}

func NewService() contracts.IGeolocationService {
	return &Service{}
}

func (s *Service) Init() error {
	path, _ := filepath.Abs("thirdparty/ip2location/ip2location.BIN")
	db, err := ip2location.OpenDB(path)
	if err != nil {
		return err
	}
	s.db = db

	countryInfoDb, err := ip2location.OpenCountryInfo("thirdparty/ip2location/ip2locationCountryInfo.CSV")
	if err != nil {
		return err
	}
	s.countryInfoDb = countryInfoDb

	return nil
}

func (s *Service) GetCountryByIp(ip string) (string, error) {
	record, err := s.db.Get_country_short(ip)
	if err != nil {
		return "", err
	}

	return record.Country_short, nil
}

func (s *Service) GetLanguageByCountry(country string) (string, error) {
	record, err := s.countryInfoDb.GetCountryInfo(country)
	if err != nil {
		return "", err
	}

	return strings.ToLower(record[0].Lang_code), nil
}

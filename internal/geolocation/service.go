package geolocation

import (
	"github.com/ip2location/ip2location-go/v9"
	"project-survey-proceeder/internal/geolocation/contracts"
)

type Service struct {
	db *ip2location.DB
}

func NewService() contracts.IGeolocationService {
	return &Service{}
}

func (s *Service) Init() error {
	db, err := ip2location.OpenDB("..\\..\\ip2location\\ip2location.bin")
	if err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Service) GetCountryByIp(ip string) (string, error) {
	record, err := s.db.Get_country_short(ip)
	if err != nil {
		return "", err
	}

	return record.Country_short, nil
}

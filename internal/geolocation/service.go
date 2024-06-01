package geolocation

import (
	"github.com/ip2location/ip2location-go/v9"
	"path/filepath"
	"project-survey-proceeder/internal/geolocation/contracts"
)

type Service struct {
	db *ip2location.DB
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
	return nil
}

func (s *Service) GetCountryByIp(ip string) (string, error) {
	record, err := s.db.Get_country_short(ip)
	if err != nil {
		return "", err
	}

	return record.Country_short, nil
}

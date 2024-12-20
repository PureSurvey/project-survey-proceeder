package contracts

type IGeolocationService interface {
	Init() error
	GetCountryByIp(ip string) (string, error)
	GetLanguageByCountry(country string) (string, error)
}

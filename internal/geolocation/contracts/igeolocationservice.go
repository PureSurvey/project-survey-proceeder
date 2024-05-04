package contracts

type IGeolocationService interface {
	GetCountryByIp(ip string) (string, error)
}

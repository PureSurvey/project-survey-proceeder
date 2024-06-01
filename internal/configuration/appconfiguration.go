package configuration

type AppConfiguration struct {
	SurveyGeneratorAddress string                `json:"surveyGeneratorAddress"`
	DbCacheConfiguration   *DbCacheConfiguration `json:"dbCacheConfiguration"`
}

type DbCacheConfiguration struct {
	ConnectionRetryCount     int    `json:"connectionRetryCount"`
	ConnectionRetrySleepTime int    `json:"connectionRetryTimeout"`
	ConnectionString         string `json:"connectionString"`
	StoredProcedure          string `json:"storedProcedure"`
	ReloadSleepTime          int    `json:"reloadSleepTime"`
}

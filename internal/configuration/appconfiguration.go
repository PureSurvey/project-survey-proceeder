package configuration

type AppConfiguration struct {
	Host                   string                `json:"host"`
	SurveyGeneratorAddress string                `json:"surveyGeneratorAddress"`
	EncryptionSecret       string                `json:"encryptionSecret"` // 32 bytes
	DbCacheConfiguration   *DbCacheConfiguration `json:"dbCacheConfiguration"`
	EventsConfiguration    *EventsConfiguration  `json:"eventsConfiguration"`
}

type DbCacheConfiguration struct {
	ConnectionRetryCount     int    `json:"connectionRetryCount"`
	ConnectionRetrySleepTime int    `json:"connectionRetrySleepTime"`
	ConnectionString         string `json:"connectionString"`
	StoredProcedure          string `json:"storedProcedure"`
	ReloadSleepTime          int    `json:"reloadSleepTime"`
}

type EventsConfiguration struct {
	BrokerUrl string `json:"brokerUrl"`
	Topic     string `json:"topic"`
}

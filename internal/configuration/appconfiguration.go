package configuration

type AppConfiguration struct {
	DbConnectionString     string `json:"dbConnectionString"`
	SurveyGeneratorAddress string `json:"surveyGeneratorAddress"`
}

package configuration

import (
	"encoding/json"
	"fmt"
	"os"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(file string) (*AppConfiguration, error) {
	fileContent, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error when reading appsettings: %v", err)
	}

	var cfg AppConfiguration
	err = json.Unmarshal(fileContent, &cfg)
	if err != nil {
		return nil, fmt.Errorf("error when parsing appsettings: %v", err)
	}

	return &cfg, nil
}

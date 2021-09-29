package config

import (
	"testing"

	"github.com/spf13/viper"
)

func initConfiguration(configName string, t *testing.T) *Configuration {
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")
	viper.SetDefault("Keyword", "")
	var configuration Configuration

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	viper.Unmarshal(&configuration)
	t.Logf("Configuration is:\n%s", configuration)

	return &configuration
}

func TestGeneralConfigValues(t *testing.T) {
	configuration := initConfiguration("keyword_config.yaml", t)

	if configuration.Images != 40 {
		t.Errorf("Configuration Images did not equal to 40, got %d", configuration.Images)
	}

	if configuration.DownloadPath != "../downloads" {
		t.Errorf("Configuration DownloadPath did not equal to ../downloads, got %s", configuration.DownloadPath)
	}

	if configuration.NumRoutines != 4 {
		t.Errorf("Configuration NumRoutines did not equal to 4, got %d", configuration.NumRoutines)
	}

	if configuration.Site != "bing" {
		t.Errorf("Configuration Site did not equal to bing, got %s", configuration.Site)
	}
}

func TestUnpackingKeywordConfig(t *testing.T) {
	configuration := initConfiguration("keyword_config.yaml", t)

	if len(configuration.Species) != 0 {
		t.Errorf("Expected Species list to be of length 0, but instead has a length of %d", len(configuration.Species))
	}

	if configuration.Keyword == "" {
		t.Errorf("Configuration Keyword was an empty string, but it should be tree!")
	}
}

func TestUnpackingSpeciesConfig(t *testing.T) {
	configuration := initConfiguration("confg.yaml", t)

	if len(configuration.Species) != 26 {
		t.Errorf("Configuation Species list does not equal to length 26, it has lenght of %d", len(configuration.Species))
	}

	if configuration.Keyword != "" {
		t.Errorf("Configuration Keyword is not empty, got %s", configuration.Keyword)
	}

}

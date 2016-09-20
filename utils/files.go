package utils

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const exampleconfig = `# This is an example config
HeadOfStation = "station.manager"
AssistantHeadOfStation = "assistant.station.manager"
ApiKey = "apikeygoeshere"
StandDownPeriod = 28 #days`

type Configurer interface {
	IsHistoricalOfficerValid(now, to time.Time) (bool, error)
	GetHeadOfStation() string
	GetAssistantHeadOfStation() string
	GetApiKey() string
}

type configData struct {
	HeadOfStation          string
	AssistantHeadOfStation string
	ApiKey                 string
	StandDownPeriod        int
}

type Config struct {
	Configurer
	configData
}

func (c Config) GetHeadOfStation() string {
	return c.configData.HeadOfStation
}

func (c Config) GetAssistantHeadOfStation() string {
	return c.configData.AssistantHeadOfStation
}

func (c Config) GetApiKey() string {
	return c.configData.ApiKey
}

func (c Config) IsHistoricalOfficerValid(now, to time.Time) (bool, error) {
	var delta, err = time.ParseDuration(fmt.Sprintf("%dh", c.configData.StandDownPeriod*24))
	if err != nil {
		return false, err
	}
	return (!to.IsZero() && now.Before(to.Add(delta))), nil
}

// NewConfigFromFile loads a config file and returns it.
func NewConfigFromFile(path string) (c Config, err error) {
	absPath, _ := filepath.Abs(path)
	b, err := ioutil.ReadFile(absPath)
	if err != nil {
		return
	}
	s := string(b)
	var cd configData
	_, err = toml.Decode(s, &cd)
	c = Config{configData: cd}
	return
}

func WriteExampleConfigToFile(path string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	_, err = f.WriteString(exampleconfig)
	return
}

func WriteAliasesToFile(aliases, file string) (err error) {
	f, err := os.Create(file)
	if err != nil {
		return
	}
	defer f.Close()
	t := time.Now()
	_, err = f.WriteString(fmt.Sprintf("# Generated: %s\n%s", t.String(), aliases))
	return
}

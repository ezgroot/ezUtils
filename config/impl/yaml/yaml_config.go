package yaml

import (
	"github.com/ezgroot/ezUtils/config/impl/data"
	"github.com/ghodss/yaml"
	jsoniter "github.com/json-iterator/go"
)

type yamlConfig struct{}

// TranslateBytes Converts byte stream data into corresponding configuration structure objects.
func (y *yamlConfig) TranslateBytes(bytesInfo []byte, config interface{}) error {
	jsonBytes, err := yaml.YAMLToJSON(bytesInfo)
	if err != nil {
		return err
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(jsonBytes, config)
	if err != nil {
		return err
	}

	return nil
}

// GetYaml get yaml config.
func GetYaml() data.Translator {
	return &yamlConfig{}
}

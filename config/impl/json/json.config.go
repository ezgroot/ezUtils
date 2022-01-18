package json

import (
	"github.com/ezgroot/ezUtils/config/impl/data"

	jsoniter "github.com/json-iterator/go"
)

type jsonConfig struct{}

// TranslateBytes Converts byte stream data into corresponding configuration structure objects.
func (j *jsonConfig) TranslateBytes(bytesInfo []byte, config interface{}) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(bytesInfo, config)
	if err != nil {
		return err
	}

	return nil
}

// GetJSON get json config.
func GetJSON() data.Translator {
	return &jsonConfig{}
}

package toml

import (
	"github.com/BurntSushi/toml"
	"github.com/ezgroot/ezUtils/config/impl/data"
)

type tomlConfig struct{}

// TranslateBytes Converts byte stream data into corresponding configuration structure objects.
func (j *tomlConfig) TranslateBytes(bytesInfo []byte, config interface{}) error {
	err := toml.Unmarshal(bytesInfo, config)
	if err != nil {
		return err
	}

	return nil
}

// GetJSON get json config.
func GetToml() data.Translator {
	return &tomlConfig{}
}

package ini

import (
	"github.com/ezgroot/ezUtils/config/impl/data"
	"gopkg.in/gcfg.v1"
)

type iniConfig struct{}

// TranslateBytes Converts byte stream data into corresponding configuration structure objects.
func (i *iniConfig) TranslateBytes(bytesInfo []byte, config interface{}) error {
	err := gcfg.ReadStringInto(config, string(bytesInfo))
	if err != nil {
		return err
	}

	return nil
}

// GetIni get ini config.
func GetIni() data.Translator {
	return &iniConfig{}
}

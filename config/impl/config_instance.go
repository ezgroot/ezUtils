package impl

import (
	"fmt"
	"io/ioutil"
	"path"
	"sync"

	"github.com/ezgroot/ezUtils/config/impl/data"
	"github.com/ezgroot/ezUtils/config/impl/ini"
	"github.com/ezgroot/ezUtils/config/impl/json"
	"github.com/ezgroot/ezUtils/config/impl/yaml"
)

// config file type.
const (
	fileTypeIni  = ".ini"
	fileTypeJSON = ".json"
	fileTypeYaml = ".yaml"
)

// Config is core of config.
type Config struct {
	bytesInfo  []byte
	filePath   string
	fileType   string
	translator data.Translator
}

// ConfigInit read config file.
func (c *Config) ConfigInit(filePath string) error {
	var err error
	c.filePath = filePath

	c.bytesInfo, err = ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var fileSuffix = path.Ext(filePath)
	if fileSuffix == fileTypeIni {
		c.fileType = fileTypeIni
		c.translator = ini.GetIni()
	} else if fileSuffix == fileTypeJSON {
		c.fileType = fileTypeJSON
		c.translator = json.GetJSON()
	} else if fileSuffix == fileTypeYaml {
		c.fileType = fileTypeYaml
		c.translator = yaml.GetYaml()
	} else {
		return fmt.Errorf("config file type = %s not support", fileSuffix)
	}

	return nil
}

// GetAllConfig get all config option.
func (c *Config) GetAllConfig(configType interface{}) error {
	err := c.translator.TranslateBytes(c.bytesInfo, configType)
	if err != nil {
		return err
	}

	return nil
}

func newConfig() *Config {
	return &Config{}
}

var instance *Config
var onceGet sync.Once

// GetConfigInstance get singleton of config.
func GetConfigInstance() *Config {
	onceGet.Do(func() {
		instance = newConfig()
	})

	return instance
}

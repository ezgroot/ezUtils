package config

import "github.com/ezgroot/ezUtils/config/impl"

// Init init to give config file.
func Init(filePath string) error {
	return impl.GetConfigInstance().ConfigInit(filePath)
}

// GetAll get all configuration option.
func GetAll(allConfigType interface{}) error {
	return impl.GetConfigInstance().GetAllConfig(allConfigType)
}

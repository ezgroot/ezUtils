package data

// Translator Converts byte stream data into corresponding configuration structure objects.
type Translator interface {
	TranslateBytes(bytesInfo []byte, config interface{}) error
}

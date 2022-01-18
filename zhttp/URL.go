package zhttp

import (
	"net/url"
)

// AddParameterToURL add parameter to URL
func AddParameterToURL(URL string, parameter map[string]string) (string, error) {
	if len(parameter) == 0 {
		return URL, nil
	}

	pu, err := url.Parse(URL)
	if err != nil {
		return "", err
	}

	kvList := pu.Query()
	for k, v := range parameter {
		kvList.Set(k, v)
	}

	pu.RawQuery = kvList.Encode()

	return pu.String(), nil
}

// GetParameterFromURL get parameter from URL
func GetParameterFromURL(URL *url.URL) (url.Values, error) {
	return url.ParseQuery(URL.RawQuery)
}

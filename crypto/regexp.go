package crypto

import (
	"regexp"
)

// IsIPv4 is IPv4 address, 0.0.0.0 - 255.255.255.255.
func IsIPv4(str string) (bool, error) {
	patten := "^([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\\." +
		"([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\\." +
		"([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\\." +
		"([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])$"
	ret, err := regexp.MatchString(patten, str)
	if err != nil {
		return false, err
	}

	return ret, nil
}

// IsEmail is email address.
func IsEmail(str string) (bool, error) {
	patten := "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	ret, err := regexp.MatchString(patten, str)
	if err != nil {
		return false, err
	}

	return ret, nil
}

// IsPhone is phone number.
func IsPhone(str string) (bool, error) {
	patten := "^1[0-9]{10}$"
	ret, err := regexp.MatchString(patten, str)
	if err != nil {
		return false, err
	}

	return ret, nil
}

// IsAccountEn is english account, windows not support file or dir name include: \/:*?"<>|
func IsAccountEn(str string) (bool, error) {
	patten := "^[a-zA-Z0-9]{1}([a-zA-Z0-9_-]||[.@#$%&]){5,127}$"
	ret, err := regexp.MatchString(patten, str)
	if err != nil {
		return false, err
	}

	return ret, nil
}

// IsAccountZh is chinese account.
func IsAccountZh(str string) (bool, error) {
	patten := "^([a-zA-Z0-9_\u4e00-\u9fa5]|[-.@#$%&]){2,127}$"
	ret, err := regexp.MatchString(patten, str)
	if err != nil {
		return false, err
	}

	return ret, nil
}

// IsIPPort is net port, 1-65535.
func IsIPPort(str string) (bool, error) {
	patten := "^(([1-9]{1}[0-9]{0,3})|([1-5]{1}[0-9]{0,4})|(6[0-4]{1}[0-9]{3})|" +
		"(65[0-4]{1}[0-9]{2})|(655[0-2]{1}[0-9]{1})|(6553[0-5]{1}))$"
	ret, err := regexp.MatchString(patten, str)
	if err != nil {
		return false, err
	}

	return ret, nil
}

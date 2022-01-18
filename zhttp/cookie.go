package zhttp

import "net/http"

// ClientSetCookie client set cookie.
func ClientSetCookie(request *http.Request, data map[string]string) {
	for name, value := range data {
		cookie := &http.Cookie{Name: name, Value: value, Path: "/", MaxAge: 86400}
		request.AddCookie(cookie)
	}
}

// ServerSetCookie server set cookie.
func ServerSetCookie(rw http.ResponseWriter, data map[string]string) {
	for name, value := range data {
		cookie := &http.Cookie{Name: name, Value: value, Path: "/", MaxAge: 86400}
		http.SetCookie(rw, cookie)
	}
}

// ServerDelCookie server revoke cookie.
func ServerDelCookie(rw http.ResponseWriter, name string) {
	cookie := &http.Cookie{Name: name, Path: "/", MaxAge: -1}
	http.SetCookie(rw, cookie)
}

// ServerGetCookie server get cookie value from request by key
func ServerGetCookie(request *http.Request, key string) (string, error) {
	cookie, err := request.Cookie(key)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

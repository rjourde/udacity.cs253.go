package tools

import (
	"net/http"
	"github.com/gorilla/securecookie"
)

func StoreCookie(w http.ResponseWriter, r *http.Request, storedCookie *securecookie.SecureCookie, cookieName, cookieValue string) {
	value := map[string]string{
		cookieName : cookieValue,
	}
	if encoded, err := storedCookie.Encode(cookieName, value); err == nil {
		cookie := &http.Cookie{
			Name:  cookieName,
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

func FetchCookie(r *http.Request, storedCookie *securecookie.SecureCookie, cookieName string) string {
	if cookie, err := r.Cookie(cookieName); err == nil {
		value := make(map[string]string)
		if(cookie != nil) {
			err = storedCookie.Decode(cookieName, cookie.Value, &value)
			if (len(value[cookieName]) > 0 && err == nil) {
				return value[cookieName]
			}
		}
	}

	return ""
}

func ClearCookie(w http.ResponseWriter, cookieName string) {
	cookie := &http.Cookie{
		Name:  cookieName,
		Value: "",
		Path:  "/",
	}
	http.SetCookie(w, cookie)
}
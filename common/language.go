package common

import (
	"net/http"
	"strings"
)

// GetLang Get Language
func GetLang(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("lang")
	lang := "en"
	if err != nil {
		lang = r.Header.Get("Accept-Language")[:2]
	} else {
		lang = cookie.Value
	}

	langArr := strings.Split(LangStr, ",")
	supportLang := false
	for i := 0; i < len(langArr); i++ {
		if langArr[i] == lang {
			supportLang = true
		}
	}

	if !supportLang {
		lang = "en"
	}
	return lang
}

// SetLang Set Language
func SetLang(w http.ResponseWriter, r *http.Request, lang string) {

	cookie := &http.Cookie{
		Name:     "lang",
		Value:    lang,
		MaxAge:   86400 * 365,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}

// func IsLetter(s string) bool {
// 	for _, r := range s {
// 		if !unicode.IsLetter(r) {
// 			return false
// 		}
// 	}
// 	return true
// }

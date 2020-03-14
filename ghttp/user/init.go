package user

import (
	"encoding/gob"

	"github.com/gorilla/sessions"
)

const (
	cookieStoreAuthKey    = "abc"
	cookieStoreEncryptKey = "1234567891234567"
)

var sessionStore *sessions.CookieStore

type Person struct {
	Name string
}

func init() {
	gob.Register(&Person{})
	sessionStore = sessions.NewCookieStore(
		[]byte(cookieStoreAuthKey),
		[]byte(cookieStoreEncryptKey),
	)

	sessionStore.Options = &sessions.Options{
		HttpOnly: true,
		MaxAge:   60 * 15,
	}
}

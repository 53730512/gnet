package user

import (
	"fmt"
	"net/http"
)

func Secret(w http.ResponseWriter, r *http.Request) {
	session, _ := sessionStore.Get(r, sessionCookieName)

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	println(session)

	fmt.Fprintln(w, "这里还是空空如也!")
}

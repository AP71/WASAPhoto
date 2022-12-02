package auth

import (
	"net/http"
	"strings"
	"wasa-photo/service/api/structures"
	"wasa-photo/service/database"
)

func CheckAuth(db database.AppDatabase, r *http.Request) (bool, structures.User) {

	var user structures.User

	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		return false, structures.User{}
	}

	splitToken := strings.Split(reqToken, "Bearer ")
	user.Id.Value = splitToken[1]

	if db.VerifyToken(&user) {
		return true, user
	}
	return false, structures.User{}
}

package structures

import (
	"regexp"
)

type Username struct {
	Value string `json:"username"`
}

type NewUsername struct {
	Value string `json:"newUsername"`
}

type Identifier struct {
	Value string `json:"identifier"`
}

type User struct {
	Id       Identifier `json:"identifier"`
	Username Username   `json:"username"`
}

type UserData struct {
	Id       string `json:"identifier"`
	Username string `json:"username"`
}

type Users struct {
	List            [10]UserData `json:"users"`
	NextUsersPageId int64        `json:"nextUsersPageId"`
}

func CheckUsername(s string) bool {
	res, _ := regexp.MatchString("^[a-zA-Z0-9_.]*$", s)
	return ((len(s) >= 3) && (len(s) <= 16) && res)
}

func (id Identifier) CheckID() bool {
	return len(id.Value) == 32
}

func (u *User) IsValid() bool {
	var id, username bool
	if u.Id.Value != "" {
		id = u.Id.CheckID()
	} else {
		id = true
	}

	if u.Username.Value != "" {
		username = CheckUsername(u.Username.Value)
	} else {
		username = true
	}

	return id && username
}

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

type UserPage struct {
	Id               string    `json:"id"`
	Username         string    `json:"username"`
	Photos           []PhotoID `json:"photos"`
	NextPhotosPageId int64     `json:"nextPhotosPageId"`
	Followers        int64     `json:"followers"`
	Following        int64     `json:"following"`
	PhotoCounter     int64     `json:"photoCounter"`
}

type UserData struct {
	Id       string `json:"identifier"`
	Username string `json:"username"`
}

type Users struct {
	List            []UserData `json:"users"`
	NextUsersPageId int64      `json:"nextUsersPageId"`
}

func CheckUsername(s string) bool {
	l := len(s)
	res, _ := regexp.MatchString("^[a-zA-Z0-9_.]*$", s)
	return ((l >= 3) && (l <= 16) && res)
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

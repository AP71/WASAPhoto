package database

import (
	"strconv"
	"wasa-photo/service/api/structures"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) GetUserId(username string) (string, error) {
	var id uuid.UUID
	err := db.c.QueryRow(`SELECT Id FROM Users WHERE Username="` + username + `";`).Scan(&id)
	if err != nil {
		return "", nil
	}
	return id.String(), err
}

func (db *appdbimpl) CreateUser(username string) (string, error) {
	insertUserSQL := `INSERT INTO Users(Id, Username) VALUES (?, ?);`
	statement, err := db.c.Prepare(insertUserSQL)
	if err != nil {
		return "", err
	}
	uuid, _ := uuid.NewV4()
	_, err = statement.Exec(uuid, username)
	if err != nil {
		return "", err
	}
	return uuid.String(), err
}

func (db *appdbimpl) VerifyToken(user *structures.User) bool {
	err := db.c.QueryRow(`SELECT Id, Username FROM Users WHERE Id=?;`, user.Id.Value).Scan(&user.Id.Value, &user.Username.Value)
	return err == nil
}

func (db *appdbimpl) UpdateUsername(user structures.User, new structures.NewUsername) (string, error) {
	statement, err := db.c.Prepare("UPDATE Users SET Username=? WHERE Id=?")
	if err != nil {
		return "", err
	}
	_, err = statement.Exec(new.Value, user.Id.Value)
	if err != nil {
		return "", err
	}

	return new.Value, nil
}

func (db *appdbimpl) UploadFile(file structures.Photo, user string) error {
	insertUserSQL := `INSERT INTO Photo(File, User) VALUES (?, ?);`
	statement, err := db.c.Prepare(insertUserSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec(file.Value, user)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) DeleteFile(file structures.PhotoID) error {
	var id int64
	err := db.c.QueryRow(`DELETE FROM Photo WHERE Id=? RETURNING Id;`, file.Value).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetUsers(userToSearch string, pageId int64, except string) (structures.Users, error) {
	var usersList structures.Users

	rows, err := db.c.Query(`SELECT Id, Username 
								FROM Users 
								WHERE Username != "` + except + `" AND Username LIKE '%` + userToSearch + `%' 
								LIMIT 10 OFFSET ` + strconv.FormatInt((pageId*10), 10) + `;`)
	if err != nil {
		return structures.Users{}, err
	}

	i := 0
	for rows.Next() {
		err = rows.Scan(&usersList.List[i].Id, &usersList.List[i].Username)
		if err != nil {
			return structures.Users{}, err
		}
		i++
	}
	return usersList, nil
}

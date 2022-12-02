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
	var num int64

	err := db.c.QueryRow(`SELECT COUNT(Id) 
							FROM Users 
							WHERE Username != "` + except + `" AND Username LIKE '%` + userToSearch + `%';`).Scan(&num)
	if err != nil {
		return structures.Users{}, err
	}

	rows, err := db.c.Query(`SELECT Id, Username 
								FROM Users 
								WHERE Username != "` + except + `" AND Username LIKE '%` + userToSearch + `%' 
								LIMIT 10 OFFSET ` + strconv.FormatInt((pageId*10), 10) + `;`)
	if err != nil {
		return structures.Users{}, err
	}

	if num <= 10+(pageId*10) {
		usersList.NextUsersPageId = 0
		num = num % 10
	} else {
		usersList.NextUsersPageId = pageId + 1
		num = 10
	}

	i := 0
	usersList.List = make([]structures.UserData, num)
	for rows.Next() {
		err = rows.Scan(&usersList.List[i].Id, &usersList.List[i].Username)
		if err != nil {
			return structures.Users{}, err
		}
		i++
	}
	return usersList, nil
}

func (db *appdbimpl) GetUserPage(username string, pageId int64) (structures.UserPage, error) {
	var user structures.UserPage

	err := db.c.QueryRow(`SELECT u.Id, u.Username, COUNT(p.Id) 
									FROM Users AS u LEFT JOIN Photo AS p ON u.Id=p.User
									WHERE Username="`+username+`";`).Scan(&user.Id, &user.Username, &user.PhotoCounter)
	if err != nil {
		return structures.UserPage{}, err
	}

	err = db.c.QueryRow(`SELECT COUNT(*) 
								FROM Users AS u JOIN Follows AS f ON u.Id=f.Followed
								WHERE u.Id="` + user.Id + `";`).Scan(&user.Followers)
	if err != nil {
		return structures.UserPage{}, err
	}

	err = db.c.QueryRow(`SELECT COUNT(*) 
								FROM Users AS u JOIN Follows AS f ON u.Id=f.Follow
								WHERE u.Id="` + user.Id + `";`).Scan(&user.Following)
	if err != nil {
		return structures.UserPage{}, err
	}

	rows, err := db.c.Query(`SELECT p.Id
								FROM Users AS u LEFT JOIN Photo AS p ON u.Id=p.User
								WHERE u.Id="` + user.Id + `"
								LIMIT 10 OFFSET ` + strconv.FormatInt((pageId*10), 10) + `;`)
	if err != nil {
		return structures.UserPage{}, err
	}

	num := user.PhotoCounter

	if num <= 10+(pageId*10) {
		user.NextPhotosPageId = 0
		num = num % 10
	} else {
		user.NextPhotosPageId = pageId + 1
		num = 10
	}

	i := 0
	user.Photos = make([]structures.PhotoID, num)
	if num > 0 {
		for rows.Next() {
			err = rows.Scan(&user.Photos[i].Value)
			if err != nil {
				return structures.UserPage{}, err
			}
			i++
		}
	}

	return user, nil
}

package database

import (
	"errors"
	"strconv"
	"wasa-photo/service/api/structures"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) GetUserId(username string) (string, error) {
	var id uuid.UUID
	err := db.c.QueryRow(`SELECT Id FROM Users WHERE Username=?;`, username).Scan(&id)
	if err != nil {
		return "", err
	}
	return id.String(), nil
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
	defer statement.Close()
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

	defer statement.Close()
	return new.Value, nil
}

func (db *appdbimpl) UploadFile(file structures.Image, user structures.User) error {
	insertUserSQL := `INSERT INTO Photo(File, User) VALUES (?, ?);`
	statement, err := db.c.Prepare(insertUserSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec(file.Value, user.Id.Value)
	if err != nil {
		return err
	}
	defer statement.Close()
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

func (db *appdbimpl) GetUsers(userToSearch string, pageId int64, user structures.User) (structures.Users, error) {
	var usersList structures.Users
	var num int64
	var banned int64

	err := db.c.QueryRow(`SELECT COUNT(Id) 
							FROM Users 
							WHERE Username != ? AND Username LIKE '%`+userToSearch+`%';`, user.Username.Value).Scan(&num)
	if err != nil {
		return structures.Users{}, err
	}

	err = db.c.QueryRow(`SELECT COUNT(b.User) 
							FROM Banned AS b JOIN Users AS u ON b.User=u.Id
							WHERE u.Username LIKE '%`+userToSearch+`%' AND b.Banned = ?;`, user.Id.Value).Scan(&banned)
	if err != nil {
		return structures.Users{}, err
	}

	rows, err := db.c.Query(`SELECT Id, Username 
								FROM Users 
								WHERE Username!=? AND Username LIKE '%`+userToSearch+`%' AND Id NOT IN (SELECT User FROM Banned WHERE Banned = ?)
								LIMIT 10 OFFSET ?;`, user.Username.Value, user.Id.Value, strconv.FormatInt((pageId*10), 10))
	if err != nil {
		return structures.Users{}, err
	}
	defer rows.Close()

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

	err = rows.Err()
	if err != nil {
		return structures.Users{}, err
	}

	return usersList, nil
}

func (db *appdbimpl) GetUserPage(username string, pageId int64) (structures.UserPage, error) {
	var user structures.UserPage

	err := db.c.QueryRow(`SELECT u.Id, u.Username, COUNT(p.Id) 
									FROM Users AS u LEFT JOIN Photo AS p ON u.Id=p.User
									WHERE Username=?;`, username).Scan(&user.Id, &user.Username, &user.PhotoCounter)
	if err != nil {
		return structures.UserPage{}, err
	}

	err = db.c.QueryRow(`SELECT COUNT(*) 
								FROM Users AS u JOIN Follows AS f ON u.Id=f.Followed
								WHERE u.Id=?;`, user.Id).Scan(&user.Followers)
	if err != nil {
		return structures.UserPage{}, err
	}

	err = db.c.QueryRow(`SELECT COUNT(*) 
								FROM Users AS u JOIN Follows AS f ON u.Id=f.Follow
								WHERE u.Id=?;`, user.Id).Scan(&user.Following)
	if err != nil {
		return structures.UserPage{}, err
	}

	rows, err := db.c.Query(`SELECT p.Id, p.Data
								FROM Users AS u LEFT JOIN Photo AS p ON u.Id=p.User
								WHERE u.Id=?
								ORDER BY p.Data DESC
								LIMIT 10 OFFSET ?;`, user.Id, strconv.FormatInt((pageId*10), 10))
	if err != nil {
		return structures.UserPage{}, err
	}
	defer rows.Close()

	num := user.PhotoCounter

	if num <= 10+(pageId*10) {
		user.NextPhotosPageId = 0
		num = num % 10
	} else {
		user.NextPhotosPageId = pageId + 1
		num = 10
	}

	i := 0
	user.Photos = make([]structures.PhotoDetails, num)
	if num > 0 {
		for rows.Next() {
			err = rows.Scan(&user.Photos[i].Id, &user.Photos[i].Data)
			if err != nil {
				return structures.UserPage{}, err
			}

			var image structures.Photo
			image.Id = user.Photos[i].Id
			err = db.getNumberOfLikesAndNumberOfComments(&image)
			if err != nil {
				return structures.UserPage{}, err
			}
			user.Photos[i].NumLikes = image.NumLikes
			user.Photos[i].NumComments = image.NumComments
			i++
		}
		err = rows.Err()
		if err != nil {
			return structures.UserPage{}, err
		}
	}

	return user, nil
}

func (db *appdbimpl) BanUser(username string, byUsername structures.User) error {

	usernameId, err := db.GetUserId(username)
	if err != nil && usernameId == "" {
		return errors.New("user not found")
	}

	insertUserSQL := `INSERT INTO Banned(User, Banned) VALUES (?, ?);`
	statement, err := db.c.Prepare(insertUserSQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(byUsername.Id.Value, usernameId)
	if err != nil {
		return err
	}

	deleteUserSQL := `DELETE FROM Follows WHERE (Follow=? AND Followed=?) OR (Follow=? AND Followed=?)`
	statement, err = db.c.Prepare(deleteUserSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec(byUsername.Id.Value, usernameId, usernameId, byUsername.Id.Value)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) UnbanUser(username string, byUsername structures.User) error {

	usernameId, err := db.GetUserId(username)
	if err != nil && usernameId == "" {
		return errors.New("user not found")
	}

	insertUserSQL := `DELETE FROM Banned WHERE User=? AND Banned=?;`
	statement, err := db.c.Prepare(insertUserSQL)
	if err != nil {
		return err
	}
	res, err := statement.Exec(byUsername.Id.Value, usernameId)
	if err != nil {
		return err
	}

	i, err := res.RowsAffected()
	if i == 0 {
		return errors.New("relationship not found")
	} else if err != nil {
		return err
	}

	defer statement.Close()
	return nil
}

func (db *appdbimpl) FollowUser(username string, byUsername structures.User) error {

	usernameId, err := db.GetUserId(username)
	if err != nil && usernameId == "" {
		return errors.New("user not found")
	}

	insertUserSQL := `INSERT INTO Follows(Follow, Followed) VALUES (?, ?);`
	statement, err := db.c.Prepare(insertUserSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec(byUsername.Id.Value, usernameId)
	if err != nil {
		return err
	}

	defer statement.Close()
	return nil
}

func (db *appdbimpl) UnfollowUser(username string, byUsername structures.User) error {

	usernameId, err := db.GetUserId(username)
	if err != nil && usernameId == "" {
		return errors.New("user not found")
	}

	insertUserSQL := `DELETE FROM Follows WHERE Follow=? AND Followed=? RETURNING Follow;`
	statement, err := db.c.Prepare(insertUserSQL)
	if err != nil {
		return err
	}
	res, err := statement.Exec(byUsername.Id.Value, usernameId)
	if err != nil {
		return err
	}

	i, err := res.RowsAffected()
	if i == 0 {
		return errors.New("relationship not found")
	} else if err != nil {
		return err
	}

	defer statement.Close()
	return nil
}

func (db *appdbimpl) GetFeed(user structures.User, pageId int64) (structures.Photos, error) {
	var feed structures.Photos
	var num int64

	err := db.c.QueryRow(`SELECT COUNT(*) 
								FROM Users u JOIN Follows f ON u.Id=f.Follow
											JOIN Photo p ON p.User=f.Followed
								WHERE 	u.Id=?
										AND p.User NOT IN (SELECT b.Banned 
																FROM Banned b 
																WHERE b.User=f.Follow);`, user.Id.Value).Scan(&num)
	if err != nil {
		return structures.Photos{}, err
	}
	if num == 0 {
		return structures.Photos{}, nil
	}

	rows, err := db.c.Query(`SELECT p.Id, u.Username, p.User, p.Data
								FROM Follows f JOIN Photo p ON p.User=f.Followed
											 JOIN Users u ON p.User=u.Id
								WHERE f.Follow=? AND p.User NOT IN (SELECT b.Banned FROM Banned b WHERE b.User=f.Follow)
								ORDER BY p.Data DESC
								LIMIT 10 OFFSET ?;`, user.Id.Value, strconv.FormatInt((pageId*10), 10))

	if err != nil {
		return structures.Photos{}, err
	}
	defer rows.Close()

	if num <= 10+(pageId*10) {
		feed.NextFeedPageId = 0
		num = num % 10
	} else {
		feed.NextFeedPageId = pageId + 1
		num = 10
	}

	i := 0
	feed.Post = make([]structures.Photo, num)
	for rows.Next() {
		err = rows.Scan(&feed.Post[i].Id, &feed.Post[i].Username, &feed.Post[i].Identifier, &feed.Post[i].Data)
		if err != nil {
			return structures.Photos{}, err
		}
		err = db.getNumberOfLikesAndNumberOfComments(&feed.Post[i])
		if err != nil {
			return structures.Photos{}, err
		}
		i++
	}
	err = rows.Err()
	if err != nil {
		return structures.Photos{}, err
	}

	return feed, nil
}

func (db *appdbimpl) getNumberOfLikesAndNumberOfComments(image *structures.Photo) error {
	err := db.c.QueryRow(`SELECT COUNT(IdPhoto) FROM Likes WHERE IdPhoto=?;`, strconv.Itoa(int(image.Id))).Scan(&image.NumLikes)
	if err != nil {
		return err
	}
	err = db.c.QueryRow(`SELECT COUNT(IdPhoto) FROM Comment WHERE IdPhoto=?;`, strconv.Itoa(int(image.Id))).Scan(&image.NumComments)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetPhoto(photoId int64, image *structures.Image) error {
	err := db.c.QueryRow(`SELECT File FROM Photo WHERE Id=?;`, strconv.Itoa(int(photoId))).Scan(&image.Value)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) SetLike(photoId structures.PhotoID, user structures.User) error {

	err := db.c.QueryRow(`SELECT Id FROM Photo WHERE Id=?;`, strconv.Itoa(int(photoId.Value))).Scan(&photoId.Value)
	if err != nil {
		return errors.New("image not found")
	}

	insertUserSQL := `INSERT INTO Likes(IdPhoto, User) VALUES (?, ?);`
	statement, err := db.c.Prepare(insertUserSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec(photoId.Value, user.Id.Value)
	if err != nil {
		return err
	}

	defer statement.Close()
	return nil
}

func (db *appdbimpl) RemoveLike(photoId structures.PhotoID, user structures.User) error {

	err := db.c.QueryRow(`SELECT Id FROM Photo WHERE Id=?;`, strconv.Itoa(int(photoId.Value))).Scan(&photoId.Value)
	if err != nil {
		return errors.New("image not found")
	}

	insertUserSQL := `DELETE FROM Likes WHERE IdPhoto=? AND User=?;`
	statement, err := db.c.Prepare(insertUserSQL)
	if err != nil {
		return err
	}
	res, err := statement.Exec(photoId.Value, user.Id.Value)
	if err != nil {
		return err
	}

	v, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if v == 0 {
		return errors.New("0 changes")
	}

	defer statement.Close()
	return nil
}

func (db *appdbimpl) WriteComment(photoId structures.PhotoID, user structures.User, comment structures.Comment) error {

	err := db.c.QueryRow(`SELECT Id FROM Photo WHERE Id=?;`, strconv.Itoa(int(photoId.Value))).Scan(&photoId.Value)
	if err != nil {
		return errors.New("image not found")
	}

	insertUserSQL := `INSERT INTO Comment(IdPhoto, User, Text) VALUES (?, ?, ?);`
	statement, err := db.c.Prepare(insertUserSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec(photoId.Value, user.Id.Value, comment.Text)
	if err != nil {
		return err
	}

	defer statement.Close()
	return nil
}

func (db *appdbimpl) DeleteComment(comment structures.CommentId) error {

	insertUserSQL := `DELETE FROM Comment WHERE Id=?;`
	statement, err := db.c.Prepare(insertUserSQL)
	if err != nil {
		return err
	}
	res, err := statement.Exec(comment.Value)
	if err != nil {
		return err
	}

	v, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if v == 0 {
		return errors.New("0 changes")
	}

	defer statement.Close()
	return nil
}

func (db *appdbimpl) GetComments(photoId structures.PhotoID, pageId int64, user structures.User) (structures.Comments, error) {
	var comments structures.Comments
	var num int64

	err := db.c.QueryRow(`SELECT Id FROM Photo WHERE Id=?;`, strconv.Itoa(int(photoId.Value))).Scan(&photoId.Value)
	if err != nil {
		return structures.Comments{}, errors.New("image not found")
	}
	err = db.c.QueryRow(`SELECT COUNT(*) 
								FROM Comment c JOIN Photo p ON c.IdPhoto=p.Id
								WHERE p.Id=?
										AND c.User NOT IN (SELECT b.Banned 
																FROM Banned b 
																WHERE b.User=?)
										AND c.User NOT IN (SELECT b.User
																FROM Banned b 
																WHERE b.Banned=?);`, strconv.FormatInt(photoId.Value, 10), user.Id.Value, user.Id.Value).Scan(&num)
	if err != nil {
		return structures.Comments{}, err
	}
	if num == 0 {
		return structures.Comments{}, nil
	}

	rows, err := db.c.Query(`SELECT u.Id, u.Username, c.Id, c.Data, c.Text
								FROM Comment c JOIN Photo p ON c.IdPhoto=p.Id JOIN Users u ON c.User=u.Id
								WHERE p.Id=?
										AND u.Id NOT IN (SELECT b.Banned 
															FROM Banned b 
															WHERE b.User=?)
										AND u.Id NOT IN (SELECT b.User
															FROM Banned b 
															WHERE b.Banned=?) 
								ORDER BY c.Data DESC
								LIMIT 10 OFFSET ?;`, strconv.FormatInt(photoId.Value, 10), user.Id.Value, user.Id.Value, strconv.FormatInt((pageId*10), 10))
	if err != nil {
		return structures.Comments{}, err
	}
	defer rows.Close()

	if num <= 10+(pageId*10) {
		comments.NextCommentPageId = 0
		num = num % 10
	} else {
		comments.NextCommentPageId = pageId + 1
		num = 10
	}

	i := 0
	comments.Comments = make([]structures.CommentData, num)
	for rows.Next() {
		err = rows.Scan(&comments.Comments[i].IdUser, &comments.Comments[i].Username, &comments.Comments[i].Id, &comments.Comments[i].Data, &comments.Comments[i].Text)
		if err != nil {
			return structures.Comments{}, err
		}
		i++
	}
	err = rows.Err()
	if err != nil {
		return structures.Comments{}, err
	}
	return comments, nil
}

func (db *appdbimpl) GetBanStatus(username structures.User, byUsername structures.User) error {

	var err error

	username.Id.Value, err = db.GetUserId(username.Username.Value)
	if err != nil && username.Id.Value == "" {
		return errors.New("user not found")
	}

	err = db.c.QueryRow(`SELECT * FROM Banned WHERE User=? AND Banned=?;`, byUsername.Id.Value, username.Id.Value).Scan(&byUsername.Id.Value, &username.Id.Value)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetFollowStatus(username structures.User, byUsername structures.User) error {

	var err error

	username.Id.Value, err = db.GetUserId(username.Username.Value)
	if err != nil && username.Id.Value == "" {
		return errors.New("user not found")
	}

	err = db.c.QueryRow(`SELECT Follow, Followed FROM Follows WHERE Follow=? AND Followed=?;`, byUsername.Id.Value, username.Id.Value).Scan(&byUsername.Id.Value, &username.Id.Value)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetLikeStatus(username structures.User, photoId structures.PhotoID) error {
	err := db.c.QueryRow(`SELECT IdPhoto, User FROM Likes WHERE User=? AND IdPhoto=?;`, username.Id.Value, strconv.Itoa(int(photoId.Value))).Scan(&photoId.Value, &username.Id.Value)
	if err != nil {
		return err
	}
	return nil
}

/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"wasa-photo/service/api/structures"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetUserId(username string) (string, error)
	CreateUser(username string) (string, error)
	VerifyToken(user *structures.User) bool
	UpdateUsername(user structures.User, new structures.NewUsername) (string, error)
	UploadFile(file structures.Image, user structures.User) error
	DeleteFile(file structures.PhotoID) error
	GetUsers(userToSearch string, pageId int64, except string) (structures.Users, error)
	GetUserPage(username string, pageId int64) (structures.UserPage, error)
	BanUser(username string, byUsername structures.User) error
	UnbanUser(username string, byUsername structures.User) error
	FollowUser(username string, byUsername structures.User) error
	UnfollowUser(username string, byUsername structures.User) error
	GetFeed(user structures.User, pageId int64) (structures.Photos, error)
	GetPhoto(photoId int64, image *structures.Image) error
	SetLike(photoId structures.PhotoID, user structures.User) error
	RemoveLike(photoId structures.PhotoID, user structures.User) error
	WriteComment(photoId structures.PhotoID, user structures.User, comment structures.Comment) error
	DeleteComment(comment structures.CommentId) error
	GetComments(photoId structures.PhotoID, pageId int64, user structures.User) (structures.Comments, error)
	GetBanStatus(username structures.User, byUsername structures.User) error
	GetFollowStatus(username structures.User, byUsername structures.User) error
	GetLikeStatus(username structures.User, photoId structures.PhotoID) error
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	return &appdbimpl{c: db}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

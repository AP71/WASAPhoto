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
	getNumberOfLikesAndNumberOfComments(image *structures.Photo) error
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

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE "Users" (
						"Id"	TEXT NOT NULL,
						"Username"	TEXT NOT NULL UNIQUE,
						PRIMARY KEY("Id")
					);
					
					CREATE TABLE "Follows" (
						"Follow"	TEXT NOT NULL,
						"Followed"	TEXT NOT NULL,
						"Data"	TEXT NOT NULL DEFAULT current_timestamp,
						PRIMARY KEY("Follow","Followed"),
						FOREIGN KEY("Follow") REFERENCES "Users"("Id") ON DELETE CASCADE,
						FOREIGN KEY("Followed") REFERENCES "Users"("Id") ON DELETE CASCADE
					);
					
					CREATE TABLE "Banned" (
						"User"	TEXT NOT NULL,
						"Banned"	TEXT NOT NULL,
						FOREIGN KEY("Banned") REFERENCES "Users"("Id") ON DELETE CASCADE,
						FOREIGN KEY("User") REFERENCES "Users"("Id") ON DELETE CASCADE,
						PRIMARY KEY("User","Banned")
					);
					
					CREATE TABLE "Photo" (
						"Id"	INTEGER NOT NULL,
						"File"	BLOB NOT NULL,
						"User"	TEXT NOT NULL,
						"Data"	TEXT NOT NULL DEFAULT current_timestamp,
						PRIMARY KEY("Id" AUTOINCREMENT),
						FOREIGN KEY("User") REFERENCES "Users"("Id") ON DELETE CASCADE
					);
					
					CREATE TABLE "Likes" (
						"IdPhoto"	INTEGER NOT NULL,
						"User"	TEXT NOT NULL,
						"Data"	TEXT NOT NULL DEFAULT current_timestamp,
						FOREIGN KEY("User") REFERENCES "Users"("Id") ON DELETE CASCADE,
						PRIMARY KEY("IdPhoto","User"),
						FOREIGN KEY("IdPhoto") REFERENCES "Photo"("Id") ON DELETE CASCADE
					);
					
					CREATE TABLE "Comment" (
						"Id"	INTEGER NOT NULL,
						"IdPhoto"	INTEGER NOT NULL,
						"User"	TEXT NOT NULL,
						"Data"	TEXT NOT NULL DEFAULT current_timestamp,
						"Text"	TEXT NOT NULL,
						FOREIGN KEY("User") REFERENCES "Users"("Id") ON DELETE CASCADE,
						FOREIGN KEY("IdPhoto") REFERENCES "Photo"("Id") ON DELETE CASCADE,
						PRIMARY KEY("Id" AUTOINCREMENT)
					);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, errors.New("error creating database structure: " + err.Error())
		}
	}

	return &appdbimpl{c: db}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

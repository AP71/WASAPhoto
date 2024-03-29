package structures

import (
	"net/http"
	"regexp"
)

type PhotoID struct {
	Value int64 `json:"photoId"`
}

type Image struct {
	Value []byte `json:"data"`
}

type PhotoDetails struct {
	Id          int64  `json:"photo"`
	Data        string `json:"data"`
	NumLikes    int64  `json:"numberOfLikes"`
	NumComments int64  `json:"numberOfComments"`
}

type Photo struct {
	Id          int64  `json:"photo"`
	Data        string `json:"data"`
	Username    string `json:"username"`
	Identifier  string `json:"identifier"`
	NumLikes    int64  `json:"numberOfLikes"`
	NumComments int64  `json:"numberOfComments"`
}

type Photos struct {
	Post           []Photo `json:"posts"`
	NextFeedPageId int64   `json:"nextFeedPageId"`
}

type Comment struct {
	Text string `json:"text"`
}

type CommentId struct {
	Value int64 `json:"id"`
}

type CommentData struct {
	IdUser   string `json:"idUser"`
	Username string `json:"username"`
	Id       int64  `json:"id"`
	Data     string `json:"data"`
	Text     string `json:"text"`
}

type Comments struct {
	Comments          []CommentData `json:"comments"`
	NextCommentPageId int64         `json:"nextCommentPageId"`
}

func (c *Comment) IsValid() bool {
	l := len(c.Text)
	res, _ := regexp.MatchString(`^[a-zA-Z0-9_., !?:;""$%&€()[{}]*$`, c.Text)
	return (l > 0) && (l < 256) && res
}

func (i *Image) IsValid() bool {
	l := len(i.Value)
	return (l > 0) && (l <= 15728640)
}

func (i *Image) Extension() string {
	return http.DetectContentType(i.Value)
}

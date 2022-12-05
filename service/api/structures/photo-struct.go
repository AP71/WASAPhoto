package structures

type PhotoID struct {
	Value int64 `json:"photoId"`
}

type Image struct {
	Value []byte `json:"data"`
}

type Photo struct {
	Id          int64  `json:"photo"`
	Data        string `json:"data"`
	User        string `json:"user"`
	NumLikes    int64  `json:"numberOfLikes"`
	NumComments int64  `json:"numberOfComments"`
}

type Photos struct {
	Post           []Photo `json:"posts"`
	NextFeedPageId int64   `json:"nextFeedPageId"`
}

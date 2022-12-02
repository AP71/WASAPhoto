package structures

type PhotoID struct {
	Value int64 `json:"photoId"`
}

type Photo struct {
	Value []byte `json:"data"`
}

type Post struct {
	Id          PhotoID `json:"photo"`
	User        string  `json:"user"`
	Data        string  `json:"data"`
	NumLikes    int64   `json:"numberOfLikes"`
	NumComments int64   `json:"numberOfComments"`
}

package models

import (
	"encoding/json"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	LIKE            uint8 = 0
	DISLIKE         uint8 = 1
	ADD_TO_FAVORITE uint8 = 2

	TYPE_ARTICLE uint8 = 0
	TYPE_COMMENT uint8 = 1
	TYPE_REPLY   uint8 = 2
)

type Like struct {
	UserID    int   `bson:"userID" json:"userID"`
	Position  bool  `bson:"position" json:"position"` //
	TimeStamp int64 `bson:"time" json:"time"`
}

type Favorite struct {
	UserID                 int    `bson:"userID" json:"userID"`
	TimeStamp              int64  `bson:"time" json:"time"`
	FavoriteCollectionUUID string `bson:"favoriteCollectionUUID" json:"favoriteCollectionUUID"`
}

type SimpleStats struct {
	Likes    int `bson:"likes" json:"likes"`
	Dislikes int `bson:"dislikes" json:"dislikes"`
}

type PostStats struct {
	Likes     int `bson:"likes" json:"likes"`
	Dislikes  int `bson:"dislikes" json:"dislikes"`
	Favorites int `bson:"favorites" json:"favorites"`
}

type Post struct {
	UUID       string             `bson:"uuid" json:"uuid"`
	UserID     int                `bson:"userID" json:"userID"`
	Content    string             `bson:"content" json:"content"`
	DocumentID primitive.ObjectID `bson:"_id,omitempty" json:"_oid"`
	Title      string             `bson:"title" json:"title"`
	Statistics PostStats          `bson:"statistics" json:"statistics"`
	Comments   []Comment          `bson:"comments" json:"comments"`
	Likes      []Like             `bson:"likes" json:"likes"`
	Favorites  []Favorite         `bson:"favorites" json:"favorites"`
}

func (content *Post) ToIndentedJSON() string {
	obj, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(obj)
}

// 评论
type Comment struct {
	UUID       string      `bson:"uuid" json:"uuid"`
	UserID     int         `bson:"userID" json:"userID"`
	Content    string      `bson:"content" json:"content"`
	Statistics SimpleStats `bson:"statistics" json:"statistics"`
	Replies    []Reply     `bson:"replies" json:"replies"`
	Likes      []Like      `bson:"likes" json:"likes"`
}

// 回复
type Reply struct {
	UUID       string      `bson:"uuid" json:"uuid"`
	UserID     int         `bson:"userID" json:"userID"`
	Content    string      `bson:"content" json:"content"`
	Statistics SimpleStats `bson:"statistics" json:"statistics"`
	ToUUID     string      `bson:"toUUID" json:"toUUID"`
	Likes      []Like      `bson:"likes" json:"likes"`
}

// 所需用到的请求
type NewPostRequest struct {
	Content string `json:"content"`
	Title   string `json:"title"`
	UserID  int    `json:"userID"`
}

type NewCommentRequest struct {
	PostOID primitive.ObjectID `json:"postOID"`
	Content string             `json:"content"`
	UserID  int                `json:"userID"`
}

type NewReplyRequest struct {
	PostOID     primitive.ObjectID `json:"postOID"`
	CommentUUID string             `json:"commentUUID"`
	Content     string             `json:"content"`
	ToUUID      string             `json:"toUUID"`
	UserID      int                `json:"userID"`
}

type LikePostRequest struct {
	PostOID  primitive.ObjectID `json:"postOID"`
	UserID   int                `json:"userID"`
	Position bool               `json:"position"`
}

type LikeCommentRequest struct {
	PostOID     primitive.ObjectID `json:"postOID"`
	CommentUUID string             `json:"commentUUID"`
	UserID      int                `json:"userID"`
	Position    bool               `json:"position"`
}

type LikeReplyRequest struct {
	PostOID     primitive.ObjectID `json:"postOID"`
	CommentUUID string             `json:"commentUUID"`
	ReplyUUID   string             `json:"replyUUID"`
	UserID      int                `json:"userID"`
	Position    bool               `json:"position"`
}

// 所需的方法
// statistics会被直接初始化，所以无需担心。
func NewPost(req *NewPostRequest) (post *Post) {
	post = &Post{
		UUID:       uuid.NewString(),
		UserID:     req.UserID,
		Content:    req.Content,
		Title:      req.Title,
		Comments:   []Comment{},
		Likes:      []Like{},
		Favorites:  []Favorite{},
		Statistics: PostStats{0, 0, 0},
	}
	return
}

func NewComment(req *NewCommentRequest) (comment *Comment) {
	comment = &Comment{
		UUID:       uuid.NewString(),
		UserID:     req.UserID,
		Content:    req.Content,
		Replies:    []Reply{},
		Likes:      []Like{},
		Statistics: SimpleStats{0, 0},
	}
	return
}

func NewReply(req *NewReplyRequest) (comment *Reply) {
	comment = &Reply{
		UUID:       uuid.NewString(),
		ToUUID:     req.ToUUID,
		UserID:     req.UserID,
		Content:    req.Content,
		Likes:      []Like{},
		Statistics: SimpleStats{0, 0},
	}
	return
}

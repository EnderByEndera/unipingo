package models

import (
	"melodie-site/server/utils"

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
	UserID    primitive.ObjectID `bson:"userID" json:"userID"`
	Position  bool               `bson:"position" json:"position"` //
	TimeStamp int64              `bson:"time" json:"time"`
}

type Favorite struct {
	UserID                 primitive.ObjectID `bson:"userID" json:"userID"`
	TimeStamp              int64              `bson:"time" json:"time"`
	FavoriteCollectionUUID string             `bson:"favoriteCollectionUUID" json:"favoriteCollectionUUID"`
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

type PostCategory uint

type Post struct {
	CreateTime uint64               `bson:"createTime" json:"createTime"`
	UpdateTime uint64               `bson:"updateTime" json:"updateTime"`
	ID         primitive.ObjectID   `bson:"_id,omitempty" json:"_oid"`
	UserID     primitive.ObjectID   `bson:"userID" json:"userID"`
	Category   PostCategory         `bson:"category" json:"category"`
	BelongsTo  EntityWithName       `bson:"belongsTo" json:"belongsTo"`
	Topics     []string             `bson:"topics" json:"topics"`
	Statistics PostStats            `bson:"statistics" json:"statistics"`
	Liked      []primitive.ObjectID `bson:"liked" json:"liked"`         // 点赞的用户
	Favorited  []primitive.ObjectID `bson:"favorited" json:"favorited"` // 添加到收藏的用户
	Title      string               `bson:"title" json:"title"`
	Content    string               `bson:"content" json:"content"`
}

func (content *Post) ToIndentedJSON() string {
	return utils.ToIndentedJSON(content)
}

// // 评论
// type Comment struct {
// 	UUID       string             `bson:"uuid" json:"uuid"`
// 	TimeStamp  uint64             `bson:"timeStamp" json:"timeStamp"`
// 	UserID     primitive.ObjectID `bson:"userID" json:"userID"`
// 	Content    string             `bson:"content" json:"content"`
// 	Statistics SimpleStats        `bson:"statistics" json:"statistics"`
// 	Replies    []Reply            `bson:"replies" json:"replies"`
// 	Likes      []Like             `bson:"likes" json:"likes"`
// }

// // 回复
// type Reply struct {
// 	UUID       string             `bson:"uuid" json:"uuid"`
// 	TimeStamp  uint64             `bson:"timeStamp" json:"timeStamp"`
// 	UserID     primitive.ObjectID `bson:"userID" json:"userID"`
// 	Content    string             `bson:"content" json:"content"`
// 	Statistics SimpleStats        `bson:"statistics" json:"statistics"`
// 	ToUUID     string             `bson:"toUUID" json:"toUUID"`
// 	Likes      []Like             `bson:"likes" json:"likes"`
// }

// // 所需用到的请求
// type NewPostRequest struct {
// 	Content string `json:"content"`
// 	Title   string `json:"title"`
// 	UserID  string `json:"userID"`
// }

// type NewCommentRequest struct {
// 	PostOID string `json:"postOID"`
// 	Content string `json:"content"`
// 	UserID  string `json:"userID"`
// }

// type NewReplyRequest struct {
// 	PostOID     string `json:"postOID"`
// 	CommentUUID string `json:"commentUUID"`
// 	Content     string `json:"content"`
// 	ToUUID      string `json:"toUUID"`
// 	UserID      string `json:"userID"`
// }

// type LikePostRequest struct {
// 	PostOID  string `json:"postOID"`
// 	UserID   string `json:"userID"`
// 	Position bool   `json:"position"`
// }

// type LikeCommentRequest struct {
// 	PostOID     string `json:"postOID"`
// 	UserID      string `json:"userID"`
// 	Position    bool   `json:"position"`
// 	CommentUUID string `json:"commentUUID"`
// }

// type LikeReplyRequest struct {
// 	PostOID     string `json:"postOID"`
// 	UserID      string `json:"userID"`
// 	Position    bool   `json:"position"`
// 	CommentUUID string `json:"commentUUID"`
// 	ReplyUUID   string `json:"replyUUID"`
// }

// // 返回的结构体
// type PostOutline struct {
// 	PostOID    string         `json:"id"`
// 	TimeStamp  uint64         `json:"timeStamp"`
// 	User       UserPublicInfo `json:"user"`
// 	Title      string         `json:"title"`
// 	Content    string         `json:"content"`
// 	Statistics PostStats      `json:"statistics"`
// }

// // type LikePostActionMeta struct {
// // 	PostOID  primitive.ObjectID `json:"postOID"`
// // 	UserID   primitive.ObjectID `json:"userID"`
// // 	Position bool               `json:"position"`
// // }

// // type LikeCommentActionMeta struct {
// // 	PostOID     primitive.ObjectID `json:"postOID"`
// // 	UserID      primitive.ObjectID `json:"userID"`
// // 	Position    bool               `json:"position"`
// // 	CommentUUID string             `json:"commentUUID"`
// // }

// // type LikeReplyActionMeta struct {
// // 	PostOID     primitive.ObjectID `json:"postOID"`
// // 	UserID      primitive.ObjectID `json:"userID"`
// // 	Position    bool               `json:"position"`
// // 	CommentUUID string             `json:"commentUUID"`
// // 	ReplyUUID   string             `json:"replyUUID"`
// // }

// // func (req *LikePostRequest) ToActionMeta() (meta *LikePostActionMeta, err error) {
// // 	postOID, err := primitive.ObjectIDFromHex(req.PostOID)
// // 	if err != nil {
// // 		return
// // 	}
// // 	userID, err := primitive.ObjectIDFromHex(req.UserID)
// // 	meta = &LikePostActionMeta{
// // 		PostOID: postOID,
// // 		UserID:  userID,
// // 		Position: ,
// // 	}
// // 	return
// // }

// // 所需的方法
// // statistics会被直接初始化，所以无需担心。
// func NewPost(req *NewPostRequest) (post *Post, err error) {
// 	userID, err := primitive.ObjectIDFromHex(req.UserID)
// 	if err != nil {
// 		return
// 	}
// 	post = &Post{
// 		UUID:       uuid.NewString(),
// 		UserID:     userID,
// 		Content:    req.Content,
// 		Title:      req.Title,
// 		Comments:   []Comment{},
// 		Likes:      []Like{},
// 		Favorites:  []Favorite{},
// 		Statistics: PostStats{0, 0, 0},
// 	}
// 	return
// }

// func NewComment(req *NewCommentRequest) (comment *Comment, err error) {
// 	userID, err := primitive.ObjectIDFromHex(req.UserID)
// 	if err != nil {
// 		return
// 	}
// 	comment = &Comment{
// 		UUID:       uuid.NewString(),
// 		UserID:     userID,
// 		Content:    req.Content,
// 		Replies:    []Reply{},
// 		Likes:      []Like{},
// 		Statistics: SimpleStats{0, 0},
// 	}
// 	return
// }

// func NewReply(req *NewReplyRequest) (comment *Reply) {
// 	userID, err := primitive.ObjectIDFromHex(req.UserID)
// 	if err != nil {
// 		return
// 	}
// 	comment = &Reply{
// 		UUID:       uuid.NewString(),
// 		ToUUID:     req.ToUUID,
// 		UserID:     userID,
// 		Content:    req.Content,
// 		Likes:      []Like{},
// 		Statistics: SimpleStats{0, 0},
// 	}
// 	return
// }

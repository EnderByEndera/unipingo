package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	StudentIdentityUnauthenticated int = 0 //
	StudentIdentityPhotoUploaded   int = 1 // 照片已上传，等待管理员审核
	StudentIdentityAuthenticated   int = 2 // 审核通过
	StudentIdentityUnQualified     int = 3 // 被打回
)

type StudentIdentityAuthenticationPhotographs struct {
	IDCard struct {
		Obverse string `json:"obverse" bson:"obverse"`
		Reverse string `json:"reverse" bson:"reverse"`
	} `json:"idCard" bson:"idCard"`
	StudentID struct {
		Cover string `json:"cover" bson:"cover"` // 学生证封面
		Inner string `json:"inner" bson:"inner"` // 学生证内页
	} `json:"studentID" bson:"studentID"`
}

// 学生认证流程
type StudentIdentityAuthentication struct {
	UserID      primitive.ObjectID                       `json:"userID" bson:"userID"`
	SchoolName  string                                   `json:"schoolName" bson:"schoolName"`
	Status      int                                      `json:"status" bson:"status"`
	Photographs StudentIdentityAuthenticationPhotographs `json:"photos" bson:"photos"`
	Suggestions string                                   `json:"suggestions" bson:"suggestions"`
}

type NewStudentIdentityAuthenticationRequest struct {
	SchoolName  string                                   `json:"schoolName"`
	Photographs StudentIdentityAuthenticationPhotographs `json:"photos"`
}

type ModifyStuIDAuthStatRequest struct {
	UserID      string `json:"userID"`
	Status      int    `json:"status"`
	Suggestions string `json:"suggestions"`
}

func (req *NewStudentIdentityAuthenticationRequest) ToAuthStruct(userID string) (auth StudentIdentityAuthentication, err error) {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return
	}
	auth = StudentIdentityAuthentication{
		UserID:      oid,
		Photographs: req.Photographs,
		Status:      StudentIdentityPhotoUploaded,
		Suggestions: "",
		SchoolName:  req.SchoolName,
	}
	return
}

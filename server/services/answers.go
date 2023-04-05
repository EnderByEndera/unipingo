package services

import (
	"context"
	"errors"
	"fmt"
	"melodie-site/server/db"
	"melodie-site/server/models"
	"melodie-site/server/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var answerService *AnswerService

type AnswerService struct {
	userAndAnswerMutex *utils.KeyedMutex
}

type IDInterface struct {
	UserID   primitive.ObjectID
	AnswerID primitive.ObjectID
}

func (ansService *AnswerService) LockUserAndAnswer(userID, ansID primitive.ObjectID) {
	key := IDInterface{UserID: userID, AnswerID: ansID}
	ansService.userAndAnswerMutex.Lock(key)
}

func (ansService *AnswerService) UnlockUserAndAnswer(userID, ansID primitive.ObjectID) {
	key := IDInterface{UserID: userID, AnswerID: ansID}
	ansService.userAndAnswerMutex.Unlock(key)
}

type ApprovedStatus uint64

const NotApprovedOrDisapproved ApprovedStatus = 0
const AlreadyApproved ApprovedStatus = 1
const AlreadyDisapproved ApprovedStatus = 2

func (service *AnswerService) NewAnswer(answer *models.Answer) (docID primitive.ObjectID, err error) {
	// conn := db.GetMongoConn()
	answer.Init()
	res1, err := db.GetCollection("answers").InsertOne(context.Background(), answer)
	if err != nil {
		return
	}
	docID = res1.InsertedID.(primitive.ObjectID)
	return
}

func (service *AnswerService) GetAnswerByID(oid primitive.ObjectID) (answer *models.Answer, err error) {
	answer = &models.Answer{}

	filter := bson.D{{"_id", oid}}
	err = db.GetCollection("answers").FindOne(context.TODO(), filter).Decode(answer)
	return
}

func (service *AnswerService) DeleteAnswerByID(oid primitive.ObjectID) (err error) {
	filter := bson.D{{"_id", oid}}
	_, err = db.GetCollection("answers").DeleteOne(context.TODO(), filter)
	return
}

func exists(filter bson.M) bool {
	result := db.GetCollection("answers").FindOne(context.TODO(), filter)
	return result.Err() == nil
}

func (service *AnswerService) CheckIfAlreadyLiked(ansID, userID primitive.ObjectID) (stat ApprovedStatus) {
	filterCreator := func(key string) bson.M {

		return bson.M{
			"_id": ansID,
			key: bson.M{
				"$elemMatch": bson.M{
					"$eq": userID,
				},
			},
		}
	}
	approved, disapproved := exists(filterCreator("approvedUsers")), exists(filterCreator("disapprovedUsers"))
	if (!approved) && (!disapproved) {
		return NotApprovedOrDisapproved
	} else if approved && (!disapproved) {
		return AlreadyApproved
	} else if (!approved) && disapproved {
		return AlreadyDisapproved
	} else {
		panic("this answer has been approved and disapproved at the same time!")
	}
}

func (service *AnswerService) GetLikedStatus(ansID, userID primitive.ObjectID) {}

func (service *AnswerService) AnswerExists(ansID primitive.ObjectID) bool {
	res := db.GetCollection("answers").FindOne(context.TODO(), bson.M{"_id": ansID})
	return res.Err() == nil
}

// 伪三目运算符。
func pseudoTernaryOp[T any](condition bool, valueOnTrue, valueOnFalse T) T {
	if condition {
		return valueOnTrue
	} else {
		return valueOnFalse
	}
}

// 取消赞
// 如果未点赞也未点踩，返回“未点赞”错误
// 如果已点赞，则尝试取消赞操作。成功则返回“取消赞成功”,失败则返回数据库操作失败。
// 如果已经点踩，返回“回答已经踩过”
func (service *AnswerService) cancelLikeInAnswer(userID primitive.ObjectID, ansID primitive.ObjectID) models.StatusReport {
	likedStatus := service.CheckIfAlreadyLiked(ansID, userID)
	answer, err := service.GetAnswerByID(ansID)
	if err != nil {
		return models.StatusReport{err, models.ApproveAnswerStatus.AnswerDoesNotExist}
	}
	isAlumn, _ := GetAuthService().IsAlumn(userID, answer.BelongsTo.ID)
	var statement bson.M
	var filter bson.M
	if likedStatus == AlreadyApproved {
		statement = bson.M{
			"$pull": bson.M{"approvedUsers": userID},
			"$inc": bson.M{
				"statistics.approves":      -1,
				"statistics.alumnApproves": pseudoTernaryOp(isAlumn, -1, 0),
			},
		}
		filter = bson.M{
			"_id": ansID,
		}
	} else if likedStatus == NotApprovedOrDisapproved {
		err = errors.New("answer not approved!")
		return models.StatusReport{err, models.ApproveAnswerStatus.AnswerNotApproved}
	} else {
		err = errors.New("already disapproved this answer, cannot cancel approval!")
		return models.StatusReport{err, models.ApproveAnswerStatus.AnswerAlreadyDisapproved}
	}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	res := db.GetCollection("answers").FindOneAndUpdate(context.TODO(), filter, statement, opts)
	err = res.Err()
	if err != nil {
		return models.StatusReport{err, models.ApproveAnswerStatus.OperationFailed}
	} else {
		return models.StatusReport{err, models.ApproveAnswerStatus.CancelApproveSucceeded}
	}
}

// 取消踩
// 如果未点赞也未点踩，返回“未点踩”错误
// 如果已点踩，则尝试取消踩操作。成功则返回“取消踩成功”,失败则返回数据库操作失败。
// 如果已经点赞，返回“回答已经赞过”
func (service *AnswerService) cancelDislikeInAnswer(userID primitive.ObjectID, ansID primitive.ObjectID) models.StatusReport {
	likedStatus := service.CheckIfAlreadyLiked(ansID, userID)
	answer, err := service.GetAnswerByID(ansID)
	if err != nil {
		return models.StatusReport{err, models.ApproveAnswerStatus.AnswerDoesNotExist}
	}
	isAlumn, _ := GetAuthService().IsAlumn(userID, answer.BelongsTo.ID)
	var statement bson.M
	var filter bson.M
	if likedStatus == AlreadyDisapproved {
		statement = bson.M{
			"$pull": bson.M{"approvedUsers": userID},
			"$inc": bson.M{
				"statistics.disapproves":      -1,
				"statistics.alumnDisapproves": pseudoTernaryOp(isAlumn, -1, 0),
			},
		}
		filter = bson.M{
			"_id": ansID,
		}
	} else if likedStatus == NotApprovedOrDisapproved {
		err = errors.New("answer not disapproved!")
		return models.StatusReport{err, models.ApproveAnswerStatus.AnswerNotDisapproved}
	} else {
		err = errors.New("already approved this answer, cannot cancel disapproval!")
		return models.StatusReport{err, models.ApproveAnswerStatus.AnswerAlreadyApproved}
	}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	res := db.GetCollection("answers").FindOneAndUpdate(context.TODO(), filter, statement, opts)
	err = res.Err()
	if err != nil {
		return models.StatusReport{err, models.ApproveAnswerStatus.OperationFailed}
	} else {
		return models.StatusReport{err, models.ApproveAnswerStatus.CancelDisapproveSucceeded}
	}
}

// 如果赞过了，就返回
func (service *AnswerService) giveLikeToAnswer(userID primitive.ObjectID, ansID primitive.ObjectID) models.StatusReport {
	likedStatus := service.CheckIfAlreadyLiked(ansID, userID)
	answer, err := service.GetAnswerByID(ansID)
	if err != nil {
		return models.StatusReport{err, models.ApproveAnswerStatus.AnswerDoesNotExist}
	}
	isAlumn, _ := GetAuthService().IsAlumn(userID, answer.BelongsTo.ID)

	var statement bson.M
	var filter bson.M
	if likedStatus == AlreadyApproved {
		err = errors.New("already approved this answer")
		return models.StatusReport{err, models.ApproveAnswerStatus.AnswerAlreadyApproved}
	} else if likedStatus == NotApprovedOrDisapproved {
		statement = bson.M{
			"$push": bson.M{"approvedUsers": userID},
			"$inc": bson.M{
				"statistics.approves":      1,
				"statistics.alumnApproves": pseudoTernaryOp(isAlumn, 1, 0),
			},
		}
		filter = bson.M{
			"_id": ansID,
		}
	} else {
		statement = bson.M{
			"$pull": bson.M{"disapprovedUsers": userID},
			"$push": bson.M{"approvedUsers": userID},
			"$inc": bson.M{
				"statistics.approves":         1,
				"statistics.disapproves":      -1,
				"statistics.alumnApproves":    pseudoTernaryOp(isAlumn, 1, 0),
				"statistics.alumnDisapproves": pseudoTernaryOp(isAlumn, -1, 0),
			},
		}
		filter = bson.M{
			"_id": ansID,
		}
	}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	res := db.GetCollection("answers").FindOneAndUpdate(context.TODO(), filter, statement, opts)
	err = res.Err()
	if err != nil {
		return models.StatusReport{err, models.ApproveAnswerStatus.OperationFailed}
	} else {
		return models.StatusReport{err, models.ApproveAnswerStatus.ApproveAnswerSucceeded}
	}
}

// 点踩
func (service *AnswerService) giveDislikeToAnswer(userID primitive.ObjectID, ansID primitive.ObjectID) models.StatusReport {
	likedStatus := service.CheckIfAlreadyLiked(ansID, userID)
	answer, err := service.GetAnswerByID(ansID)
	if err != nil {
		return models.StatusReport{err, models.ApproveAnswerStatus.AnswerDoesNotExist}
	}
	isAlumn, err := GetAuthService().IsAlumn(userID, answer.BelongsTo.ID)

	fmt.Printf("liked: %+v\n", likedStatus)
	var statement bson.M
	var filter bson.M
	if likedStatus == AlreadyDisapproved {
		err = errors.New("already disapproved this answer")
		return models.StatusReport{err, models.ApproveAnswerStatus.AnswerAlreadyDisapproved}
	} else if likedStatus == NotApprovedOrDisapproved {
		statement = bson.M{
			"$push": bson.M{"disapprovedUsers": userID},
			"$inc": bson.M{
				"statistics.disapproves":      1,
				"statistics.alumnDisapproves": pseudoTernaryOp(isAlumn, 1, 0),
			},
		}
		filter = bson.M{
			"_id": ansID,
		}
	} else {
		statement = bson.M{
			"$pull": bson.M{"approvedUsers": userID},
			"$push": bson.M{"disapprovedUsers": userID},
			"$inc": bson.M{
				"statistics.approves":         -1,
				"statistics.disapproves":      1,
				"statistics.alumnApproves":    pseudoTernaryOp(isAlumn, -1, 0),
				"statistics.alumnDisapproves": pseudoTernaryOp(isAlumn, 1, 0),
			},
		}
		filter = bson.M{
			"_id": ansID,
		}
	}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	res := db.GetCollection("answers").FindOneAndUpdate(context.TODO(), filter, statement, opts)
	err = res.Err()
	if err != nil {
		return models.StatusReport{err, models.ApproveAnswerStatus.OperationFailed}
	} else {
		return models.StatusReport{err, models.ApproveAnswerStatus.DisapproveAnswerSucceeded}
	}
}

func (ansService *AnswerService) ApproveAnswer(userID, ansID primitive.ObjectID) models.StatusReport {
	ansService.LockUserAndAnswer(userID, ansID)
	defer ansService.UnlockUserAndAnswer(userID, ansID)
	return ansService.giveLikeToAnswer(userID, ansID)
}

func (ansService *AnswerService) DisApproveAnswer(userID, ansID primitive.ObjectID) models.StatusReport {
	ansService.LockUserAndAnswer(userID, ansID)
	defer ansService.UnlockUserAndAnswer(userID, ansID)
	return ansService.giveDislikeToAnswer(userID, ansID)
}

func (ansService *AnswerService) CancelApprovalOfAnswer(userID, ansID primitive.ObjectID) models.StatusReport {
	ansService.LockUserAndAnswer(userID, ansID)
	defer ansService.UnlockUserAndAnswer(userID, ansID)
	return ansService.cancelLikeInAnswer(userID, ansID)
}

func (ansService *AnswerService) CancelDisApprovalOfAnswer(userID, ansID primitive.ObjectID) models.StatusReport {
	ansService.LockUserAndAnswer(userID, ansID)
	defer ansService.UnlockUserAndAnswer(userID, ansID)
	return ansService.cancelDislikeInAnswer(userID, ansID)
}

func GetAnswersService() *AnswerService {
	if answerService == nil {
		answerService = &AnswerService{userAndAnswerMutex: &utils.KeyedMutex{}}
	}
	return answerService
}

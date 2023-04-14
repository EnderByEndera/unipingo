package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"melodie-site/server/models"
	"melodie-site/server/services"
	"sync"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TCreateTempUser(t *testing.T) (user *models.User, hei *models.HEI) {
	hei, err := services.GetHEIService().GetHEIByName("北京航空航天大学")
	fmt.Printf("%+v\n", hei)
	assert.Equal(t, err, nil)
	userName := "user-" + uuid.NewString()
	user_, err := services.GetAuthService().InternalAddUser(userName, "123456", models.RoleUnpaidUser, func(u *models.User) {
		u.EducationalBackground = make([]models.EduBGItem, 0)
		u.EducationalBackground = append(u.EducationalBackground, models.EduBGItem{HEIID: hei.ID, HEIName: hei.Name})
	})
	user = &user_
	fmt.Printf("%+v %+v\n", user, user.ID)
	assert.Equal(t, err, nil)
	isAlumn, err := services.GetAuthService().IsAlumn(user.ID, hei.ID)

	assert.Equal(t, err, nil)
	assert.Equal(t, isAlumn, true)
	return
}

func TGetAnswer(t *testing.T, answerID primitive.ObjectID) (ans *models.Answer) {
	ans, err := services.GetAnswersService().GetAnswerByID(answerID)
	assert.Equal(t, err, nil)
	return
}

func TNewAnswer(t *testing.T, user models.User, hei models.HEI) (answer *models.Answer) {
	byts, _ := ioutil.ReadFile("answer.json")
	answer = &models.Answer{
		UserID:     user.ID,
		CreateTime: uint64(time.Now().UnixMilli()),
		// Category: ,
		BelongsTo: models.EntityWithName{Name: hei.Name, ID: hei.ID},
	}

	err := json.Unmarshal(byts, answer)
	if err != nil {
		fmt.Println("marshal failed:", err)
		t.FailNow()
	}
	insertedDocID, err := services.GetAnswersService().NewAnswer(answer)
	if err != nil {
		t.FailNow()
	}
	answer, err = services.GetAnswersService().GetAnswerByID(insertedDocID)
	return
}

func TGiveLike(t *testing.T, userID primitive.ObjectID, answerID primitive.ObjectID) (err error) {
	// 模拟点赞
	err = services.GetAnswersService().ApproveAnswer(userID, answerID).Error
	ans := TGetAnswer(t, answerID)
	assert.Equal(t, err, nil)
	assert.Equal(t, ans.Statistics.Approves, 1)
	assert.Equal(t, ans.Statistics.Disapproves, 0)
	assert.Equal(t, ans.Statistics.AlumnApproves, 1)
	assert.Equal(t, ans.Statistics.AlumnDisapproves, 0)

	// 重复点赞，假设点赞不成功。
	stat := services.GetAnswersService().ApproveAnswer(userID, answerID)
	err = stat.Error
	assert.Equal(t, stat.Status, models.ApproveAnswerStatus.AnswerAlreadyApproved)
	assert.NotEqual(t, err, nil)
	ans = TGetAnswer(t, answerID)
	assert.Equal(t, ans.Statistics.Approves, 1)
	assert.Equal(t, ans.Statistics.Disapproves, 0)
	assert.Equal(t, ans.Statistics.AlumnApproves, 1)
	assert.Equal(t, ans.Statistics.AlumnDisapproves, 0)

	// 取消赞
	stat = services.GetAnswersService().CancelApprovalOfAnswer(userID, answerID)
	assert.Equal(t, stat.Status, models.ApproveAnswerStatus.CancelApproveSucceeded)
	assert.Equal(t, stat.Error, nil)
	ans = TGetAnswer(t, answerID)
	fmt.Printf("%+v\n", ans)
	assert.Equal(t, ans.Statistics.Approves, 0)
	assert.Equal(t, ans.Statistics.AlumnApproves, 0)
	assert.Equal(t, len(ans.ApprovedUsers), 0)

	// 点个踩
	stat = services.GetAnswersService().DisApproveAnswer(userID, answerID)
	err = stat.Error
	assert.Equal(t, stat.Status, models.ApproveAnswerStatus.DisapproveAnswerSucceeded)
	assert.Equal(t, err, nil)
	ans = TGetAnswer(t, answerID)
	assert.Equal(t, ans.Statistics.Approves, 0)
	assert.Equal(t, ans.Statistics.Disapproves, 1)
	assert.Equal(t, ans.Statistics.AlumnApproves, 0)
	assert.Equal(t, ans.Statistics.AlumnDisapproves, 1)

	// 再点个赞。
	stat = services.GetAnswersService().ApproveAnswer(userID, answerID)
	err = stat.Error
	assert.Equal(t, stat.Status, models.ApproveAnswerStatus.ApproveAnswerSucceeded)
	assert.Equal(t, err, nil)
	ans = TGetAnswer(t, answerID)
	assert.Equal(t, ans.Statistics.Approves, 1)
	assert.Equal(t, ans.Statistics.Disapproves, 0)
	assert.Equal(t, ans.Statistics.AlumnApproves, 1)
	assert.Equal(t, ans.Statistics.AlumnDisapproves, 0)
	return
}

func TReleaseSource(t *testing.T, answer *models.Answer, user *models.User) {
	answerPtr, err := services.GetAnswersService().GetAnswerByID(answer.ID)
	assert.Equal(t, err, nil)
	fmt.Printf("%+v\n", answerPtr)

	services.GetAnswersService().DeleteAnswerByID(answer.ID)

	err = services.GetAuthService().InternalRemoveUser(user.Name)
	assert.Equal(t, err, nil)

	_, err = services.GetAnswersService().GetAnswerByID(answer.ID)
	assert.NotEqual(t, err, nil)
}

// 测试并发写入与加锁。
func TParallel(t *testing.T, answerID primitive.ObjectID, heiID primitive.ObjectID) {

	const UserNum = 20
	const chLength = UserNum * 10
	users := make([]*models.User, 0)
	ch := make(chan *models.User, chLength)
	mtx := sync.Mutex{}
	var counter = 0

	for i := 0; i < 4; i++ {
		go func() {
			for {
				u := <-ch
				// for i:=0;
				services.GetAnswersService().ApproveAnswer(u.ID, answerID)
				mtx.Lock()
				counter += 1
				mtx.Unlock()
			}
		}()
	}

	fmt.Println("Put!")
	for i := 0; i < UserNum; i++ {
		user, _ := TCreateTempUser(t)
		users = append(users, user)
		// fmt.Println("put", i)
		for j := 0; j < 10; j++ {
			ch <- user
		}
	}
	for {
		if counter == chLength {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	for i := 0; i < UserNum; i++ {
		err := services.GetAuthService().InternalRemoveUserByID(users[i].ID)
		assert.Equal(t, err, nil)
	}
	ans := TGetAnswer(t, answerID)
	fmt.Println(ans.Statistics)
}
//测试的函数名称要以Test打头
func TestModelAnswers(t *testing.T) {

	// 创建一名临时用户，属于某个学校
	user, hei := TCreateTempUser(t)

	// 模拟此用户，创建一个新的回答
	answer := TNewAnswer(t, *user, *hei)

	// 模拟点赞操作
	TGiveLike(t, user.ID, answer.ID)

	TParallel(t, answer.ID, hei.ID)
	defer TReleaseSource(t, answer, user)
}

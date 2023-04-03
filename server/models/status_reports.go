package models

type ProgramStatusType string

// 点赞或者点踩的时候，返回的数据
var ApproveAnswerStatus = struct {
	ApproveAnswerSucceeded    ProgramStatusType
	DisapproveAnswerSucceeded ProgramStatusType
	AnswerDoesNotExist        ProgramStatusType
	AnswerAlreadyApproved     ProgramStatusType
	AnswerAlreadyDisapproved  ProgramStatusType
	OperationFailed           ProgramStatusType
}{
	ApproveAnswerSucceeded:    "approve_answer_succeeded",
	DisapproveAnswerSucceeded: "disapprove_answer_succeeded",
	AnswerDoesNotExist:        "answer_does_not_exist",
	AnswerAlreadyApproved:     "answer_already_approved",
	AnswerAlreadyDisapproved:  "answer_already_disapproved",
	OperationFailed:           "operation_failed",
}

// var UserServiceStatus = struct{

// }

type StatusReport struct {
	Error  error             `json:"error"`
	Status ProgramStatusType `json:"status"`
}

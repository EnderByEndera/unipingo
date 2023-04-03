package models

type ProgramStatusType string

// 点赞或者点踩的时候，返回的数据
var ApproveAnswerStatus = struct {
	ApproveAnswerSucceeded    ProgramStatusType
	DisapproveAnswerSucceeded ProgramStatusType
	CancelApproveSucceeded    ProgramStatusType
	CancelDisapproveSucceeded ProgramStatusType
	AnswerDoesNotExist        ProgramStatusType
	AnswerAlreadyApproved     ProgramStatusType
	AnswerAlreadyDisapproved  ProgramStatusType
	OperationFailed           ProgramStatusType
}{
	ApproveAnswerSucceeded:    "approve_answer_succeeded",    // 赞同成功
	DisapproveAnswerSucceeded: "disapprove_answer_succeeded", // 点踩成功
	CancelApproveSucceeded:    "cancel_approve_succeeded",    // 取消赞成功
	CancelDisapproveSucceeded: "cancel_disapprove_succeeded", // 取消踩成功
	AnswerDoesNotExist:        "answer_does_not_exist",       // 回答不存在
	AnswerAlreadyApproved:     "answer_already_approved",     // 回答已经赞过
	AnswerAlreadyDisapproved:  "answer_already_disapproved",  // 回答已经踩过
	OperationFailed:           "operation_failed",            // 内部错误，操作失败
}

type StatusReport struct {
	Error  error             `json:"error"`
	Status ProgramStatusType `json:"status"`
}

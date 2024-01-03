package server

type UserReq struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

type LoginResp struct {
	Token  string `json:"token"`
	UserId uint64 `json:"userid"`
}

type NoteReq struct {
	Note string `json:"note"`
}

type ShareNoteReq struct {
	SenderId   uint64 `json:"senderid"`
	RecieverId uint64 `json:"recieverid"`
}

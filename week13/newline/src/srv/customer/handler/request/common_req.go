package request

type MsgMap struct {
	PaMsg PaMsg
	YzMsg YzMsg
	WmMsg WmMsg
	MbMsg MbMsg
	DzReq DzInfo
}

type Tag struct {
	Id   int32
	Name string
}

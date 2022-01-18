package request

import "time"

type EmBasicInfo struct {
	HeadImageUrl string `json:"avatar"`
	NickName     string `json:"nick"`
	Gender       int    `json:"sex"`
	Phone        string `json:"mobile"`
	Info         struct {
		Birthday       *time.Time `json:"birthday"`
		ContactAddress struct {
			Province string `json:"province"`
			City     string `json:"city"`
			District string `json:"county"`
			Address  string `json:"address"`
		} `json:"contact_address"`
	} `json:"info"`
	UnionId string `json:"union_id" gorm:"size:255;not null;unique_index;"      comment:"微信开放平台的唯一身份标识"`
}

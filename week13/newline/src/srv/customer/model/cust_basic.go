package model

import (
	"bytes"
	"github.com/corona10/goimagehash"
	"github.com/golang/protobuf/ptypes/timestamp"
	"go.uber.org/zap"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"newline.com/newline/src/common/log"
	"newline.com/newline/src/common/utils"
	"newline.com/newline/src/models"
	"strconv"
	"time"
)

//客户基本信息表
type CustBasicInfo struct {
	models.Model     `table-comment:"用户基础信息"`
	HeadImageUrl     string `json:"head_image_url" gorm:"type:varchar(255)"     comment:"微信头像地址" `
	HeadHash         string `json:"head_hash" gorm:"type:varchar(255);index;"       comment:"微信头像hash"`
	NickName         string `json:"nick_name" gorm:"type:varchar(50);index;"      comment:"微信昵称"`
	Gender           int32  `json:"gender" gorm:"type:int(2)" comment:"性别"`
	Phone            string `json:"phone" gorm:"type:varchar(50);"      comment:"手机号"`
	Country          string `json:"country" gorm:"type:varchar(50);"        comment:"国家"`
	Province         string `json:"province" gorm:"type:varchar(50);"        comment:"省"`
	City             string `json:"city" gorm:"type:varchar(50);"        comment:"市"`
	District         string `json:"district" gorm:"type:varchar(50);"        comment:"区"`
	Address          string `json:"address" gorm:"type:varchar(500);"       comment:"详细地址"`
	Birthday         string `json:"birthday" gorm:"type:varchar(50);"     comment:"出生日期"`
	UnionId          string `json:"union_id" gorm:"size:255;not null;unique_index;"      comment:"开放平台的唯一身份标识"`
	CreatedBy        string `json:"created_by" gorm:"size:10;"        comment:"用户创建渠道"`
	UpdatedBy        string `json:"updated_by" gorm:"size:10;"        comment:"用户创建渠道"`
	OwnerAccountId   int32  `json:"owner_account_id"  comment:"归属账户ID"`
	OwnerAccountName string `json:"owner_account_name" gorm:"type:varchar(50);"  comment:"归属账户ID名称"`
	OwnerFrom        string `json:"owner_from" gorm:"type:varchar(10);"  comment:"归属账户来源"`
	//audited.AuditedModel
	//loggable.LoggableModel
}

// TableName TableName
func (*CustBasicInfo) TableName() string {
	return "cust_basic_info"
}

// 用于orm查询时的模型
type CustomerBasicParam struct {
	CustBasicInfo
}

func (u *CustBasicInfo) BeforeSave() (err error) {
	if u.HeadImageUrl != "" {
		headImageUrl := utils.ReplaceLastWord(u.HeadImageUrl, "/96", "/132")
		res, err := http.Get(headImageUrl)
		if err != nil {
			//log.GetLogger().Error("Get Img Hash Step1", zap.Error(err))
			time.Sleep(time.Second * 10)
			res, err = http.Get(headImageUrl)
			if err != nil {
				log.GetLogger().Error("headImageUrl get error", zap.String("headImageUrl", headImageUrl), zap.Error(err))
				return err
			}
			return err
		}
		img, err := typeImg(res)
		if err != nil {
			log.GetLogger().Error("Get Img Hash Step2", zap.Error(err), zap.String("headImageUrl", headImageUrl))
			return err
		}
		defer res.Body.Close()
		hash, err := goimagehash.AverageHash(img)
		if err != nil {
			log.GetLogger().Error("Get Img Hash Step3", zap.Error(err))
			return err
		}
		u.HeadHash = strconv.Itoa(int(hash.GetHash()))
	}
	return
}
func typeImg(resp *http.Response) (image.Image, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.GetLogger().Error("typeImg", zap.Error(err))
	}
	pngI, err := png.Decode(bytes.NewReader(body))
	if err == nil {
		return pngI, nil
	}
	jpegI, err := jpeg.Decode(bytes.NewReader(body))
	if err == nil {
		return jpegI, nil
	}
	gifI, err := gif.Decode(bytes.NewReader(body))
	if err == nil {
		return gifI, nil
	}
	return nil, err

}

type CustChannelInfo struct {
	CustId               uint32
	UnionId              string
	PaOpenId             string
	PaSubscribed         bool
	PaFirstSubscribeTime *timestamp.Timestamp
	PaLastSubscribeTime  *timestamp.Timestamp
	MbOpenId             string
	MbBindPhone          string
	MbBindTime           *timestamp.Timestamp
	DzWeixinId           string
	OwnerAccountId       int32
	OwnerAccountName     string
}

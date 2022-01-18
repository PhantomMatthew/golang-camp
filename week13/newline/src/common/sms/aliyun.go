package internal

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

// SMSAliyun SMSAliyun
type SMSAliyun struct {
	AccessKeyID, AccessKeySecret, RegionID, CaptchaSignName, CaptchaTemplateCode string
}

// Init Init
func (sa *SMSAliyun) Init(accessKeyID, accessKeySecret, regionID, captchaSignName, captchaTemplateCode string) {
	sa.AccessKeyID = accessKeyID
	sa.AccessKeySecret = accessKeySecret
	sa.RegionID = regionID
	sa.CaptchaSignName = captchaSignName
	sa.CaptchaTemplateCode = captchaTemplateCode
}

// SendSMS SendSMS
func (sa *SMSAliyun) SendSMS(signName, templateCode, mobile, content string) (string, error) {
	client, err := sdk.NewClientWithAccessKey(sa.RegionID, sa.AccessKeyID, sa.AccessKeySecret)
	if err != nil {
		panic(err)
	}

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = sa.RegionID
	request.QueryParams["PhoneNumbers"] = mobile
	request.QueryParams["SignName"] = signName
	request.QueryParams["TemplateCode"] = templateCode
	request.QueryParams["TemplateParam"] = content

	response, err := client.ProcessCommonRequest(request)
	return response.GetHttpContentString(), err
}

// SendCaptcha SendCaptcha
func (sa *SMSAliyun) SendCaptcha(mobile, captcha string) (string, error) {
	return sa.SendSMS(sa.CaptchaSignName, sa.CaptchaTemplateCode, mobile, `{"code":"`+captcha+`"}`)
}

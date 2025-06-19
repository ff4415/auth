package sms_provider

import (
	"fmt"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type TencentAuth struct {
	SecretId   string
	SecretKey  string
	SdkAppId   string
	SignName   string
	TemplateId string
	client     *sms.Client
}

func NewTencentAuth(secretId, secretKey, sdkAppId, signName, templateId string) (*TencentAuth, error) {
	// 创建认证对象
	credential := common.NewCredential(secretId, secretKey)

	// 创建客户端配置
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"

	// 创建腾讯云 SMS 客户端
	client, err := sms.NewClient(credential, "ap-beijing", cpf)
	if err != nil {
		return nil, fmt.Errorf("failed to create tencent sms client: %v", err)
	}

	return &TencentAuth{
		SecretId:   secretId,
		SecretKey:  secretKey,
		SdkAppId:   sdkAppId,
		SignName:   signName,
		TemplateId: templateId,
		client:     client,
	}, nil
}

func (t *TencentAuth) SendSMS(phoneNumber string, otp string) error {
	// 生成验证码
	code := otp

	// 创建发送短信请求
	request := sms.NewSendSmsRequest()

	// 设置手机号码（需要包含国家码）
	if phoneNumber[0] != '+' {
		phoneNumber = "+86" + phoneNumber
	}
	request.PhoneNumberSet = []*string{&phoneNumber}

	// 设置应用ID
	request.SmsSdkAppId = &t.SdkAppId

	// 设置模板ID
	request.TemplateId = &t.TemplateId

	// 设置签名
	request.SignName = &t.SignName

	// 设置模板参数（验证码和有效期）
	// 根据调试结果，模板需要2个参数：验证码和有效期
	validityPeriod := "5" // 5分钟有效期
	request.TemplateParamSet = []*string{&code, &validityPeriod}

	// 发送短信
	response, err := t.client.SendSms(request)
	if err != nil {
		return fmt.Errorf("failed to send sms: %v, phoneNumber: %s, code: %s, utc: %s", err, phoneNumber, code, time.Now().UTC().Format(time.RFC3339))
	}

	// 检查发送状态
	if response.Response == nil || len(response.Response.SendStatusSet) == 0 {
		return fmt.Errorf("send sms failed: empty response")
	}

	sendStatus := response.Response.SendStatusSet[0]
	if sendStatus.Code == nil || *sendStatus.Code != "Ok" {
		message := "unknown error"
		if sendStatus.Message != nil {
			message = *sendStatus.Message
		}
		code := "unknown"
		if sendStatus.Code != nil {
			code = *sendStatus.Code
		}
		return fmt.Errorf("send sms failed, code: %s, message: %s", code, message)
	}

	return nil
}

// SendSMSWithTemplate 发送带自定义模板参数的短信
func (t *TencentAuth) SendSMSWithTemplate(phoneNumber string, templateParams []string) error {
	// 创建发送短信请求
	request := sms.NewSendSmsRequest()

	// 设置手机号码
	if phoneNumber[0] != '+' {
		phoneNumber = "+86" + phoneNumber
	}
	request.PhoneNumberSet = []*string{&phoneNumber}

	// 设置应用ID
	request.SmsSdkAppId = &t.SdkAppId

	// 设置模板ID
	request.TemplateId = &t.TemplateId

	// 设置签名
	request.SignName = &t.SignName

	// 设置模板参数
	paramPtrs := make([]*string, len(templateParams))
	for i, param := range templateParams {
		paramPtrs[i] = &param
	}
	request.TemplateParamSet = paramPtrs

	// 发送短信
	response, err := t.client.SendSms(request)
	if err != nil {
		return fmt.Errorf("failed to send sms: %v", err)
	}

	// 检查发送状态
	if response.Response == nil || len(response.Response.SendStatusSet) == 0 {
		return fmt.Errorf("send sms failed: empty response")
	}

	sendStatus := response.Response.SendStatusSet[0]
	if sendStatus.Code == nil || *sendStatus.Code != "Ok" {
		message := "unknown error"
		if sendStatus.Message != nil {
			message = *sendStatus.Message
		}
		code := "unknown"
		if sendStatus.Code != nil {
			code = *sendStatus.Code
		}
		return fmt.Errorf("send sms failed, code: %s, message: %s", code, message)
	}

	return nil
}

// GetSendStatus 获取短信发送状态
func (t *TencentAuth) GetSendStatus(limit uint64) ([]*sms.PullSmsSendStatus, error) {
	request := sms.NewPullSmsSendStatusRequest()
	request.Limit = &limit
	request.SmsSdkAppId = &t.SdkAppId

	response, err := t.client.PullSmsSendStatus(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get send status: %v", err)
	}

	if response.Response == nil {
		return nil, fmt.Errorf("empty response")
	}

	return response.Response.PullSmsSendStatusSet, nil
}

// GetReplyStatus 获取短信回复状态
func (t *TencentAuth) GetReplyStatus(limit uint64) ([]*sms.PullSmsReplyStatus, error) {
	request := sms.NewPullSmsReplyStatusRequest()
	request.Limit = &limit
	request.SmsSdkAppId = &t.SdkAppId

	response, err := t.client.PullSmsReplyStatus(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get reply status: %v", err)
	}

	if response.Response == nil {
		return nil, fmt.Errorf("empty response")
	}

	return response.Response.PullSmsReplyStatusSet, nil
}

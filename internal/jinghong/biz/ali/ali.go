package ali

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"jonghong/internal/pkg/log"
	"time"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/spf13/viper"
)

type AliBiz interface {
	// 上传和下载文件的功能放在前端直接调用oss的api
	GetPolicyToken(username string, purpose string) (string, error)
	GetOssHost() string
}

type aliBiz struct {
	region     string
	bucketName string
	product    string
}

type policyToken struct {
	Policy           string `json:"policy"`
	SecurityToken    string `json:"security_token"`
	SignatureVersion string `json:"x_oss_signature_version"`
	Credential       string `json:"x_oss_credential"`
	Date             string `json:"x_oss_date"`
	Signature        string `json:"signature"`
	Host             string `json:"host"`
	Dir              string `json:"dir"`
	Callback         string `json:"callback"`
}

type callbackParam struct {
	CallbackUrl      string `json:"callbackUrl"`
	CallbackBody     string `json:"callbackBody"`
	CallbackBodyType string `json:"callbackBodyType"`
}

// 先写死初始化，等后面再抽出来写成配置项
func NewAliBiz() AliBiz {
	return &aliBiz{
		region:     "cn-shanghai",
		bucketName: "ai-tang",
		product:    "oss",
	}
}

func (ab *aliBiz) GetPolicyToken(username string, purpose string) (string, error) {

	host := fmt.Sprintf("https://%s.oss-%s.aliyuncs.com", ab.bucketName, ab.region)

	dir := "jinghong"

	callbackUrl := fmt.Sprintf("https://api.honghouse.cn/ali/%s/callback", purpose)
	// 安全信息配置
	config := new(credentials.Config).
		SetType("ram_role_arn").                                          //设置认证类型
		SetAccessKeyId(viper.GetString("ali.oss-access-key-id")).         //设置Id
		SetAccessKeySecret(viper.GetString("ali.oss-access-key-secret")). //设置secret
		SetRoleArn(viper.GetString("ali.oss-sts-role-arn")).              // 设置权限角色
		SetRoleSessionName("Role_Session_Name").
		SetPolicy("").
		SetRoleSessionExpiration(3600)

	provider, err := credentials.NewCredential(config)

	if err != nil {
		log.Fatalw("NewCredential fail, err:%v", err)
		return "", err
	}

	cred, err := provider.GetCredential()

	if err != nil {
		log.Fatalw("GetCredential fail, err:%v", err)
		return "", err
	}

	// 构建token策略
	utcTime := time.Now().UTC()
	date := utcTime.Format("20060102")
	expiration := utcTime.Add(1 * time.Hour)
	// 都是一些默认必要参数，没有额外的需求就不用变
	policyMap := map[string]any{
		"expiration": expiration.Format("2006-01-02T15:04:05.000Z"),
		"conditions": []any{
			map[string]string{"bucket": ab.bucketName},
			map[string]string{"x-oss-signature-version": "OSS4-HMAC-SHA256"},
			map[string]string{"x-oss-credential": fmt.Sprintf("%v/%v/%v/%v/aliyun_v4_request", *cred.AccessKeyId, date, ab.region, ab.product)},
			map[string]string{"x-oss-date": utcTime.Format("20060102T150405Z")},
			map[string]string{"x-oss-security-token": *cred.SecurityToken},
		},
	}

	policy, err := json.Marshal(policyMap)
	if err != nil {
		log.Fatalw("json.Marshal fail, err:%v", err)
		return "", err
	}

	// 构造待签名字符串
	stringToSign := base64.StdEncoding.EncodeToString([]byte(policy))
	hmacHash := func() hash.Hash {
		return sha256.New()
	}

	// 构建signing key  固定步骤按照文档写就ok了
	signingKey := "aliyun_v4" + *cred.AccessKeySecret
	h1 := hmac.New(hmacHash, []byte(signingKey))
	io.WriteString(h1, date)
	h1Key := h1.Sum(nil)

	h2 := hmac.New(hmacHash, h1Key)
	io.WriteString(h2, ab.region)
	h2Key := h2.Sum(nil)

	h3 := hmac.New(hmacHash, h2Key)
	io.WriteString(h3, ab.product)
	h3Key := h3.Sum(nil)

	h4 := hmac.New(hmacHash, h3Key)
	io.WriteString(h4, "aliyun_v4_request")
	h4Key := h4.Sum(nil)

	// 生成签名
	h := hmac.New(hmacHash, h4Key)
	io.WriteString(h, stringToSign)
	signature := hex.EncodeToString(h.Sum(nil))

	// 设置回调信息
	var callbackParam callbackParam
	callbackParam.CallbackUrl = callbackUrl
	callbackParam.CallbackBody = fmt.Sprintf("filename=${object}&size=${size}&mimeType=${mimeType}&username=%s", username)
	callbackParam.CallbackBodyType = "application/x-www-form-urlencoded"
	callback_str, err := json.Marshal(callbackParam)
	if err != nil {
		fmt.Println("callback json err:", err)
		return "", err
	}
	callbackBase64 := base64.StdEncoding.EncodeToString(callback_str)
	policyToken := policyToken{
		Policy:           stringToSign,
		SecurityToken:    *cred.SecurityToken,
		SignatureVersion: "OSS4-HMAC-SHA256",
		Credential:       fmt.Sprintf("%v/%v/%v/%v/aliyun_v4_request", *cred.AccessKeyId, date, ab.region, ab.product),
		Date:             utcTime.Format("20060102T150405Z"),
		Signature:        signature,
		Host:             host,           // 返回 OSS 上传地址
		Dir:              dir,            // 返回上传目录
		Callback:         callbackBase64, // 返回上传回调参数
	}
	response, err := json.Marshal(policyToken)
	if err != nil {
		fmt.Println("json err:", err)
		return "", err
	}
	return string(response), nil
}

func (ab *aliBiz) GetOssHost() string {
	return fmt.Sprintf("https://%s.oss-%s.aliyuncs.com", ab.bucketName, ab.region)
}

package aliclient

import (
	"context"
	"sync"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	openapicred "github.com/aliyun/credentials-go/credentials"
	"github.com/spf13/viper"
)

var (
	OssClient *oss.Client
	once      sync.Once
)

func NewOssClient() *oss.Client {
	once.Do(
		func() {
			config := new(openapicred.Config).
				// 填写Credential类型，固定值为ram_role_arn
				SetType("ram_role_arn").
				// 从环境变量中获取RAM用户的访问密钥（AccessKeyId和AccessKeySecret）
				SetAccessKeyId(viper.GetString("ali.oss-access-key-id")).         //设置Id
				SetAccessKeySecret(viper.GetString("ali.oss-access-key-secret")). //设置secret
				SetRoleArn(viper.GetString("ali.oss-sts-role-arn")).              // 设置权限角色
				// 以下操作默认直接填入参数数值，您也可以通过添加环境变量，并使用os.Getenv("<变量名称>")的方式来set对应参数
				// 自定义角色会话名称，用于区分不同的令牌
				SetRoleSessionName("ALIBABA_CLOUD_ROLE_SESSION_NAME"). // RoleSessionName默认环境变量规范名称ALIBABA_CLOUD_ROLE_SESSION_NAME
				//（可选）限制STS Token的有效时间
				SetRoleSessionExpiration(3600)

			arnCredential, gerr := openapicred.NewCredential(config)
			provider := credentials.CredentialsProviderFunc(func(ctx context.Context) (credentials.Credentials, error) {
				if gerr != nil {
					return credentials.Credentials{}, gerr
				}
				cred, err := arnCredential.GetCredential()
				if err != nil {
					return credentials.Credentials{}, err
				}
				return credentials.Credentials{
					AccessKeyID:     *cred.AccessKeyId,
					AccessKeySecret: *cred.AccessKeySecret,
					SecurityToken:   *cred.SecurityToken,
				}, nil
			})

			// 加载默认配置并设置凭证提供者和region
			cfg := oss.LoadDefaultConfig().
				WithCredentialsProvider(provider).
				WithRegion(viper.GetString("ali.oss-region"))

			// 创建OSS客户端
			OssClient = oss.NewClient(cfg)
		},
	)

	return OssClient

}

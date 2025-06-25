package ali

import (
	"bytes"
	"context"
	"io"
	aliclient "jonghong/internal/pkg/ali_client"
	"jonghong/internal/pkg/errno"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type AliBiz interface {
	// 上传和下载文件的功能放在前端直接调用oss的api
	PutObject(fileBody []byte, fileName string, contentType string, ctx context.Context) error
	GetObject(fileName string, ctx context.Context) ([]byte, string, error)
}

type aliBiz struct {
	bucketName string
	dirName    string
	ossClient  *oss.Client
}

type policyToken struct {
	Policy           string `json:"policy"`
	SecurityToken    string `json:"x-oss-security-token"`
	SignatureVersion string `json:"x-oss-signature-version"`
	Credential       string `json:"x-oss-credential"`
	Date             string `json:"x-oss-date"`
	Signature        string `json:"x-oss-signature"`
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
		bucketName: "ai-tang",
		dirName:    "jinghong/",
		ossClient:  aliclient.NewOssClient(),
	}
}

func (ab *aliBiz) PutObject(fileBody []byte, fileName string, contentType string, ctx context.Context) error {
	request := &oss.PutObjectRequest{
		Bucket:      oss.Ptr(ab.bucketName),
		Key:         oss.Ptr(ab.dirName + fileName),
		ContentType: oss.Ptr(contentType),
		Body:        bytes.NewReader(fileBody),
	}

	result, err := ab.ossClient.PutObject(ctx, request)

	if err != nil {
		return err
	}

	if result.StatusCode != 200 {
		return errno.InternalServerError.SetMessage("%s", result.Status)
	}

	return nil
}

func (ab *aliBiz) GetObject(fileName string, ctx context.Context) ([]byte, string, error) {

	request := &oss.GetObjectRequest{
		Bucket: oss.Ptr(ab.bucketName),
		Key:    oss.Ptr(ab.dirName + fileName),
	}

	result, err := ab.ossClient.GetObject(ctx, request)

	if err != nil {
		return nil, "", err
	}
	defer result.Body.Close()
	data, err := io.ReadAll(result.Body)

	if err != nil {
		return nil, "", err
	}

	return data, *result.ContentType, nil
}

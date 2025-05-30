package token

import (
	"errors"
	"fmt"
	"jonghong/internal/pkg/errno"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 因为是外部包所以单独设置错误

var ErrMissingHeader = errors.New("missing `Authorization` header")

// 采用jwt的方式进行token的生成和验证

type Config struct {
	key   string
	idKey string // 用于存储用户名信息
}

var (
	config = Config{"f0fe25ede84b9d3e1216417341c3f766", "idKey"}
	once   sync.Once //用于创建应用的单例对象
)

// 创建应用的单例对象

func Init(key, idKey string) error {
	once.Do(
		func() {
			if key != "" {
				config.key = key
			}
			if idKey != "" {
				config.idKey = idKey
			}
		})
	return nil
}

// 1. 签发token的关键信息主要有
// 	1）加密算法
//	2）载荷（claims）标准字段有 nbf 生效时间, iat 签发时间, exp 失效时间; 可以额外增加信息，以便后续解析出来
//  3）密钥，留存资服务端进行加密使用

func Sign(idKey string, expTime int64) (tokenString string, err error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			config.idKey: idKey,
			"nbf":        time.Now().Unix(),
			"iat":        time.Now().Unix(),
			"exp":        time.Now().Add(7 * 24 * time.Hour).Unix(),
		},
	)
	tokenString, err = token.SignedString([]byte(config.key))
	return
}

// 2. 解析token
// 1）验证加密算法是否为服务端设置的算法
// 2）验证服务端的密钥
// 3）提取额外的用户信息
func Parse(tokenString string) (string, error) {

	key := config.key

	// 先从tokenString中解析出token相关的信息，如加密方法，载荷信息
	// 解析的第二个参数是提供验证token相关信息并提供密钥的函数
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		// 验证加密算法一致
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(key), nil
	})
	if err != nil {
		return "", err
	}

	var identityKey string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		identityKey = claims[config.idKey].(string)
		if exp, ok := claims["exp"].(float64); ok {
			expirationTime := time.Unix(int64(exp), 0)
			if time.Now().UTC().After(expirationTime) {
				return "", errno.ErrTokenExpired
			}
		} else {
			return "", errno.ErrTokenInvalid.SetMessage("token失效时间获取失败")
		}

	} else {
		return "", errno.ErrTokenInvalid.SetMessage("token格式无效")
	}

	return identityKey, nil
}

// 3. 对Parse进一步封装，获取网络来的数据进行parse

func ParseRequest(c *gin.Context) (string, error) {
	header := c.Request.Header.Get("Authorization")
	if len(header) == 0 {
		return "", ErrMissingHeader
	}
	var t string
	// 按照字符串模板提取出tokenString
	fmt.Sscanf(header, "Bearer %s", &t)
	return Parse(t)
}

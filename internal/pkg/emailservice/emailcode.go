package emailservice

import (
	"crypto/rand"
	"jonghong/internal/pkg/errno"
	"math/big"
	"sync"
	"time"
)

type EmailCode struct {
	codes         map[string]codeInfo
	cleanupTicker time.Ticker
	mutex         sync.Mutex
	running       bool
}

type codeInfo struct {
	code        string
	expiredTime time.Time
}

func (ec *EmailCode) generateCode(key string) (string, error) {
	const digits = "0123456789"
	code := make([]byte, 6)
	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		code[i] = digits[num.Int64()]
	}
	// 失效时间为5分钟
	expiredAt := time.Now().Add(5 * time.Minute)
	ec.codes[key] = codeInfo{code: string(code), expiredTime: expiredAt}
	return string(code), nil
}

func (ec *EmailCode) verifyCode(key string, value string) (bool, error) {
	ec.mutex.Lock()
	defer ec.mutex.Unlock()
	v, exist := ec.codes[key]
	if !exist {
		return false, errno.ErrUserNotFound
	}
	if v.code != value {
		return false, errno.ErrCodeNotExist
	}

	if time.Now().After(v.expiredTime) {
		return false, errno.ErrCodeExpired
	}
	ec.cleanCode(key)
	return true, nil
}

func (ec *EmailCode) cleanCode(key string) bool {

	delete(ec.codes, key)

	return true
}

func (ec *EmailCode) cleanExpiredCodes() {
	for range ec.cleanupTicker.C {
		if !ec.running {
			break
		}
		ec.mutex.Lock()
		now := time.Now()
		for username, code := range ec.codes {
			if now.After(code.expiredTime) {
				delete(ec.codes, username)
			}
		}
		ec.mutex.Unlock()
	}
}

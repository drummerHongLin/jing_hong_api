package store

import (
	"sync"

	"gorm.io/gorm"
)

// 定义一个总的store
// 各个功能体的store都从这个store中获取

var (
	once sync.Once
	S    IStore
)

type IStore interface {
	Users() UserStore
	DB() *gorm.DB
	Chat() ChatStore
}

type databstore struct {
	db *gorm.DB
}

func (ds *databstore) Users() UserStore {
	return newUsers(ds.db)
}

func (ds *databstore) DB() *gorm.DB {
	return ds.db
}

func (ds *databstore) Chat() ChatStore {
	return newChat(ds.db)
}

func NewStore(db *gorm.DB) IStore {
	once.Do(func() {
		S = &databstore{db: db}
	})
	return S
}

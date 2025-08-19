package store

import (
	"context"
	"jonghong/internal/pkg/model"

	"gorm.io/gorm"
)

type PlayerStore interface {
	CreatePlayConfig(ctx context.Context, playConfig *model.PlayConfigM) error
	UpdatePlayConfig(ctx context.Context, playConfig *model.PlayConfigM) error
	GetPlayConfigs(ctx context.Context, offset int, limit int, userId int) ([]model.PlayConfigM, error)
	InsertPlayConfigs(ctx context.Context, playConfigs []model.PlayConfigM) error
	DeletePlayConfig(ctx context.Context, userId int, playConfigNo string) error
}

type player struct {
	db *gorm.DB
}

func newPlayer(db *gorm.DB) PlayerStore {
	return &player{db: db}
}

func (p *player) CreatePlayConfig(ctx context.Context, playConfig *model.PlayConfigM) error {
	return p.db.Create(playConfig).Error
}

func (p *player) UpdatePlayConfig(ctx context.Context, playConfig *model.PlayConfigM) error {
	return p.db.Model(&model.PlayConfigM{}).Where("playConfigNo = ?", playConfig.PlayConfigNo).Updates(playConfig).Error
}

func (p *player) GetPlayConfigs(ctx context.Context, offset int, limit int, userId int) ([]model.PlayConfigM, error) {
	var playConfigs []model.PlayConfigM
	err := p.db.Where("userId = ?", userId).Order("updatedAt desc").Offset(offset).Limit(limit).Find(&playConfigs).Error
	return playConfigs, err
}

func (p *player) InsertPlayConfigs(ctx context.Context, playConfigs []model.PlayConfigM) error {
	return p.db.Create(&playConfigs).Error
}

func (p *player) DeletePlayConfig(ctx context.Context, userId int, playConfigNo string) error {
	return p.db.Where("userId = ?", userId).Where("playConfigNo = ?", playConfigNo).Delete(&model.PlayConfigM{}).Error
}

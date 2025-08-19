package player

import (
	"context"
	"jonghong/internal/jinghong/store"
	"jonghong/internal/pkg/model"
	v1 "jonghong/pkg/api/jinghong/v1"

	"github.com/jinzhu/copier"
)

type PlayerBiz interface {
	CreatePlayConfig(ctx context.Context, r *v1.PlayConfig, userId int) error
	UpdatePlayConfig(ctx context.Context, r *v1.PlayConfig, userId int) error
	GetPlayConfigs(ctx context.Context, offset int, limit int, userId int) (*v1.GetPlayConfigsResponse, error)
	InsertPlayConfigs(ctx context.Context, playConfigs []v1.PlayConfigWithTime, userId int) error
	DeletePlayConfig(ctx context.Context, userId int, playConfigNo string) error
}

type player struct {
	ps store.PlayerStore
}

func NewPlayerBiz(ps store.PlayerStore) PlayerBiz {
	return &player{ps: ps}
}

func (p *player) CreatePlayConfig(ctx context.Context, r *v1.PlayConfig, userId int) error {
	var newPlayConfig model.PlayConfigM
	_ = copier.Copy(&newPlayConfig, r)
	newPlayConfig.UserId = userId
	return p.ps.CreatePlayConfig(ctx, &newPlayConfig)
}

func (p *player) UpdatePlayConfig(ctx context.Context, r *v1.PlayConfig, userId int) error {
	var newPlayConfig model.PlayConfigM
	_ = copier.Copy(&newPlayConfig, r)
	newPlayConfig.UserId = userId
	return p.ps.UpdatePlayConfig(ctx, &newPlayConfig)
}

func (p *player) GetPlayConfigs(ctx context.Context, offset int, limit int, userId int) (*v1.GetPlayConfigsResponse, error) {
	playerConfigs, err := p.ps.GetPlayConfigs(ctx, offset, limit, userId)
	if err != nil {
		return nil, err
	}
	hasMore := len(playerConfigs) >= limit

	var res []v1.PlayConfigWithTime
	_ = copier.Copy(&res, playerConfigs)
	return &v1.GetPlayConfigsResponse{
		HasMore:     hasMore,
		PlayConfigs: res,
	}, nil
}

func (p *player) InsertPlayConfigs(ctx context.Context, playConfigs []v1.PlayConfigWithTime, userId int) error {
	var pm []model.PlayConfigM
	_ = copier.Copy(&pm, playConfigs)
	for i := range pm {
		pm[i].UserId = userId
	}

	return p.ps.InsertPlayConfigs(ctx, pm)

}

func (p *player) DeletePlayConfig(ctx context.Context, userId int, playConfigNo string) error {
	return p.ps.DeletePlayConfig(ctx, userId, playConfigNo)
}

package service

import (
	"context"

	"github.com/BiTaksi/drivercampaign/internal/model"
)

type ICampaignService interface {
	GetCampaignsHome(ctx context.Context, dto *model.GetCampaignsHome) (*model.Data, error)
}

type campaignService struct {
}

func NewCampaignService() ICampaignService {
	return &campaignService{}
}

func (c *campaignService) GetCampaignsHome(ctx context.Context, dto *model.GetCampaignsHome) (*model.Data, error) {
	return &model.Data{
		MultiCampaign: model.MultiCampaign{
			Title:      "Yeni kampanyalar var!",
			ButtonText: "Kampanyaları Gör",
		},
		SingleCampaign: model.SingleCampaign{
			ID:         "ahsfhs",
			Title:      "Yeni kampanyalar var!",
			ButtonText: "Kampanyaları Gör",
		},
		ActiveCampaign: []model.ActiveCampaign{
			{
				ID:             "ashfhds",
				NumberOfTasks:  8,
				CompletedTasks: 1,
				Reward:         "150 TL",
				IconName:       "YELLOW_TAXI",
			},
		},
	}, nil
}

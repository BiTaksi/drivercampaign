package request

import "github.com/BiTaksi/drivercampaign/internal/model"

type GetCampaignsHome struct {
	DriverID string `json:"-"`
}

func (gdc *GetCampaignsHome) ToEntity() *model.GetCampaignsHome {
	return &model.GetCampaignsHome{
		DriverID: gdc.DriverID,
	}
}

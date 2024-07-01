package model

type GetCampaignsHome struct {
	DriverID string `json:"-"`
}

type Data struct {
	MultiCampaign  MultiCampaign    `json:"multiCampaign"`
	SingleCampaign SingleCampaign   `json:"singleCampaign"`
	ActiveCampaign []ActiveCampaign `json:"activeCampaign"`
}

// MultiCampaign struct
type MultiCampaign struct {
	Title      string `json:"title"`
	ButtonText string `json:"buttonText"`
}

// SingleCampaign struct
type SingleCampaign struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	ButtonText string `json:"buttonText"`
}

// ActiveCampaign struct
type ActiveCampaign struct {
	ID             string `json:"id"`
	NumberOfTasks  int    `json:"numberOfTasks"`
	CompletedTasks int    `json:"completedTasks"`
	Reward         string `json:"reward"`
	IconName       string `json:"iconName"`
}

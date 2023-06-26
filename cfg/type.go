package cfg

type Configuration struct {
	Development            bool   `json:"development"`
	ServerID               string `json:"serverId"`
	ReportChannelID        string `json:"reportChannelId"`
	ResoledReportChannelID string `json:"resoledReportChannelId"`
	EmbedColors            struct {
		Default int `json:"default"`
		Error   int `json:"error"`
	} `json:"embedColors"`
}

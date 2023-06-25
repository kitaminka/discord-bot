package cfg

type Configuration struct {
	Development bool   `json:"development"`
	ServerID    string `json:"serverId"`
	EmbedColors struct {
		Default int `json:"default"`
		Error   int `json:"error"`
	} `json:"embedColors"`
}

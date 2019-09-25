package model

type Config struct {
	Space        string             `json:"space"`
	MessageQuery MessageQueryConfig `json:"message_query_config"`
	Gateway      GatewayConfig      `json:"gateway"`
}

type MessageQueryConfig struct {
	DataSource string `json:"url"`
}

type GatewayConfig struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

package configs

type Configuration struct {
	Server          ServerConfiguration
	Client          ClientConfiguration
	Telegram        TelegramBotConfiguration
	PostgreDatabase PostgreDatabaseConfiguration
	RedisDatabase   RedisDatabaseConfiguration
	Kraken          KrakenConfiguration
	KrakenWS        KrakenWSConfiguration
}

type ServerConfiguration struct {
	Port      string
	Websocket ServerWebsocketConfiguration
}

type ServerWebsocketConfiguration struct {
	ReadBufferSize  int
	WriteBufferSize int
	CheckOrigin     bool
}

type ClientConfiguration struct {
	URL string
}

type TelegramBotConfiguration struct {
	APIToken   string
	WebhookURL string
}

type PostgreDatabaseConfiguration struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type RedisDatabaseConfiguration struct {
	Host string
	Port string
}

type KrakenConfiguration struct {
	APIURL string
}

type KrakenWSConfiguration struct {
	Requests KrakenWSAPIRequestsConfiguration
	Kraken   KrakenWSAPIConfiguration
}

type KrakenWSAPIConfiguration struct {
	WSAPIURL string
}

type KrakenWSAPIRequestsConfiguration struct {
	WriteWaitInSeconds  int
	PongWaitInSeconds   int
	PingPeriodInSeconds int
	MaxMessageSize      int
}

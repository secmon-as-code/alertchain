package types

type Config struct {
	DB       DBConfig
	Server   ServerConfig
	LogLevel string
}

type DBConfig struct {
	Type   string
	Config string `zlog:"secret"`
}

type ServerConfig struct {
	Addr string
	Port int
}

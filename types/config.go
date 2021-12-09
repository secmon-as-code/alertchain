package types

type Config struct {
	DB DBConfig
}

type DBConfig struct {
	Type   string
	Config string `zlog:"secret"`
}

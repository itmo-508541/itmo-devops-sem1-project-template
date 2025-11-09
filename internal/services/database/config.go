package database

import "fmt"

type Config struct {
	Host     string
	Port     string
	SslMode  string
	Database string
	User     string
	Password string
	Timezone string
}

func (c Config) DataSourceName() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&TimeZone=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.SslMode,
		c.Timezone,
	)
}

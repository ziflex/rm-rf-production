package database

import (
	"strconv"
	"strings"
)

type Options struct {
	Name string
	Host string
	Port int
	User string
	Pass string
}

func toConnectionString(cfg Options) string {
	sb := new(strings.Builder)

	sb.WriteString("host=")
	sb.WriteString(cfg.Host)
	sb.WriteString(" port=")
	sb.WriteString(strconv.Itoa(cfg.Port))
	sb.WriteString(" dbname=")
	sb.WriteString(cfg.Name)
	sb.WriteString(" user=")
	sb.WriteString(cfg.User)
	sb.WriteString(" password=")
	sb.WriteString(cfg.Pass)
	sb.WriteString(" sslmode=disable")

	return sb.String()
}

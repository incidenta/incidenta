package server

type Config struct {
	Addr string `envconfig:"SERVER_ADDR" default:":8080"`
}

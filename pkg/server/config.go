package server

type Config struct {
	StaticAssets string `envconfig:"STATIC_ASSETS" default:"static"`
	Addr         string `envconfig:"SERVER_ADDR" default:":8080"`
}

package main

var Config ConfigClass

type ConfigClass struct {
	Server struct {
		BindAddr string `toml:"bind"`
		WWWRoot  string `toml:"wwwroot"`
	} `toml:"server"`
	NBCP struct {
		HomeURL string `toml:"home"`
	} `toml:"nbcp"`
}

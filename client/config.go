package client

import "github.com/BurntSushi/toml"

type Config struct {
	Node struct {
		ID   string
		Port string
		Path string
	}
	Server struct {
		Port string
	}
}

func ReadConfig(path string) (c Config, err error) {
	_, err = toml.DecodeFile(path, &c)
	return c, err
}

func (c *Config) ApplyDefaults() {

	if c.Node.ID == "" {
		c.Node.ID = "Node-0001"
	}

	if c.Node.Port == "" {
		c.Node.Port = "8888"
	}
	if c.Node.Path == "" {
		c.Node.Path = "data.json"
	}
	if c.Server.Port == "" {
		c.Server.Port = "7778"
	}
}

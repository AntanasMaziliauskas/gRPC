package node

import "github.com/BurntSushi/toml"

type Config struct {
	Node struct {
		ID     string
		Listen string
		Path   string
	}
	Server struct {
		Source string
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

	if c.Node.Listen == "" {
		c.Node.Listen = "0.0.0.0:8887"
	}
	if c.Node.Path == "" {
		c.Node.Path = "data.json"
	}
	if c.Server.Source == "" {
		c.Server.Source = "0.0.0.0:7778"
	}
}

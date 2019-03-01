package node

import "github.com/BurntSushi/toml"

//Config is a structure for a configuration file
type Config struct {
	Node struct {
		ID   string
		Path string
	}
	Server struct {
		Source string
	}
}

//ReadConfig function reads through and decodes the given file.
//Returns data in Config structure
func ReadConfig(path string) (c Config, err error) {
	_, err = toml.DecodeFile(path, &c)
	return c, err
}

//ApplyDefaults function check to see if we have all neccessary info from config file
//Assigns default values if these is some valuable information missing
func (c *Config) ApplyDefaults() {

	if c.Node.ID == "" {
		c.Node.ID = "Node-"
	}

	if c.Node.Path == "" {
		c.Node.Path = "data.json"
	}
	if c.Server.Source == "" {
		c.Server.Source = "0.0.0.0:7778"
	}
}

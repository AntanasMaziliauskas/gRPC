package client

import (
	"context"
	"log"

	"github.com/AntanasMaziliauskas/grpc/api"
	"google.golang.org/grpc"
)

type Application struct {
	conn *grpc.ClientConn
}

func (a *Application) Init() {
	var err error
	//Connecting to the server
	a.conn, err = grpc.Dial(":7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

}
func (a *Application) Start() {
	c := api.NewGreetingClient(a.conn)
	response, err := c.SayHello(context.Background(), &api.Handshake{Id: "Node003", Port: "7778"})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Timeout in: %d seconds", response.Timeout)
}

func (a *Application) Stop() {
	//closing connection to the server
	a.conn.Close()
}

//Pasiruosimas connectionui su serveriu |Init Susijungimas su serveriu

//Connectionas su serveriu | Start Pasisveikinimo issiuntimas serveriui

//Pasiruosimas jungti serveri |INIT

//Listeneris serverio |START

//Pasisveikinimo priimimas is Node kliento

//Pinginimo siuntimas serveriui | Start GO rutina

//Ateinancios uzklausos priimimas ir vykdymas - GetOne

/*MAINE SITIE DALYKAI
//NUSKAITOM CONFIG
if config, err = trapserver.ReadConfig(c.GlobalString("config")); err != nil {
		return err
	}
	//Defaults nustatom
	config.ApplyDefaults()

//SITIE APPLICATIONE PALIEKA
//ReadConfig function decodes config file
func ReadConfig(path string) (c Config, err error) {
	_, err = toml.DecodeFile(path, &c)
	return c, err
}
//Struktura config
type Config struct {
	SNMP struct {
		Listen string
	}
	HTTP struct {
		Listen string
	}
	Telegram struct {
		Enable    bool
		ChannelID int64
		Token     string
	}
// defaultai
func (c *Config) ApplyDefaults() {

	//SNMP
	if c.SNMP.Listen == "" {
		c.SNMP.Listen = "0.0.0.0:8000"
	}
	//HTTP
	if c.HTTP.Listen == "" {
		c.HTTP.Listen = "0.0.0.0:9162"
	}


	PATS CONFIGAS:
	name.toml
	HTTP]
listen = "0.0.0.0:8000"

[SNMP]
listen = "0.0.0.0:9162"

#You can set Telegram enable to true if you want to received Logrus logs via Telegram
#channelID is Telegram Channel ID
#token is Telegram BOT token needed for authenticatin the bot
[Telegram]
enable = false
channelID = -1
token = " "


*/

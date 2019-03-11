package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/AntanasMaziliauskas/grpc/node"
	"github.com/AntanasMaziliauskas/grpc/node/person"
)

func main() {
	var (
		config  node.Config
		err     error
		handler person.PersonService
	)

	conf := flag.String("config", "config.toml", "Config file to be used")
	storage := flag.String("storage", "memory", "Storage to be used. Can be 'memory' or 'mongo'")
	flag.Parse()

	if config, err = node.ReadConfig(*conf); err != nil {
		log.Fatalf("Could not read config file: %s", err)
	}
	config.ApplyDefaults()

	//Random Node name
	rand.Seed(time.Now().UnixNano())
	id := "Node-" + strconv.Itoa(rand.Intn(1000))

	handler = &person.DataFromMem{
		ID: id,
	}

	if *storage == "mongo" {
		handler = &person.DataFromMgo{
			ID: id,
		}
	}

	app := node.Application{
		ID:         id,
		ServerPort: config.Server.Source,
		Person:     handler,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM /*syscall.SIGSTOP,*/, syscall.SIGKILL)

	//TODO: Errors turetu ateiti iki cia.
	app.Init()

	app.Start()

	<-stop

	app.Stop()
}

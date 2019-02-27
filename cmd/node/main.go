package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AntanasMaziliauskas/grpc/node"
	"github.com/AntanasMaziliauskas/grpc/node/person"
)

func main() {
	var config node.Config
	var err error
	//Flag
	conf := flag.String("config", "config.toml", "Config file to be used")
	flag.Parse()

	if config, err = node.ReadConfig(*conf); err != nil {
		log.Fatalf("Could not read config file: %s", err)
	}
	config.ApplyDefaults()

	app := node.Application{
		//Port:       config.Node.Listen,
		ID:         config.Node.ID,
		ServerPort: config.Server.Source,
		Person: &person.DataFromFile{
			ID:   config.Node.ID,
			Path: config.Node.Path},
	}

	app.Init()

	app.Start()

	//app.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)

	<-stop

	/*
		var conn *grpc.ClientConn
		var wg *sync.WaitGroup

		//Flag
		id := flag.String("ID", "Node-001", "Node ID that is being sent to server")
		flag.Parse()
		//Connection to server
		conn, err := grpc.Dial(":7777", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		defer conn.Close()

		c := api.NewGreetingClient(conn)
		p := api.NewPingClient(conn)

		response, err := c.SayHello(context.Background(), &api.Handshake{Id: *id, Port: "7778"})
		if err != nil {
			log.Fatalf("Error when calling SayHello: %s", err)
		}
		log.Printf("Timeout in: %d seconds", response.Timeout)

		wg = &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			ticker := time.NewTicker(time.Duration(response.Timeout) / 2 * time.Second)
			for {
				select {
				case <-ticker.C:
					log.Printf("Pinging")
					_, err := p.PingMe(context.Background(), &api.PingMessage{Id: *id})
					if err != nil {
						log.Fatalf("Error when calling PingMe: %s", err)
					}
					//log.Printf("Response from server: %d", response)
					//wg.Done()
				}
			}
		}()

		//SERVER PART
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7778))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		// create a server instance
		s := handler.Server{}

		// create a gRPC server object
		grpcServer := grpc.NewServer()

		// attach the Ping service to the server
		api.RegisterLookForDataServer(grpcServer, &s)

		// start the server
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}

		//TODO STOP isideti
		//Stop isideti

		wg.Wait()*/
}

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/AntanasMaziliauskas/grpc/client"
)

//FLAGAs Config failas
//Config Node ID, Portas, kur jungtis, Koks mano portas.
func main() {

	app := client.Application{}

	app.Init()

	app.Start()

	app.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGKILL)

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

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/AntanasMaziliauskas/grpc/server"
	"github.com/AntanasMaziliauskas/grpc/server/broker"
)

// main start a gRPC server and waits for connection
func main() {

	app := server.Application{Broker: &broker.GRPCBroker{}}

	app.Init()

	app.Start()

	//app.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)

	<-stop

	/*var conn *grpc.ClientConn

	// create a listener on TCP port 7777
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7777))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a server instance
	s := handler.Server{}

	// create a gRPC server object
	grpcServer := grpc.NewServer()

	// attach the Greeting service to the server
	api.RegisterGreetingServer(grpcServer, &s)
	// attach the Ping service to the server
	api.RegisterPingServer(grpcServer, &s)
	// start the server
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()

	//Connection to server
	conn, err = grpc.Dial(":7778", grpc.WithInsecure()) // Portas ateina is NODE
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	a := api.NewLookForDataClient(conn)

	go func() {
		ticker := time.NewTicker(handler.TimeOut * time.Second)
		for {
			select {
			case <-ticker.C:
				//TODO CHECKPING TESTAS
				s.CheckPing()
				response, err := a.FindData(context.Background(), &api.LookFor{Id: "Jonas"})
				if err != nil {
					log.Fatalf("Error when trying to get response from server: %s", err)
				}
				log.Printf("Data received: %s", response)
			}
		}
	}()

	//start HTTP server
	http.HandleFunc("/", s.HTTPHandler)     // set router
	err = http.ListenAndServe(":8080", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	*/
}

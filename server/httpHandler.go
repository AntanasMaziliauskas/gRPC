package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/AntanasMaziliauskas/grpc/api"
	"google.golang.org/grpc"
)

//Pasiruosimas http serveriui | Init
func (a *Application) StartHTTPServer() {
	var err error

	http.HandleFunc("/list", a.HTTPHandlerList)
	http.HandleFunc("/getPerson", a.HTTPHandleGet) // set router
	err = http.ListenAndServe(":8080", nil)        // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (a *Application) HTTPHandleGet(w http.ResponseWriter, r *http.Request) {
	//	p, err := a.Broker.GetOnePersonBroadcast(...)
	//w.Write(json(p))
	if name := r.FormValue("name"); name != "" {
		fmt.Println(name)
		b := a.getOnePersonBroadcast(name)
		data, _ := json.Marshal(b)
		w.Write(data)
	}

}

func (a *Application) HTTPHandlerList(w http.ResponseWriter, r *http.Request) {

	b, _ := json.Marshal(Nodes)

	//	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

//BROKERYJE TURI BUTI
func (a *Application) getOnePersonBroadcast(name string) *api.Found {
	var response *api.Found
	//FOR einu per connections
	//for k, _ := range Nodes {
	conn, err := grpc.Dial(fmt.Sprintf(":8888"), grpc.WithInsecure()) // Portas ateina is NODE
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	g := api.NewLookForDataClient(conn)

	response, err = g.FindData(context.Background(), &api.LookFor{Name: name})
	if err != nil {
		log.Fatalf("Error when trying to get response from server: %s", err)
	}
	log.Printf("Data received from : %s", response)
	//	}
	return response
}

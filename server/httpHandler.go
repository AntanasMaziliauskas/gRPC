package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/AntanasMaziliauskas/grpc/api"
)

//StartHTTPServer sets handlers and start the ListenAndServe function
func (a *Application) StartHTTPServer() error {
	http.HandleFunc("/listnodes", a.HTTPHandlerListNodes)
	http.HandleFunc("/listpersons", a.HTTPHandlerListPersons)
	http.HandleFunc("/getperson/", a.HTTPHandleGetPerson)
	http.HandleFunc("/getpersonnode/", a.HTTPHandleGetPersonNode)

	go func() {
		defer a.wg.Done()
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()

	return nil
}

//HTTPHandleGet function get the data to be displayed
func (a *Application) HTTPHandleGetPersonNode(w http.ResponseWriter, r *http.Request) {
	var data []byte

	url := r.URL.String()
	s := strings.Split(url, "/")

	list, err := a.Broker.GetOnePersonNode(context.Background(), &api.Person{Id: s[2], Node: s[3]})

	if data, err = json.Marshal(list); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

//HTTPHandleGet function get the data to be displayed
func (a *Application) HTTPHandleGetPerson(w http.ResponseWriter, r *http.Request) {
	var data []byte

	url := r.URL.String()
	s := strings.Split(url, "/")

	list, err := a.Broker.GetOnePersonBroadcast(context.Background(), &api.Person{Id: s[2]})

	if data, err = json.Marshal(list); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	/*var (
		err  error
		b    *api.Person
		data []byte
	)

	if name := r.FormValue("name"); name != "" {
		if b, err = a.Broker.GetOnePersonBroadcast(name); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if data, err = json.Marshal(b); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, err = w.Write(data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}*/

}

//HTTPHandlerList get data to be displayed
func (a *Application) HTTPHandlerListNodes(w http.ResponseWriter, r *http.Request) {
	var (
		data []byte
	)

	list, err := a.Broker.ListNodes(context.Background(), &api.Empty{})

	if data, err = json.Marshal(list); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

//HTTPHandlerList get data to be displayed
func (a *Application) HTTPHandlerListPersons(w http.ResponseWriter, r *http.Request) {
	var (
		data []byte
	)

	list, err := a.Broker.ListPersonsBroadcast(context.Background(), &api.Empty{})

	if data, err = json.Marshal(list); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

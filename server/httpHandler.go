package server

import (
	"net/http"
)

//StartHTTPServer sets handlers and start the ListenAndServe function
func (a *Application) StartHTTPServer() error {
	http.HandleFunc("/list", a.HTTPHandlerList)
	http.HandleFunc("/getPerson", a.HTTPHandleGet)
	err := http.ListenAndServe(":8080", nil)

	return err
}

//HTTPHandleGet function get the data to be displayed
func (a *Application) HTTPHandleGet(w http.ResponseWriter, r *http.Request) {
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
func (a *Application) HTTPHandlerList(w http.ResponseWriter, r *http.Request) {
	/*	var (
			err  error
			data []byte
		)

		list := a.Broker.ListNodes()

		if data, err = json.Marshal(list); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, err = w.Write(data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}*/
}

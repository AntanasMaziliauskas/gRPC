package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/AntanasMaziliauskas/grpc/api"
	"golang.org/x/net/context"
)

// Server represents the gRPC server
type Server struct {
	Nodes []Node
}
type Node struct {
	ID             string
	LastConnection time.Time
}

const TimeOut = 30

// SayHello generates response to a Ping request
func (s *Server) SayHello(ctx context.Context, in *api.Handshake) (*api.Timeout, error) {
	log.Printf("Received message: %s. Port: %s", in.Id, in.Port)

	if !sliceContainsString(in.Id, s.Nodes) {
		nodeStruct := Node{
			ID:             in.Id,
			LastConnection: time.Now(),
		}
		s.Nodes = append(s.Nodes, nodeStruct)
	}

	return &api.Timeout{Timeout: TimeOut}, nil
}

func (s *Server) PingMe(ctx context.Context, in *api.PingMessage) (*api.Empty, error) {
	var newList []Node
	log.Printf("Ping Received from %s", in.Id)
	for _, v := range s.Nodes {
		if v.ID == in.Id {
			v.LastConnection = time.Now()
		}
		newList = append(newList, v)
	}
	s.Nodes = newList
	return &api.Empty{}, nil
}

//TODO SLICE PASIKEISTI I MAP
//Ne slice o map daryti, kad galeciau istrinti
func (s *Server) CheckPing() {
	var newList []Node
	for _, v := range s.Nodes {
		diff := time.Now().Sub(v.LastConnection)
		sec := int(diff.Seconds())
		if sec < TimeOut {
			newList = append(newList, v)
		}
	}
	s.Nodes = newList
}

// SliceContainsString will return true if needle has been found in haystack.
func sliceContainsString(needle string, haystack []Node) bool {
	for _, v := range haystack {
		if v.ID == needle {
			return true
		}
	}

	return false
}

func (s *Server) HTTPHandler(w http.ResponseWriter, r *http.Request) {

	b, _ := json.Marshal(s.Nodes)

	/*if name := r.FormValue("name"); name != "" {
		fmt.Println(name)
	}*/

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

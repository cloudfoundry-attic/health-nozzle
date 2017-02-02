package server

import (
	"encoding/json"
	"net/http"

	"github.com/cloudfoundry-incubator/health-nozzle/counter"
)

type summary struct {
	Ingress uint64 `json:"ingress"`
	Egress  uint64 `json:"egress"`
	Dropped uint64 `json:"dropped"`
}

type Server struct {
	counters []*counter.HealthCounter
}

func NewServer(counters []*counter.HealthCounter) *Server {
	return &Server{
		counters: counters,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	summaries := map[string]summary{}
	for _, counter := range s.counters {
		summaries[counter.GetName()] = summary{
			Ingress: counter.Ingress.Value(),
			Egress:  counter.Egress.Value(),
			Dropped: counter.Dropped.Value(),
		}
	}

	data, _ := json.Marshal(summaries)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

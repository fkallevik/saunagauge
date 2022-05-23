package server

import (
	"fmt"
	"html/template"
	"net/http"
)

// Logger defines a generic logging interface.
type Logger interface {
	Printf(string, ...any)
	Println(...any)
}

type Server struct {
	sauna  Sauna
	logger Logger
}

type Sauna interface {
	GetTemperature() float32
	GetHumidity() float32
}

func New(sauna Sauna, logger Logger) *Server {
	return &Server{
		sauna:  sauna,
		logger: logger,
	}
}

func (s *Server) handle(f func() htmlResponse) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := f()

		if err := respond(w, resp); err != nil {
			s.logger.Println(err)
		}
	}
}

func (s *Server) getTemperature() htmlResponse {
	return htmlResponse{
		html: template.HTML(fmt.Sprintf("%.1f", s.sauna.GetTemperature())),
	}
}

func (s *Server) getHumidity() htmlResponse {
	return htmlResponse{
		html: template.HTML(fmt.Sprintf("%.1f", s.sauna.GetTemperature())),
	}
}

func (s *Server) GetTemperature() func(w http.ResponseWriter, r *http.Request) {
	return s.handle(func() htmlResponse {
		return s.getTemperature()
	})
}

func (s *Server) GetHumidity() func(w http.ResponseWriter, r *http.Request) {
	return s.handle(func() htmlResponse {
		return s.getHumidity()
	})
}

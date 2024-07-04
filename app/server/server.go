package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
	"weathersvc/app/config"
	"weathersvc/app/service"
	_ "weathersvc/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	ln     net.Listener
	server *http.Server
	router *mux.Router
	Addr   string
}

type DecimalRequest struct {
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}
type NonDecimalRequest struct {
	Latitude  string `json:"latitude,omitempty"`
	Longitude string `json:"longitude,omitempty"`
}
type Response struct {
	Message   string
	Temp      string
	Condition string
	Wind      string
}

func NewServer(conf config.App, s *service.WeatherService) *Server {
	r := mux.NewRouter()
	r.HandleFunc("/weather", weatherHandler(s)).Methods("GET")
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	return &Server{
		server: &http.Server{
			Handler:           r,
			ReadHeaderTimeout: 3 * time.Second,
		},
		router: r,
		Addr:   fmt.Sprintf("0.0.0.0:%s", conf.Port),
	}
}

// Open validates the server options and begins listening on the bind address.
func (s *Server) Open() (err error) {
	if s.ln, err = net.Listen("tcp", s.Addr); err != nil {
		return fmt.Errorf("error listening, %w", err)
	}
	log.Printf("Server started listening for new connections:Port: %v", s.Port())
	if err := s.server.Serve(s.ln); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Close gracefully shuts down the server.
func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

// Port returns the TCP port for the running server.
// This is useful in tests where we allocate a random port by using ":0".
func (s *Server) Port() int {
	if s.ln == nil {
		return 0
	}
	return s.ln.Addr().(*net.TCPAddr).Port
}

// @title WeatherService API
// @version 1.0
// @description This is a sample server.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8001
// @BasePath /

// Handlers
// @Summary Local Weather Condition
// @Description Get the local weather condition by entering your latitude/longitude coordinates.
// @Success 200 {object} Response
// @Failure 500 {string} Internal Service Failure
// @Router /weather [get]
func weatherHandler(s *service.WeatherService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var inReq DecimalRequest
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &inReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		msg := ""
		if inReq.Latitude == 0 && inReq.Longitude == 0 {
			msg = "Lat and Lon are set to zero. Are you enjoying the ocean? "
		}
		wResp, err := s.GetWeather(r.Context(), inReq.Latitude, inReq.Longitude)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		msg = fmt.Sprintf("%sOutside it is %s with %s and %s.", msg, wResp.Temp, wResp.Wind, wResp.Condition)
		json.NewEncoder(w).Encode(Response{
			Message:   msg,
			Temp:      string(wResp.Temp),
			Condition: wResp.Condition,
			Wind:      string(wResp.Wind),
		})

	}
}

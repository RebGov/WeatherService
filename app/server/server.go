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
	apperrors "weathersvc/app/app_errors"
	"weathersvc/app/config"
	"weathersvc/app/service"
	_ "weathersvc/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server interface {
	Open() (err error)
	Close() error
	Port() int
}
type server struct {
	ln     net.Listener
	server *http.Server
	router *mux.Router
	Addr   string
}

type DecimalRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
type Response struct {
	Message   string
	Temp      string
	Condition string
	Wind      string
}

func NewServer(conf *config.App, s service.Service) Server {
	r := mux.NewRouter()
	r.HandleFunc("/weather/get", weatherHandler(s)).Methods("GET")
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	return &server{
		server: &http.Server{
			Handler:           r,
			ReadHeaderTimeout: 3 * time.Second,
		},
		router: r,
		Addr:   fmt.Sprintf("0.0.0.0:%s", conf.Port),
	}
}

// Open validates the server options and begins listening on the bind address.
func (s *server) Open() (err error) {
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
func (s *server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

// Port returns the TCP port for the running server.
// This is useful in tests where we allocate a random port by using ":0".
func (s *server) Port() int {
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
// @Failure 429 {string} ErrorResponse: Limit reached
// @Failure 404 {string} ErrorResponse: Coordinates not found
// @Failure 400 {string} ErrorResponse: Request invalid and reason
// @Router /weather/get [get]
func weatherHandler(s service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var inReq DecimalRequest
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if body == nil {
			http.Error(w, apperrors.ErrNoBody.Error(), http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &inReq)
		if err != nil {
			http.Error(w, apperrors.ErrNoBody.Error(), http.StatusBadRequest)
			return
		}
		if inReq.Latitude == 0 && inReq.Longitude == 0 {
			http.Error(w, apperrors.CreateInvalidRequestError("latitude and longitude missing or null").Error(), http.StatusBadRequest)
			return
		}
		if !isValidLat(inReq.Latitude) {
			http.Error(w, apperrors.CreateInvalidRequestError("latitude is out of range").Error(), http.StatusBadRequest)
			return
		}
		if !isValidLon(inReq.Longitude) {
			http.Error(w, apperrors.CreateInvalidRequestError("longitude is out of range").Error(), http.StatusBadRequest)
			return
		}
		wResp, err := s.GetWeather(r.Context(), inReq.Latitude, inReq.Longitude)
		if err != nil {
			switch err.Error() {
			case apperrors.ErrTooManyRequests.Error():
				http.Error(w, err.Error(), http.StatusTooManyRequests)
			case apperrors.ErrNotFound.Error():
				http.Error(w, err.Error(), http.StatusNotFound)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		msg := fmt.Sprintf("Outside it is %s with %s and %s.", wResp.Temp, wResp.Wind, wResp.Condition)
		json.NewEncoder(w).Encode(Response{
			Message:   msg,
			Temp:      string(wResp.Temp),
			Condition: wResp.Condition,
			Wind:      string(wResp.Wind),
		})

	}
}

func isValidLat(l float64) bool {
	if l < -90 || l > 90 {
		return false
	}
	return true
}
func isValidLon(l float64) bool {
	if l < -180 || l > 180 {
		return false
	}
	return true
}

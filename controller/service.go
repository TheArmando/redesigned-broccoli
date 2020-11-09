package controller

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	"../maxmind"
)

type Controller interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type NewParams struct {
	GeoDB maxmind.GeoDB
}

func New(p NewParams) Controller {
	return &controller{
		geoDB: p.GeoDB,
	}
}

type controller struct {
	geoDB maxmind.GeoDB
}

// Handle ...
func (c *controller) Handle(w http.ResponseWriter, r *http.Request) {
	// decode json payload into a struct
	var body IPValidateRequest

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println("error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// parse IP then pass into maxmind service
	ip := net.ParseIP(body.IPAddress)
	if ip == nil {
		log.Println("error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if c.geoDB.IsWhitelisted(ip, body.Countries) {
		w.WriteHeader(http.StatusOK)
		return
	}
	log.Println("ip not whitelisted", ip.String())
	w.WriteHeader(http.StatusUnauthorized)

}

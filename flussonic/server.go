package flussonic

import (
	"encoding/json"
	"net/http"
	"time"
)

type Server struct {
	RequestDuration float64 `json:"-"`
	Url string `json:"-"`
	TotalClients float64 `json:"total_clients"`
}

func (f *Flussonic) GetServer() (*Server, error) {
	client := &http.Client{}
	server := Server{}
	server.Url = "/flussonic/api/server"
	req, err := http.NewRequest("GET", f.Url.String()+server.Url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(f.User, f.Password)
	startTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	server.RequestDuration = time.Since(startTime).Seconds()
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&server)
	if err != nil {
		return nil, err
	}
	return &server, nil
}

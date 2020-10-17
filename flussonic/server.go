package flussonic

import (
	"encoding/json"
	"net/http"
)

type Server struct {
	TotalClients float64 `json:"total_clients"`
}

func (f *Flussonic) GetServer() (*Server, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", f.Url.String()+"/flussonic/api/server", nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(f.User, f.Password)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	server := Server{}
	err = json.NewDecoder(resp.Body).Decode(&server)
	if err != nil {
		return nil, err
	}
	return &server, nil
}

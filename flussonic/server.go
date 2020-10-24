/*
 *    Copyright 2020 Yury Makarov
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 *
 */

package flussonic

import (
	"encoding/json"
	"net/http"
	"time"
)

type Server struct {
	RequestDuration float64 `json:"-"`
	Url             string  `json:"-"`
	TotalClients    float64 `json:"total_clients"`
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
